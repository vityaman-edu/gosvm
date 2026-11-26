# GoSVM

Yet Another Go Symbolic Virtual Machine based on Alexey Menshutin course.

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
