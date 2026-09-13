package target

import renamed "example.com/importalias/target" // want `import alias "renamed" is unnecessary; package name is "target"`

var _ = renamed.Value
