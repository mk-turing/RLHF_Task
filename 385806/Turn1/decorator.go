package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
)

func main() {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, "Turn1A.go", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	generateDecorators(fs, f)
}

func generateDecorators(fs *token.FileSet, f *ast.File) {
	for _, d := range f.Decls {
		if funcDecl, ok := d.(*ast.FuncDecl); ok {
			name := funcDecl.Name.Name
			decoratedName := fmt.Sprintf("decorated%s", name)
			// Generate wrapper function
			wrapper := generateWrapper(fs, funcDecl, decoratedName)

			// Add the wrapper function to the file
			f.Decls = append(f.Decls, wrapper)
		}
	}

	// Format the generated code
	var buf bytes.Buffer
	if err := format.Node(&buf, fs, f); err != nil {
		panic(err)
	}

	// Write the modified file back
	if err := ioutil.WriteFile("Turn1A.go", buf.Bytes(), 0644); err != nil {
		panic(err)
	}
}

func generateWrapper(fs *token.FileSet, funcDecl *ast.FuncDecl, decoratedName string) ast.Decl {
	// Create the wrapper function signature
	funcType := funcDecl.Type

	// Initialize the body as a *ast.BlockStmt
	wrapperFunc := &ast.FuncDecl{
		Name: &ast.Ident{
			Name: decoratedName,
		},
		Type: &ast.FuncType{
			Params:  funcType.Params,  // Directly access Params
			Results: funcType.Results, // Directly access Results
		},
		Body: &ast.BlockStmt{ // Correctly initialize as *ast.BlockStmt
			List: []ast.Stmt{}, // Initialize the statement list
		},
	}

	// Add the policy handling logic
	policyCall := generatePolicyCall(fs, funcDecl.Name.Name)
	wrapperFunc.Body.List = append(wrapperFunc.Body.List, policyCall)

	// Create the call to the original function
	originalFuncCall := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun:  &ast.Ident{Name: funcDecl.Name.Name}, // Call the original function
			Args: generateArguments(funcDecl),          // Pass the original function's arguments
		},
	}
	wrapperFunc.Body.List = append(wrapperFunc.Body.List, originalFuncCall)

	return wrapperFunc
}

// generateArguments generates the function arguments for the original function call
func generateArguments(funcDecl *ast.FuncDecl) []ast.Expr {
	args := []ast.Expr{}
	for _, param := range funcDecl.Type.Params.List {
		// Assuming that each parameter has an identifier
		arg := &ast.Ident{Name: param.Names[0].Name}
		args = append(args, arg)
	}
	return args
}

// generatePolicyCall generates the call to the error handling policy
func generatePolicyCall(fs *token.FileSet, funcName string) ast.Stmt {
	// Assume the policy is applied to the last result (use an underscore as a placeholder)
	lastResult := ast.NewIdent("_") // Placeholder, adjust based on function return type
	policyCall := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "LogPolicy"},
				Sel: &ast.Ident{Name: "Handle"},
			},
			Args: []ast.Expr{lastResult},
		},
	}

	return policyCall
}
