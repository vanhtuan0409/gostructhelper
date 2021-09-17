package gostructhelper

import (
	"fmt"
	"go/ast"
	"unicode"
)

func genConstructorName(typSpec *ast.TypeSpec) string {
	return fmt.Sprintf("New%s", typSpec.Name.Name)
}

func toArgName(name string) string {
	rs := []rune(name)
	if len(rs) < 1 {
		return string(rs)
	}
	rs[0] = unicode.ToLower(rs[0])
	return string(rs)
}
