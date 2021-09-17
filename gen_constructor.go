package gostructhelper

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

type constructorGenerator struct {
	typeName string
}

func NewConstructorGenerator(name string) *constructorGenerator {
	return &constructorGenerator{typeName: name}
}

func (g *constructorGenerator) Name() string {
	return "constructor"
}

func (g *constructorGenerator) Accept(c *astutil.Cursor) bool {
	spec, ok := filterTypeName(c.Node(), g.typeName)
	if !ok {
		return true
	}

	decl := getConstructor(spec)
	c.InsertAfter(decl)
	return false
}

func getConstructor(typSpec *ast.TypeSpec) *ast.FuncDecl {
	typStruct := typSpec.Type.(*ast.StructType)
	typIdent := &ast.Ident{Name: typSpec.Name.Name}

	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: genConstructorName(typSpec),
		},
		Type: &ast.FuncType{
			Params: genConstructorArgs(typStruct),
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: typIdent,
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: typIdent,
								Elts: genStructVals(typStruct),
							},
						},
					},
				},
			},
		},
	}
}

func genConstructorArgs(typStruct *ast.StructType) *ast.FieldList {
	ret := &ast.FieldList{
		List: []*ast.Field{},
	}

	for _, field := range typStruct.Fields.List {
		if len(field.Names) < 1 { // embedded field
			continue
		}
		fieldName := field.Names[0].Name

		ret.List = append(ret.List, &ast.Field{
			Names: []*ast.Ident{{Name: toArgName(fieldName)}},
			Type:  field.Type,
		})
	}
	return ret
}

func genStructVals(typStruct *ast.StructType) []ast.Expr {
	ret := []ast.Expr{}

	for _, field := range typStruct.Fields.List {
		if len(field.Names) < 1 { // embedded field
			continue
		}
		fieldName := field.Names[0].Name

		ret = append(ret, &ast.KeyValueExpr{
			Key:   &ast.Ident{Name: fieldName},
			Value: &ast.Ident{Name: toArgName(fieldName)},
		})
	}

	return ret
}
