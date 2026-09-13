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
