package a

import (
	json "encoding/json" // want `import alias "json" is unnecessary; package name is "json"`
	"example.com/importalias/underscore" // want `imported package name "underscore_pkg" contains underscore; use an alias`
	renamed "example.com/importalias/target" // want `import alias "renamed" is unnecessary; package name is "target"`
)

var (
	_ = json.Valid
	_ = underscore_pkg.Value
	_ = renamed.Value
)
