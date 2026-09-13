package a

import schemaapi "example.com/importalias/schema/v2" // want `import alias "schemaapi" is unnecessary; package name is "schema"`

var _ = schemaapi.Value
