package gostructhelper

import "go/ast"

func filterTypeName(node ast.Node, name string) (*ast.TypeSpec, bool) {
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
