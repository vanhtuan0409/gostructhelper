package main

import (
	"context"
	"sync"
)

type MyType struct {
	foo string
	Bar string
	ctx context.Context

	sync.Mutex
}

func main() {}
