import argparse
import fcntl
import json
import os
import re
import shutil
import subprocess
import sys
import tempfile
from pathlib import Path, PurePosixPath
from urllib.parse import quote


PR_FIELDS = "number,url,state,baseRefName,headRefName,headRefOid,isCrossRepository,title,body"


def run(*args, **kwargs):
    result = subprocess.run(args, check=True, stdout=subprocess.PIPE, **kwargs)
    return os.fsdecode(result.stdout).removesuffix("\n")


def git(*args, **kwargs):
    return run("git", *args, **kwargs)


def github_repo():
    url = git("config", "--get", "remote.origin.url")
    match = re.fullmatch(
        r"(?:https://github\.com/|ssh://git@github\.com/|git@github\.com:)"
        r"([a-zA-Z0-9_.-]+/[a-zA-Z0-9_.-]+?)(?:\.git)?/?",
        url,
    )
    if not match:
        raise ValueError("origin must point to a github.com repository")
    return f"github.com/{match[1]}"


def gh(repo, *args):
    return run("gh", "pr", *args, "--repo", repo)


def branch_ref(branch):
    return f"refs/heads/{branch}"


def commit_at(ref):
    return git("rev-parse", "--verify", "--end-of-options", f"{ref}^{{commit}}")


def default_branch():
    for line in git("ls-remote", "--symref", "origin", "HEAD").splitlines():
        ref, _, name = line.partition("\t")
        if name == "HEAD" and ref.startswith("ref: refs/heads/"):
            return ref.removeprefix("ref: refs/heads/")
    raise ValueError("could not find origin's default branch")


def fetch(*branches):
    tracking_refs = [f"refs/remotes/origin/{branch}" for branch in branches]
    refspecs = [
        f"+{branch_ref(branch)}:{tracking_ref}"
        for branch, tracking_ref in zip(branches, tracking_refs)
    ]
    git("fetch", "--atomic", "--no-tags", "--no-recurse-submodules", "origin", *refspecs)
    return [commit_at(ref) for ref in tracking_refs]


def changed_paths(parent, source, paths=()):
    return subprocess.check_output([
        "git", "--literal-pathspecs", "diff", "--name-only", "-z",
        "--no-renames", "--ignore-submodules=none", parent, source, "--", *paths,
    ])


def submission_commit(parent, source, paths, message):
    # Git supplies the NUL-delimited list so additions and deletions both work.
    changed = changed_paths(parent, source, paths)
    if not changed:
        return None

    with tempfile.TemporaryDirectory(prefix="homework-pr-") as temp_dir:
        env = dict(os.environ, GIT_INDEX_FILE=str(Path(temp_dir) / "index"))
        git("read-tree", parent, env=env)
        git(
            "--literal-pathspecs", "restore", f"--source={source}", "--staged",
            "--ignore-skip-worktree-bits", "--pathspec-from-file=-", "--pathspec-file-nul",
            input=changed, env=env,
        )
        tree = git("write-tree", env=env)
    if changed_paths(parent, tree) != changed:
        raise ValueError("selected paths overlap a file-to-directory change; select the parent path")
    return git("commit-tree", tree, "-p", parent, "-m", message)


def save_state(path, state):
    temporary = path.with_suffix(".tmp")
    temporary.write_text(json.dumps(state, indent=2) + "\n", encoding="utf-8")
    temporary.replace(path)


def remote_branches(*refs):
    branches = {}
    for line in git("ls-remote", "--refs", "origin", *refs).splitlines():
        oid, ref = line.split("\t")
        branches[ref] = oid
    return branches


def pr_text(name, state):
    source_branch = state["source_branch"]
    readme_url = f"https://{state['repo']}/blob/{quote(source_branch)}/README.md"
    body = (
        "This PR contains only the code selected for this homework.\n\n"
        f"To build the complete app, check out `{source_branch}` and follow the "
        f"[README]({readme_url})."
    )
    return f"Homework review ({name})", body


def reopen_pr(state, pr):
    base_ref = branch_ref(state["base_branch"])
    head_ref = branch_ref(state["head_branch"])
    branches = remote_branches(base_ref, head_ref)
    if branches.get(base_ref) != state["base"]:
        raise ValueError(f"the submission base branch is missing or has moved: {state['base_branch']}")
    if head_ref not in branches:
        git(
            "fetch", "--no-tags", "--no-recurse-submodules", "--no-write-fetch-head", "origin",
            f"refs/pull/{pr['number']}/head",
        )
        git(
            "push", "--no-follow-tags", f"--force-with-lease={head_ref}:",
            "origin", f"{pr['headRefOid']}:{head_ref}",
        )
    gh(state["repo"], "reopen", pr["url"])


def update_pr(name, state, pr):
    url = pr["url"]
    branches = pr["baseRefName"], pr["headRefName"]
    if pr["isCrossRepository"] or branches != (state["base_branch"], state["head_branch"]):
        raise ValueError(f"the PR no longer targets its original branches: {url}")
    if pr["state"] == "MERGED":
        raise ValueError(f"the submission PR has already been merged: {url}")
    if pr["state"] == "CLOSED":
        reopen_pr(state, pr)

    title, body = pr_text(name, state)
    if (pr["title"], pr["body"]) != (title, body):
        gh(state["repo"], "edit", url, "--title", title, "--body", body)


