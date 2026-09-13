package a

import (
	json "encoding/json" // want `import alias "json" is unnecessary; package name is "json"`
	_ "example.com/importalias/blank"
	. "example.com/importalias/dot"
	schemaapi "example.com/importalias/schema/v2" // want `import alias "schemaapi" is unnecessary; package name is "schema"`
	renamed "example.com/importalias/target"      // want `import alias "renamed" is unnecessary; package name is "target"`
	"example.com/importalias/underscore"          // want `imported package name "underscore_pkg" contains underscore; use an alias`
)

var (
	_ = json.Valid
	_ = underscore_pkg.Value
	_ = schemaapi.Value
	_ = renamed.Value
	_ = DotValue
)
