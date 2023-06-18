// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package inlinemethod defines an Analyzer that inlines simple methods.
package inlinemethod

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `check for simple methods that can be inlined

A method consisting of a single statement such as:
	// todo
will be removed, and its calls replaced with that single statement:
	// todo
`

var Analyzer = &analysis.Analyzer{
	Name:     "inlinemethod",
	Doc:      Doc,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		expr := n.(*ast.CallExpr)
		function, ok := expr.Fun.(*ast.Ident)
		if !ok || function.Obj == nil {
			return
		}
		functionData, ok := function.Obj.Decl.(*ast.FuncDecl)
		if !ok || functionData.Body == nil {
			return
		}
		if len(functionData.Body.List) == 1 {
			pass.Report(analysis.Diagnostic{
				Pos:     expr.Pos(),
				End:     expr.End(),
				Message: fmt.Sprintf("inline method: %s", function.Name),
				SuggestedFixes: []analysis.SuggestedFix{{
					Message: fmt.Sprintf("Inline method: '%s'", function.Name),
					TextEdits: []analysis.TextEdit{{
						Pos:     expr.Pos(),
						End:     expr.End(),
						NewText: suggestedFix(pass, expr, functionData),
					}},
				}},
			})
		}
	})
	return nil, nil
}

func suggestedFix(pass *analysis.Pass, functionCallExpr *ast.CallExpr, functionData *ast.FuncDecl) []byte {
	paramIndexMap := make(map[string]ast.Expr)
	for i, argExp := range functionData.Type.Params.List {
		paramIndexMap[argExp.Names[0].Name] = functionCallExpr.Args[i]
	}
	var resultNode ast.Node

	// Extract return expression
	newStmt := functionData.Body.List[0]
	returnStmt, ok := newStmt.(*ast.ReturnStmt)
	if ok {
		resultNode = returnStmt.Results[0]
	} else {
		resultNode = newStmt
	}

	// Replace identifiers with their values

	var b bytes.Buffer
	printer.Fprint(&b, pass.Fset, resultNode)
	return b.Bytes()
}
