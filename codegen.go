package gostructhelper

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"

	"golang.org/x/tools/go/ast/astutil"
)

func genInternal(reg *Registry, tree *ast.File, tokenSet *token.FileSet, out io.WriteSeeker) {
	astutil.Apply(tree, nil, func(c *astutil.Cursor) bool {
		for _, g := range reg.generators {
			if !reg.isDone(g) {
				reg.done[g.Name()] = !g.Accept(c)
			}
		}
		return !reg.isAllDone()
	})

	out.Seek(0, 0)
	printer.Fprint(out, tokenSet, tree)
}

func Gen(reg *Registry, s *source, out io.WriteSeeker) error {
	tokenSet := token.NewFileSet()
	tree, err := parser.ParseFile(tokenSet, s.path, s.r, parser.ParseComments)
	if err != nil {
		return err
	}
	genInternal(reg, tree, tokenSet, out)
	return nil
}
