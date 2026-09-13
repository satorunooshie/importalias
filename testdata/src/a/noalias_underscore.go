package a

import "example.com/importalias/underscore" // want `imported package name "underscore_pkg" contains underscore; use an alias`

var _ = underscore_pkg.Value
