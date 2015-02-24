// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	fset, files, err := parsePackage("github.com/josharian/mirror/sample")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		fmt.Println("---")
		if err := Fprint(os.Stdout, f); err != nil {
			log.Fatal(err)
		}
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

func Fprint(w io.Writer, n ast.Node) error {
	p := &printer{w: w}
	return p.printRoot(n)
}

type printErr error

type printer struct {
	w io.Writer
}

func (p *printer) print(s interface{}) {
	p.printf("%v", s)
}

func (p *printer) printf(format string, a ...interface{}) {
	_, err := fmt.Fprintf(p.w, format, a...)
	if err != nil {
		panic(err)
	}
}

func (p *printer) printRoot(n ast.Node) (err error) {
	defer func() {
		if perr := recover(); perr != nil {
			var ok bool
			if err, ok = perr.(printErr); !ok {
				panic(perr)
			}
		}
	}()
	p.printNode(n)
	return nil
}

func (p *printer) printNode(n ast.Node) {
	switch n := n.(type) {
	case *ast.File:
		p.printf("package %s\n", n.Name)
		p.printDeclList(n.Decls)
	case *ast.ValueSpec:
		p.printValueSpec(n)
	default:
		fmt.Fprintf(os.Stderr, "printNode not implemented for %T\n", n)
	}
}

func (p *printer) printDeclList(list []ast.Decl) {
	for _, x := range list {
		p.printDecl(x)
	}
}

func (p *printer) printGenDecl(n *ast.GenDecl) {
	switch n.Tok {
	case token.VAR:
		p.print("var ")
	default:
		fmt.Fprintf(os.Stderr, "printGenDecl not implemented for %s\n", n.Tok)
	}
	// TODO:
	//     Tok    token.Token   // IMPORT, CONST, TYPE, VAR
	//     Specs  []Spec
	for _, spec := range n.Specs {
		p.printSpec(spec)
	}
}

func (p *printer) printFuncDecl(n *ast.FuncDecl) {
	// TODO: Deal with complications around exported vs non-exported

	// TODO:
	// if n.Recv != nil {
	// 	Walk(v, n.Recv)
	// }
	// Walk(v, n.Name)
	// Walk(v, n.Type)
	if n.Body != nil {
		p.printBlockStmt(n.Body)
	}
}

func (p *printer) printSpec(n ast.Spec) {
	p.printNode(n)
}

func (p *printer) printValueSpec(n *ast.ValueSpec) {
	// TODO
	// Names   []*Ident      // value names (len(Names) > 0)
	// Type    Expr          // value type; or nil
	// Values  []Expr        // initial values; or nil
	for i := 0; i < len(n.Names); i++ {
		name := n.Names[i]
		if i > 0 {
			p.print(", ")
		}
		p.printIdent(name)
	}
	p.print(" ")
	p.printExpr(n.Type)
	// TODO: Values
	p.print("\n")
}

func (p *printer) printStmt(n ast.Stmt) {
	switch n := n.(type) {
	case *ast.AssignStmt:
		p.printAssignStmt(n)
	case *ast.DeclStmt:
		p.printDeclStmt(n)
	default:
		fmt.Fprintf(os.Stderr, "printStmt not implemented for %T\n", n)
	}
}

func (p *printer) printDeclStmt(n *ast.DeclStmt) {
	p.printDecl(n.Decl)
}

func (p *printer) printAssignStmt(n *ast.AssignStmt) {
	for _, expr := range n.Lhs {
		p.printExpr(expr)
	}
	p.print("=")
	for _, expr := range n.Rhs {
		p.printExpr(expr)
	}
	p.print("\n")
}

func (p *printer) printBlockStmt(n *ast.BlockStmt) {
	for _, stmt := range n.List {
		p.printStmt(stmt)
	}
}

func (p *printer) printIdent(n *ast.Ident) {
	p.print(n.String())
}

func (p *printer) printDecl(n ast.Decl) {
	switch n := n.(type) {
	case *ast.FuncDecl:
		p.printFuncDecl(n)
	case *ast.GenDecl:
		p.printGenDecl(n)
	default:
		fmt.Fprintf(os.Stderr, "printDecl not implemented for %T\n", n)
	}
}

func (p *printer) printExpr(n ast.Expr) {
	switch n := n.(type) {
	case *ast.BinaryExpr:
		p.printBinaryExpr(n)
	case *ast.Ident:
		p.printIdent(n)
	default:
		// TODO
		fmt.Fprintf(os.Stderr, "printExpr not implemented for %T\n", n)
	}
}

func (p *printer) printBinaryExpr(n *ast.BinaryExpr) {
	op, ok := binOp[n.Op]
	if !ok {
		panic(n.Op)
	}
	p.printf("mirror.%s(", op)
	p.printExpr(n.X)
	p.print(", ")
	p.printExpr(n.Y)
	p.print(")")
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
