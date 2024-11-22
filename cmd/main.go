package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// isTestFunction checks if a function is a test function.
func isTestFunction(funcDecl *ast.FuncDecl) bool {
	// Check if the function is exported and starts with "Test"
	if !funcDecl.Name.IsExported() || !strings.HasPrefix(funcDecl.Name.Name, "Test") {
		return false
	}

	// Ensure the function has exactly one parameter
	if funcDecl.Type.Params.NumFields() != 1 {
		return false
	}

	// Check the parameter type
	paramType := funcDecl.Type.Params.List[0].Type
	switch expr := paramType.(type) {
	case *ast.StarExpr: // Handle pointer types
		if sel, ok := expr.X.(*ast.SelectorExpr); ok {
			return sel.Sel.Name == "T" // Check if it matches *testing.T
		} else if ident, ok := expr.X.(*ast.Ident); ok {
			return ident.Name == "T" // Check if it matches T
		}
	}
	return false
}

// containsTRun checks if a function body contains any t.Run() calls.
// containsReflection checks if a function body contains specific reflection usage, like reflect.Indirect.
func containsReflection(body *ast.BlockStmt) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		if found {
			return false
		}
		// Check for calls to reflect methods like reflect.Indirect
		if callExpr, ok := n.(*ast.CallExpr); ok {
			if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := selExpr.X.(*ast.Ident); ok {
					if ident.Name == "reflect" && (selExpr.Sel.Name == "Indirect" || selExpr.Sel.Name == "ValueOf") {
						found = true
						return false
					}
				}
			}
		}
		return true
	})
	return found
}

// containsTRun checks if a function body contains any t.Run() calls.
func containsTRun(body *ast.BlockStmt) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		if found {
			return false
		}
		if exprStmt, ok := n.(*ast.ExprStmt); ok {
			if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
				if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
					if selExpr.Sel.Name == "Run" {
						if ident, ok := selExpr.X.(*ast.Ident); ok {
							if ident.Name == "t" {
								found = true
								return false
							}
						}
					}
				}
			}
		}
		return true
	})
	return found
}

// processFile parses a Go test file and checks for test functions without t.Run().
// It skips functions that use reflection but does not print them.
func processFile(path string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return err
	}

	for _, decl := range node.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok && isTestFunction(funcDecl) {
			// Skip functions using reflection silently
			if containsReflection(funcDecl.Body) {
				continue
			}
			// Get the line number of the function
			line := fset.Position(funcDecl.Pos()).Line
			// Print functions without t.Run() with line number
			if !containsTRun(funcDecl.Body) {
				fmt.Printf("File: %s:%d, Function: %s\n", path, line, funcDecl.Name.Name)
			}
		}
	}
	return nil
}

// walkDir recursively walks through directories to find *_test.go files.
func walkDir(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), "_test.go") {
			if err := processFile(path); err != nil {
				return err
			}
		}
		return nil
	})
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: main2 <directory>")
		os.Exit(1)
	}
	dir := os.Args[1]
	if err := walkDir(dir); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