def publish(name, state):
    repo = state["repo"]
    base_branch = state["base_branch"]
    head_branch = state["head_branch"]
    base_ref, head_ref = branch_ref(base_branch), branch_ref(head_branch)
    expected = {base_ref: state["base"], head_ref: state["initial_head"]}
    existing = remote_branches(base_ref, head_ref)

    # Recover the PR if GitHub created it before an earlier invocation exited.
    prs = json.loads(gh(
        repo, "list", "--base", base_branch, "--head", head_branch, "--state", "all",
        "--json", PR_FIELDS,
    ))
    for pr in prs:
        if not pr["isCrossRepository"]:
            matching_head = pr["headRefOid"] == state["initial_head"]
            missing_head = pr["state"] == "CLOSED" and existing == {base_ref: state["base"]}
            if not matching_head or (existing != expected and not missing_head):
                raise ValueError(f"submission branches already have a PR at {pr['url']}")
            update_pr(name, state, pr)
            return pr["url"]

    if existing and existing != expected:
        raise ValueError("submission branch names are already in use on origin")
    if not existing:
        git(
            "push", "--atomic", "--no-follow-tags",
            f"--force-with-lease={base_ref}:", f"--force-with-lease={head_ref}:",
            "origin", f"{expected[base_ref]}:{base_ref}", f"{expected[head_ref]}:{head_ref}",
        )
    title, body = pr_text(name, state)
    return gh(
        repo, "create", "--base", base_branch, "--head", head_branch,
        "--title", title, "--body", body,
    )


def refresh_pr(name, state):
    pr = json.loads(gh(state["repo"], "view", state["url"], "--json", PR_FIELDS))
    update_pr(name, state, pr)


def open_submission(args, path, state, repo):
    paths = sorted({str(PurePosixPath(path)) for path in args.paths})
    for selected in paths:
        parts = PurePosixPath(selected)
        if parts.is_absolute() or ".." in parts.parts or selected == ".":
            raise ValueError(f"expected a repository-relative file or directory: {selected!r}")

    source_branch = state["source_branch"] if state else default_branch()
    source, = fetch(source_branch)
    base = commit_at(args.base)
    if state:
        if (state["base"], state["paths"]) != (base, paths):
            raise ValueError("submission already uses a different base or path selection")
    else:
        head = submission_commit(base, source, paths, f"Homework submission: {args.name}")
        if head is None:
            raise ValueError(f"the selected paths have no changes between {base} and {source_branch}")
        state = {
            "repo": repo,
            "base": base,
            "base_branch": f"submission/{args.name}/base",
            "head_branch": f"submission/{args.name}/head",
            "initial_head": head,
            "paths": paths,
            "source_branch": source_branch,
        }
        save_state(path, state)

    if "url" not in state:
        state["url"] = publish(args.name, state)
        save_state(path, state)
    else:
        refresh_pr(args.name, state)
    print(state["url"])


def sync_submission(name, state):
    if not state or "url" not in state:
        raise ValueError(f"finish opening the submission with 'open {name}' first")

    base_branch = state["base_branch"]
    head_branch = state["head_branch"]
    url = state["url"]
    refresh_pr(name, state)

    source, base, head = fetch(state["source_branch"], base_branch, head_branch)
    if base != state["base"]:
        raise ValueError(f"the submission base branch has moved: {base_branch}")

    updated = submission_commit(head, source, state["paths"], f"Sync homework submission: {name}")
    if updated is None:
        print("Submission is up to date.")
    else:
        git("push", "--no-follow-tags", "origin", f"{updated}:{branch_ref(head_branch)}")
    print(url)


def main():
    parser = argparse.ArgumentParser(prog="./scripts/hw/pr.sh")
    commands = parser.add_subparsers(dest="command", required=True)
    open_command = commands.add_parser("open", help="create and open a homework PR")
    open_command.add_argument("name", help="submission name, for example hw1")
    open_command.add_argument("--base", required=True, help="Git revision before the homework")
    open_command.add_argument("paths", nargs="+", help="repository-relative files or directories")
    sync_command = commands.add_parser("sync", help="update a homework PR from origin's default branch")
    sync_command.add_argument("name", help="submission name used with open")
    args = parser.parse_args()
    if not re.fullmatch(r"[a-zA-Z0-9][a-zA-Z0-9_-]*", args.name):
        parser.error("submission name must start with a letter or digit and use letters, digits, _ or -")

    try:
        if not shutil.which("gh"):
            raise ValueError("GitHub CLI (gh) is required")
        os.chdir(git("rev-parse", "--show-toplevel"))
        repo = github_repo()
        git_dir = Path(git("rev-parse", "--path-format=absolute", "--git-common-dir"))
        state_dir = git_dir / "homework-pr"
        state_dir.mkdir(exist_ok=True)
        path = state_dir / f"{args.name}.json"
        with path.with_suffix(".lock").open("a") as lock:
            fcntl.flock(lock, fcntl.LOCK_EX)
            state = json.loads(path.read_text(encoding="utf-8")) if path.exists() else None
            if state and state["repo"] != repo:
                raise ValueError("submission belongs to a different GitHub repository")
            if args.command == "open":
                open_submission(args, path, state, repo)
            else:
                sync_submission(args.name, state)
        return 0
    except subprocess.CalledProcessError as error:
        print(f"pr.sh: {error.cmd[0]} failed with exit status {error.returncode}", file=sys.stderr)
    except (OSError, ValueError) as error:
        print(f"pr.sh: {error}", file=sys.stderr)
    return 1


if __name__ == "__main__":
    sys.exit(main())
