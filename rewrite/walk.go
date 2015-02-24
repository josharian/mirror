// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
)

// Helper functions for common node lists. They may be empty.

func walkIdentList(list []*ast.Ident, buf *bytes.Buffer) {
	for _, x := range list {
		Walk(x, buf)
	}
}

func walkExprList(list []ast.Expr, buf *bytes.Buffer) {
	for _, x := range list {
		Walk(x, buf)
	}
}

func walkStmtList(list []ast.Stmt, buf *bytes.Buffer) {
	for _, x := range list {
		Walk(x, buf)
	}
}

func walkDeclList(list []ast.Decl, buf *bytes.Buffer) {
	for _, x := range list {
		Walk(x, buf)
	}
}

var binaryOps = map[token.Token]string{
	token.ADD: "Add",
	token.SUB: "Sub",
}

func Walk(node ast.Node, buf *bytes.Buffer) {
	// TODO: Do work on node here?

	// walk children
	// (the order of the cases matches the order
	// of the corresponding node types in ast.go)
	switch n := node.(type) {
	// Comments and fields
	case *ast.Comment:
		// nothing to do

	case *ast.CommentGroup:
		// for _, c := range n.List {
		// 	Walk(c, buf)
		// }

	case *ast.Field:
		if n.Doc != nil {
			Walk(n.Doc, buf)
		}
		walkIdentList(n.Names, buf)
		Walk(n.Type, buf)
		if n.Tag != nil {
			Walk(n.Tag, buf)
		}
		if n.Comment != nil {
			Walk(n.Comment, buf)
		}

	case *ast.FieldList:
		for _, f := range n.List {
			Walk(f, buf)
		}

	// Expressions
	case *ast.BadExpr, *ast.Ident, *ast.BasicLit:
		// nothing to do

	case *ast.Ellipsis:
		if n.Elt != nil {
			Walk(n.Elt, buf)
		}

	case *ast.FuncLit:
		Walk(n.Type, buf)
		Walk(n.Body, buf)

	case *ast.CompositeLit:
		if n.Type != nil {
			Walk(n.Type, buf)
		}
		walkExprList(n.Elts, buf)

	case *ast.ParenExpr:
		Walk(n.X, buf)

	case *ast.SelectorExpr:
		Walk(n.X, buf)
		Walk(n.Sel, buf)

	case *ast.IndexExpr:
		Walk(n.X, buf)
		Walk(n.Index, buf)

	case *ast.SliceExpr:
		Walk(n.X, buf)
		if n.Low != nil {
			Walk(n.Low, buf)
		}
		if n.High != nil {
			Walk(n.High, buf)
		}
		if n.Max != nil {
			Walk(n.Max, buf)
		}

	case *ast.TypeAssertExpr:
		Walk(n.X, buf)
		if n.Type != nil {
			Walk(n.Type, buf)
		}

	case *ast.CallExpr:
		Walk(n.Fun, buf)
		walkExprList(n.Args, buf)

	case *ast.StarExpr:
		Walk(n.X, buf)

	case *ast.UnaryExpr:
		Walk(n.X, buf)

	case *ast.BinaryExpr:
		buf.WriteString("mirror.")
		if _, ok := binOp[n.Op]; !ok {
			panic(n.Op)
		}
		buf.WriteString(binOp[n.Op])
		buf.WriteString("(")
		Walk(n.X, buf)
		buf.WriteString(",")
		Walk(n.Y, buf)
		buf.WriteString(")")

	case *ast.KeyValueExpr:
		Walk(n.Key, buf)
		Walk(n.Value, buf)

	// Types
	case *ast.ArrayType:
		if n.Len != nil {
			Walk(n.Len, buf)
		}
		Walk(n.Elt, buf)

	case *ast.StructType:
		Walk(n.Fields, buf)

	case *ast.FuncType:
		if n.Params != nil {
			Walk(n.Params, buf)
		}
		if n.Results != nil {
			Walk(n.Results, buf)
		}

	case *ast.InterfaceType:
		Walk(n.Methods, buf)

	case *ast.MapType:
		Walk(n.Key, buf)
		Walk(n.Value, buf)

	case *ast.ChanType:
		Walk(n.Value, buf)

	// Statements
	case *ast.BadStmt:
		// nothing to do

	case *ast.DeclStmt:
		Walk(n.Decl, buf)

	case *ast.EmptyStmt:
		// nothing to do

	case *ast.LabeledStmt:
		Walk(n.Label, buf)
		Walk(n.Stmt, buf)

	case *ast.ExprStmt:
		Walk(n.X, buf)

	case *ast.SendStmt:
		Walk(n.Chan, buf)
		Walk(n.Value, buf)

	case *ast.IncDecStmt:
		Walk(n.X, buf)

	case *ast.AssignStmt:
		walkExprList(n.Lhs, buf)
		walkExprList(n.Rhs, buf)

	case *ast.GoStmt:
		Walk(n.Call, buf)

	case *ast.DeferStmt:
		Walk(n.Call, buf)

	case *ast.ReturnStmt:
		walkExprList(n.Results, buf)

	case *ast.BranchStmt:
		if n.Label != nil {
			Walk(n.Label, buf)
		}

	case *ast.BlockStmt:
		walkStmtList(n.List, buf)

	case *ast.IfStmt:
		if n.Init != nil {
			Walk(n.Init, buf)
		}
		Walk(n.Cond, buf)
		Walk(n.Body, buf)
		if n.Else != nil {
			Walk(n.Else, buf)
		}

	case *ast.CaseClause:
		walkExprList(n.List, buf)
		walkStmtList(n.Body, buf)

	case *ast.SwitchStmt:
		if n.Init != nil {
			Walk(n.Init, buf)
		}
		if n.Tag != nil {
			Walk(n.Tag, buf)
		}
		Walk(n.Body, buf)

	case *ast.TypeSwitchStmt:
		if n.Init != nil {
			Walk(n.Init, buf)
		}
		Walk(n.Assign, buf)
		Walk(n.Body, buf)

	case *ast.CommClause:
		if n.Comm != nil {
			Walk(n.Comm, buf)
		}
		walkStmtList(n.Body, buf)

	case *ast.SelectStmt:
		Walk(n.Body, buf)

	case *ast.ForStmt:
		if n.Init != nil {
			Walk(n.Init, buf)
		}
		if n.Cond != nil {
			Walk(n.Cond, buf)
		}
		if n.Post != nil {
			Walk(n.Post, buf)
		}
		Walk(n.Body, buf)

	case *ast.RangeStmt:
		if n.Key != nil {
			Walk(n.Key, buf)
		}
		if n.Value != nil {
			Walk(n.Value, buf)
		}
		Walk(n.X, buf)
		Walk(n.Body, buf)

	// Declarations
	case *ast.ImportSpec:
		if n.Doc != nil {
			Walk(n.Doc, buf)
		}
		if n.Name != nil {
			Walk(n.Name, buf)
		}
		Walk(n.Path, buf)
		if n.Comment != nil {
			Walk(n.Comment, buf)
		}

	case *ast.ValueSpec:
		if n.Doc != nil {
			Walk(n.Doc, buf)
		}
		walkIdentList(n.Names, buf)
		if n.Type != nil {
			Walk(n.Type, buf)
		}
		walkExprList(n.Values, buf)
		if n.Comment != nil {
			Walk(n.Comment, buf)
		}

	case *ast.TypeSpec:
		if n.Doc != nil {
			Walk(n.Doc, buf)
		}
		Walk(n.Name, buf)
		Walk(n.Type, buf)
		if n.Comment != nil {
			Walk(n.Comment, buf)
		}

	case *ast.BadDecl:
		// nothing to do

	case *ast.GenDecl:
		if n.Doc != nil {
			Walk(n.Doc, buf)
		}
		for _, s := range n.Specs {
			Walk(s, buf)
		}

	case *ast.FuncDecl:
		if n.Doc != nil {
			Walk(n.Doc, buf)
		}
		if n.Recv != nil {
			Walk(n.Recv, buf)
		}
		Walk(n.Name, buf)
		Walk(n.Type, buf)
		if n.Body != nil {
			Walk(n.Body, buf)
		}

	// Files and packages
	case *ast.File:
		if n.Doc != nil {
			Walk(n.Doc, buf)
		}
		Walk(n.Name, buf)
		walkDeclList(n.Decls, buf)
		// don't walk n.Comments - they have been
		// visited already through the individual
		// nodes

	case *ast.Package:
		for _, f := range n.Files {
			Walk(f, buf)
		}

	default:
		fmt.Printf("Walk: unexpected node type %T", n)
		panic("Walk")
	}

	return
}

