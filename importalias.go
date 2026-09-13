// Package importalias reports unnecessary import aliases and import names
// containing underscores.
package importalias

import (
	"fmt"
	"go/ast"
	"go/types"
	"path"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/satorunooshie/importalias/interceptor"
	"golang.org/x/tools/go/analysis"
)

var analyzer = &analysis.Analyzer{
	Name: "importalias",
	Doc:  "reports unnecessary import aliases and import names containing underscores",
	Run:  run,
}

// Analyzer is the importalias analyzer.
var Analyzer = interceptor.With(analyzer, interceptor.SkipGeneratedFile, interceptor.SkipByDirective("importalias"))

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		uses := packageUses(pass, file)
		for _, spec := range file.Imports {
			alias := ""
			if spec.Name != nil {
				alias = spec.Name.Name
			}
			if alias == "." || alias == "_" {
				continue
			}

			pkgName := importedPackageName(pass, spec)
			if pkgName == nil || pkgName.Imported() == nil {
				continue
			}
			actual := pkgName.Imported().Name()
			if alias != "" && strings.Contains(alias, "_") {
				pass.Reportf(spec.Name.Pos(), "import alias %q should not contain underscore", alias)
				continue
			}
			if alias == "" {
				if strings.Contains(actual, "_") {
					pass.Reportf(spec.Path.Pos(), "imported package name %q contains underscore; use an alias", actual)
				}
				continue
			}
			if alias != actual && strings.Contains(actual, "_") {
				continue
			}

			importPath, ok := unquotedImportPath(spec)
			if !ok || importPathImpliedName(importPath) != actual {
				continue
			}
			if alias != actual && !canUsePackageName(pass, file, spec, actual) {
				continue
			}

			reportUnnecessaryAlias(pass, uses, pkgName, spec, alias, actual)
		}
	}
	return nil, nil
}

func unquotedImportPath(spec *ast.ImportSpec) (string, bool) {
	importPath, err := strconv.Unquote(spec.Path.Value)
	return importPath, err == nil
}

func importPathImpliedName(importPath string) string {
	element := path.Base(importPath)
	if isMajorVersionElement(element) && path.Dir(importPath) != "." {
		element = path.Base(path.Dir(importPath))
	}
	return identifierPrefix(strings.TrimPrefix(element, "go-"))
}

func isMajorVersionElement(element string) bool {
	if len(element) < 2 || element[0] != 'v' {
		return false
	}
	_, err := strconv.Atoi(element[1:])
	return err == nil
}

func identifierPrefix(value string) string {
	for index, r := range value {
		if !isIdentifierRune(r) {
			return value[:index]
		}
	}
	return value
}

func isIdentifierRune(r rune) bool {
	return r == '_' || 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || '0' <= r && r <= '9' ||
		r >= utf8.RuneSelf && (unicode.IsLetter(r) || unicode.IsDigit(r))
}

func reportUnnecessaryAlias(pass *analysis.Pass, uses map[*types.PkgName][]*ast.Ident, pkgName *types.PkgName, spec *ast.ImportSpec, alias, actual string) {
	edits := []analysis.TextEdit{{Pos: spec.Name.Pos(), End: spec.Path.Pos()}}
	if alias != actual {
		for _, ident := range uses[pkgName] {
			edits = append(edits, analysis.TextEdit{Pos: ident.Pos(), End: ident.End(), NewText: []byte(actual)})
		}
	}

	pass.Report(analysis.Diagnostic{
		Pos: spec.Name.Pos(), End: spec.Name.End(),
		Message: fmt.Sprintf("import alias %q is unnecessary; package name is %q", alias, actual),
		SuggestedFixes: []analysis.SuggestedFix{{
			Message:   "Remove import alias",
			TextEdits: edits,
		}},
	})
}

func packageUses(pass *analysis.Pass, file *ast.File) map[*types.PkgName][]*ast.Ident {
	fileInfo := pass.Fset.File(file.Pos())
	uses := make(map[*types.PkgName][]*ast.Ident)
	for ident, object := range pass.TypesInfo.Uses {
		if ident == nil || pass.Fset.File(ident.Pos()) != fileInfo {
			continue
		}
		if pkgName, ok := object.(*types.PkgName); ok {
			uses[pkgName] = append(uses[pkgName], ident)
		}
	}
	for _, idents := range uses {
		sort.Slice(idents, func(i, j int) bool { return idents[i].Pos() < idents[j].Pos() })
	}
	return uses
}

func canUsePackageName(pass *analysis.Pass, file *ast.File, target *ast.ImportSpec, name string) bool {
	for _, spec := range file.Imports {
		if spec == target {
			continue
		}
		if importName(pass, spec) == name {
			return false
		}
		if pkgName := importedPackageName(pass, spec); pkgName != nil && pkgName.Imported() != nil && pkgName.Imported().Name() == name {
			return false
		}
	}

	fileInfo := pass.Fset.File(file.Pos())
	for ident, object := range pass.TypesInfo.Defs {
		if ident == nil || object == nil || ident.Name != name || pass.Fset.File(ident.Pos()) != fileInfo {
			continue
		}
		if variable, ok := object.(*types.Var); ok && variable.IsField() {
			continue
		}
		return false
	}
	return true
}

func importName(pass *analysis.Pass, spec *ast.ImportSpec) string {
	if spec.Name != nil {
		return spec.Name.Name
	}
	if pkgName := importedPackageName(pass, spec); pkgName != nil && pkgName.Imported() != nil {
		return pkgName.Imported().Name()
	}
	return ""
}

func importedPackageName(pass *analysis.Pass, spec *ast.ImportSpec) *types.PkgName {
	if spec.Name != nil {
		if object, ok := pass.TypesInfo.Defs[spec.Name].(*types.PkgName); ok {
			return object
		}
	}
	if object, ok := pass.TypesInfo.Implicits[spec].(*types.PkgName); ok {
		return object
	}
	return nil
}
