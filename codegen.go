package gostructhelper

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"

	"golang.org/x/tools/go/ast/astutil"
)

func genInternal(reg *Registry, tree *ast.File, fs *token.FileSet, out io.WriteSeeker) {
	astutil.Apply(tree, nil, func(c *astutil.Cursor) bool {
		for _, g := range reg.generators {
			if !reg.isDone(g) {
				reg.done[g.Name()] = !g.Accept(c)
			}
		}
		return !reg.isAllDone()
	})

	out.Seek(0, 0)
	format.Node(out, fs, tree)
}

func Gen(reg *Registry, s *source, fs *token.FileSet, out io.WriteSeeker) error {
	tree, err := parser.ParseFile(fs, s.path, s.r, parser.ParseComments)
	if err != nil {
		return err
	}
	genInternal(reg, tree, fs, out)
	return nil
}
