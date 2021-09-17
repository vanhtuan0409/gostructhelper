package gostructhelper

import (
	"fmt"
	"go/ast"
	"unicode"
)

func genConstructorName(typSpec *ast.TypeSpec) string {
	typName := typSpec.Name.Name
	return fmt.Sprintf("New%s", toExportedName(typName))
}

func genGetterReceiver(name string) string {
	rs := []rune(name)
	firstChar := rs[0]
	return string([]rune{unicode.ToLower(firstChar)})
}

func genGetterName(field *ast.Field) string {
	if len(field.Names) < 1 {
		return ""
	}
	fieldName := field.Names[0]
	if fieldName.IsExported() {
		return fmt.Sprintf("Get%s", fieldName.Name)
	}
	return toExportedName(fieldName.Name)
}

func toArgName(name string) string {
	rs := []rune(name)
	if len(rs) < 1 {
		return string(rs)
	}
	rs[0] = unicode.ToLower(rs[0])
	return string(rs)
}

func toExportedName(name string) string {
	rs := []rune(name)
	if len(rs) < 1 {
		return string(rs)
	}
	rs[0] = unicode.ToUpper(rs[0])
	return string(rs)
}
