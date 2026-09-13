package a

import json_alias "encoding/json" // want `import alias "json_alias" should not contain underscore`

var _ = json_alias.Valid
