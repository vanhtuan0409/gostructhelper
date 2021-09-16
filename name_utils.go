package gostructhelper

import (
	"fmt"
	"go/ast"
)

func genConstructorName(typSpec *ast.TypeSpec) string {
	return fmt.Sprintf("New%s", typSpec.Name.Name)
}
