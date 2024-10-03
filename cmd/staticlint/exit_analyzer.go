package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// OsExitChecker is an analyzer that checks for os.Exit calls in the main function of the main package.
var OsExitChecker = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "check for os.Exit calls in the main function of the main package",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Name.Name == "main" && funcDecl.Recv == nil {
				ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
					if callExpr, ok := n.(*ast.CallExpr); ok {
						if fun, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
							if ident, ok := fun.X.(*ast.Ident); ok {
								if ident.Name == "os" && fun.Sel.Name == "Exit" {
									pass.Reportf(callExpr.Pos(), "direct call to os.Exit in main function of main package is prohibited")
								}
							}
						}
					}
					return true
				})
			}
		}
	}
	return nil, nil
}
