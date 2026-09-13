package a

import (
	custompkg "example.com/importalias/customname"
	actual "example.com/importalias/legacy"
)

var (
	_ = custompkg.Value
	_ = actual.Value
)
