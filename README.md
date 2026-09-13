# importalias

`importalias` is a conservative [`go/analysis`](https://pkg.go.dev/golang.org/x/tools/go/analysis) linter for import names.

It reports unnecessary aliases when the import path and declared package name agree, and reports package or alias names containing underscores. It skips generated files, blank/dot imports, path-required aliases, and unsafe selector rewrites. Safe diagnostics include suggested fixes.

## Install

```sh
go install github.com/satorunooshie/importalias/cmd/importalias@latest
```

## Run

```sh
importalias ./...
```

The command uses the standard analysis driver and supports `-fix`, `-diff`, `-json`, and editor integrations.

## Use with `go vet`

Build the analyzer once and pass it to `go vet` as a vettool:

```sh
go build -o ./bin/importalias ./cmd/importalias
go vet -vettool=./bin/importalias ./...
```

This makes `importalias` usable in existing `go vet` and CI workflows.

## Ignoring a diagnostic

Use a Go directive to suppress one analyzer diagnostic on a line, or all
diagnostics in a file:

```go
//importalias:ignore importalias
package example
```

The directive follows the standard `ast.ParseDirective` syntax and is handled
by the reusable `interceptor` package.

See [the analyzer specification](docs/analyzer-spec.md) for the complete rule
set and suggested-fix safety conditions.
