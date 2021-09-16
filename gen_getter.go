package gostructhelper

import "golang.org/x/tools/go/ast/astutil"

type getterGenerator struct {
	typeName string
}

func NewGetterGenerator(name string) *getterGenerator {
	return &getterGenerator{
		typeName: name,
	}
}

func (g *getterGenerator) Name() string {
	return "getter"
}

func (g *getterGenerator) Accept(c *astutil.Cursor) bool {
	_, ok := filterTypeName(c.Node(), g.typeName)
	if !ok {
		return true
	}

	return false
}
