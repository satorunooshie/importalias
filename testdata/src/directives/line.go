package directives

import (
	json "encoding/json"                     // want `import alias "json" is unnecessary; package name is "json"`
	renamed "example.com/importalias/target" //importalias:ignore importalias
)

var (
	_ = renamed.Value
	_ = json.Valid
)
