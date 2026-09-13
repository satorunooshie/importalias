package a

import (
	diskcache "example.com/importalias/disk/cache"
	memorycache "example.com/importalias/memory/cache"
)

var (
	_ = diskcache.Value
	_ = memorycache.Value
)
