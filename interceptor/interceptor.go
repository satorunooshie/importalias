// Package interceptor provides reusable wrappers for analysis.Analyzer values.
package interceptor

import (
	"go/ast"
	"go/token"
	"slices"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Func wraps an Analyzer and returns the wrapped Analyzer.
type Func func(*analysis.Analyzer) *analysis.Analyzer

// With applies fs in order and returns the resulting Analyzer.
func With(a *analysis.Analyzer, fs ...Func) *analysis.Analyzer {
	for _, f := range fs {
		a = f(a)
	}
	return a
}

// SkipGeneratedFile drops diagnostics in generated files.
func SkipGeneratedFile(a *analysis.Analyzer) *analysis.Analyzer {
	return wrapReport(a, func(pass *analysis.Pass, diagnostic analysis.Diagnostic) bool {
		file := fileAt(pass, diagnostic.Pos)
		return file != nil && ast.IsGenerated(file)
	})
}

// SkipByDirective drops diagnostics suppressed by //<tool>:ignore <analyzer>
// directives. A directive before the package clause suppresses the file;
// one on the same or previous line suppresses diagnostics on that line.
func SkipByDirective(tool string) Func {
	return func(a *analysis.Analyzer) *analysis.Analyzer {
		return wrapReport(a, func(pass *analysis.Pass, diagnostic analysis.Diagnostic) bool {
			file := fileAt(pass, diagnostic.Pos)
			if file == nil {
				return false
			}
			if hasFileDirective(file, tool, a.Name) {
				return true
			}
			line := pass.Fset.Position(diagnostic.Pos).Line
			for _, group := range file.Comments {
				for _, comment := range group.List {
					commentLine := pass.Fset.Position(comment.Pos()).Line
					if (commentLine == line || commentLine == line-1) && matches(comment, tool, a.Name) {
						return true
					}
				}
			}
			return false
		})
	}
}

func wrapReport(a *analysis.Analyzer, skip func(*analysis.Pass, analysis.Diagnostic) bool) *analysis.Analyzer {
	originalRun := a.Run
	a.Run = func(pass *analysis.Pass) (any, error) {
		originalReport := pass.Report
		pass.Report = func(diagnostic analysis.Diagnostic) {
			if !skip(pass, diagnostic) {
				originalReport(diagnostic)
			}
		}
		return originalRun(pass)
	}
	return a
}

func fileAt(pass *analysis.Pass, pos token.Pos) *ast.File {
	tokenFile := pass.Fset.File(pos)
	if tokenFile == nil {
		return nil
	}
	for _, file := range pass.Files {
		if pass.Fset.File(file.Pos()) == tokenFile {
			return file
		}
	}
	return nil
}

func hasFileDirective(file *ast.File, tool, analyzerName string) bool {
	for _, group := range file.Comments {
		if group.End() > file.Package {
			break
		}
		for _, comment := range group.List {
			if matches(comment, tool, analyzerName) {
				return true
			}
		}
	}
	return false
}

func matches(comment *ast.Comment, tool, analyzerName string) bool {
	directive, ok := ast.ParseDirective(comment.Pos(), comment.Text)
	return ok && directive.Tool == tool && directive.Name == "ignore" && slices.Contains(strings.Fields(directive.Args), analyzerName)
}
