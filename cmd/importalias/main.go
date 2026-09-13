package main

import (
	"github.com/satorunooshie/importalias"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(importalias.Analyzer)
}