func main() {
	fset, files, err := parsePackage("github.com/josharian/impl")
	if err != nil {
		log.Fatal(err)
	}
	var buf bytes.Buffer
	for _, f := range files {
		fmt.Println("---")
		Walk(f, &buf)
		fmt.Println(buf.String())
		buf.Reset()
	}
	_ = fset
}

func parsePackage(path string) (*token.FileSet, []*ast.File, error) {
	bpkg, err := build.Import(path, "", 0)
	if err != nil {
		return nil, nil, err
	}

	fset := token.NewFileSet()
	var files []*ast.File
	for _, file := range bpkg.GoFiles {
		f, err := parser.ParseFile(fset, filepath.Join(bpkg.Dir, file), nil, 0)
		if err != nil {
			return nil, nil, err
		}
		files = append(files, f)
	}
	return fset, files, nil
}

var binOp = map[token.Token]string{
	token.ADD:     "Add",
	token.AND:     "And",
	token.AND_NOT: "AndNot",
	token.EQL:     "Eql",
	token.GEQ:     "Geq",
	token.GTR:     "Gtr",
	token.LAND:    "Land",
	token.LEQ:     "Leq",
	token.LOR:     "Lor",
	token.LSS:     "Lss",
	token.MUL:     "Mul",
	token.NEQ:     "Neq",
	token.OR:      "Or",
	token.QUO:     "Quo",
	token.REM:     "Rem",
	token.SHL:     "Shl",
	token.SHR:     "Shr",
	token.SUB:     "Sub",
	token.XOR:     "Xor",
}
