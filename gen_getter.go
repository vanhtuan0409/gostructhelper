package gostructhelper

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"os"

	"golang.org/x/tools/go/ast/astutil"
)

type getterGenerator struct {
	typeName string
	fset     *token.FileSet
}

func NewGetterGenerator(name string, fset *token.FileSet) *getterGenerator {
	return &getterGenerator{
		typeName: name,
		fset:     fset,
	}
}

func (g *getterGenerator) Name() string {
	return "getter"
}

func (g *getterGenerator) Accept(c *astutil.Cursor) bool {
	spec, ok := filterTypeName(c.Node(), g.typeName)
	if !ok {
		return true
	}

	getters := g.getGetters(spec)
	for _, getter := range getters {
		c.InsertAfter(getter)
	}

	return false
}

func (g *getterGenerator) getGetters(typSpec *ast.TypeSpec) []*ast.FuncDecl {
	typStruct := typSpec.Type.(*ast.StructType)
	getters := []*ast.FuncDecl{}
	fieldLen := len(typStruct.Fields.List)
	for i := fieldLen - 1; i >= 0; i-- {
		field := typStruct.Fields.List[i]
		getterFn := g.getGetterForField(typSpec, field)
		if getterFn != nil {
			getters = append(getters, getterFn)
		}
	}
	return getters
}

func (g *getterGenerator) getGetterForField(typSpect *ast.TypeSpec, field *ast.Field) *ast.FuncDecl {
	recvName := genGetterReceiver(typSpect.Name.Name)
	getterName := genGetterName(field)
	if getterName == "" {
		fmt.Fprintf(os.Stderr, "Cannot generate getter name. Probably because field is an embedded field\n")
		return nil
	}
	var fieldTyp bytes.Buffer
	if err := printer.Fprint(&fieldTyp, g.fset, field.Type); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot generate field type expr. ERR: %+v\n", err)
		return nil
	}

	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: recvName,
							Obj: &ast.Object{
								Kind: ast.Var,
								Name: recvName,
							},
						},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{
							Name: typSpect.Name.Name,
						},
					},
				},
			},
		},
		Name: &ast.Ident{
			Name: getterName,
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{List: []*ast.Field{}},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.Ident{
							Name: fieldTyp.String(),
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.SelectorExpr{
							X: &ast.Ident{
								Name: recvName,
							},
							Sel: &ast.Ident{
								Name: field.Names[0].Name,
							},
						},
					},
				},
			},
		},
	}
}
