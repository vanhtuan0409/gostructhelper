package gostructhelper

import (
	"golang.org/x/tools/go/ast/astutil"
)

type Generator interface {
	Name() string

	// Accept return true if want to continue
	Accept(c *astutil.Cursor) bool
}

type Registry struct {
	done       map[string]bool
	generators []Generator
}

func NewRegistry() *Registry {
	return &Registry{
		done:       map[string]bool{},
		generators: []Generator{},
	}
}

func (r *Registry) Register(g Generator) {
	r.generators = append(r.generators, g)
}

func (r *Registry) isDone(g Generator) bool {
	return r.done[g.Name()]
}

func (r *Registry) isAllDone() bool {
	for _, done := range r.done {
		if !done {
			return false
		}
	}
	return true
}
