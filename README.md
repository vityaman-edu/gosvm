# GoSVM

Yet Another Go Symbolic Virtual Machine based on Alexey Menshutin course.

> Time forks perpetually toward innumerable futures.
>
> -- Jorge Luis Borges, The Garden of Forking Paths (1941)

<!-- markdownlint-disable MD033 -->
<img src="./docs/asset/mascot.png" alt="Mascot" width="600">
<!-- markdownlint-restore -->

## Prepare

Enter the Nix development shell.

```sh
nix develop
```

I recommend using Zed editor.

```sh
zed .
```

Install precommit hooks.

```sh
pre-commit install
```

## Test

```sh
pre-commit run --all-files
```

## Build

```sh
go build -o ./gosvm ./cmd/gosvm
```

## Run

```sh
# Dry run
./gosvm -function simpleLoop < internal/examples/test_functions.go

# See interpretation
./gosvm -function simpleLoop -print-smt < internal/examples/test_functions.go
```

## Learning

See the [course/README.md](course/README.md).

## References

- <https://github.com/CaelmBleidd/symbolic_execution_2026>
