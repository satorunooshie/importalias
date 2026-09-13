# Analyzer specification

`importalias` checks import names in non-generated Go files.

## Diagnostics

### Unnecessary aliases

An explicit alias is reported when all of the following are true:

- the import is not `_` or `.`,
- the imported package can be resolved,
- the package name matches the name implied by the import path, and
- removing or replacing the alias is safe in that file.

The import path name is taken from its final element. A semantic-major suffix
such as `/v2` is ignored, a `go-` prefix is removed, and the name is truncated
at the first non-identifier character.

When the alias equals the package name, the fix removes only the alias. When
the alias differs, the fix also rewrites package selectors, but only if no
import or non-field declaration in the file already uses that package name.

Aliases are not reported when the package declaration intentionally differs
from the import path, such as an import path ending in `customname` with a
`package custompkg` declaration.

### Underscores

The analyzer reports an explicit alias containing `_`. It also reports an
unaliased import whose declared package name contains `_`. These diagnostics do
not provide automatic fixes because the appropriate replacement is
project-specific.

An explicit alias without `_` is allowed to hide a package name containing `_`.

## Exclusions

Diagnostics are suppressed for generated files and for the following standard
directive forms:

```go
//importalias:ignore importalias
```

A directive before the `package` clause suppresses the whole file. A directive
on the same line as, or immediately before, a diagnostic suppresses that
diagnostic. Other tools, other analyzer names, and comments with whitespace
between `//` and the tool name are ignored.

Blank and dot imports are always ignored. The analyzer also skips unresolved
imports and never reports aliases that are required by the package declaration
or by a name collision.
