package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"

	"golang.org/x/tools/go/ast/astutil"
)

func main() {
	tokenSet := token.NewFileSet()
	fPath := "/home/tuan/Workspaces/go-test/main.go"

	astTree, err := parser.ParseFile(tokenSet, fPath, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	astutil.Apply(astTree, nil, func(c *astutil.Cursor) bool {
		node := c.Node()
		typSpec, ok := isTypeSpec(node, "Yolo")
		if !ok {
			return true
		}

		cstrDecl := getConstructor(typSpec)
		c.InsertAfter(cstrDecl)
		return false
	})

	printer.Fprint(os.Stdout, tokenSet, astTree)
}

func isTypeSpec(node ast.Node, name string) (*ast.TypeSpec, bool) {
	// Should stop add top-level declarations
	// This allow to insert type funcs after type decl
	genDecl, ok := node.(*ast.GenDecl)
	if !ok {
		return nil, false
	}
	if len(genDecl.Specs) == 0 {
		return nil, false
	}
	spec := genDecl.Specs[0]

	typSpec, ok := spec.(*ast.TypeSpec)
	if !ok {
		return nil, false
	}

	if typSpec.Name.Name != name {
		return nil, false
	}
	return typSpec, true
}

func getConstructor(typSpec *ast.TypeSpec) *ast.FuncDecl {
	typStruct := typSpec.Type.(*ast.StructType)
	fieldLen := len(typStruct.Fields.List)

	funcArgs := make([]*ast.Field, fieldLen)
	structVal := make([]ast.Expr, fieldLen)
	for i := 0; i < fieldLen; i++ {
		field := typStruct.Fields.List[i]
		fieldName := field.Names[0].Name

		funcArgs[i] = &ast.Field{
			Names: []*ast.Ident{{Name: fieldName}},
			Type:  field.Type,
		}
		structVal[i] = &ast.KeyValueExpr{
			Key:   &ast.Ident{Name: fieldName},
			Value: &ast.Ident{Name: fieldName},
		}
	}

	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: fmt.Sprintf("New%s", typSpec.Name.Name),
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: funcArgs,
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.Ident{
								Name: typSpec.Name.Name,
							},
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
								Type: &ast.Ident{
									Name: typSpec.Name.Name,
								},
								Elts: structVal,
							},
						},
					},
				},
			},
		},
	}
}
