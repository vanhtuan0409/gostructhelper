package main

import (
	"flag"
	"fmt"
	"go/token"
	"os"
	"strings"

	helper "github.com/vanhtuan0409/gostructhelper"
)

func usage() {
	fmt.Fprintln(os.Stderr, "Usage of gostructhelper:")
	flag.PrintDefaults()
}

func main() {
	name := flag.String("type", "", "Struct name to generate")
	path := flag.String("file", "", "Source code")
	shouldWrite := flag.Bool("write", false, "Overwrite source file")
	stdin := flag.Bool("stdin", false, "Read from stdin")
	genConstructor := flag.Bool("constructor", false, "Generate constructor")
	genGetter := flag.Bool("getter", false, "Generate getter")
	flag.Parse()
	if *name == "" {
		usage()
		return
	}
	*name = strings.TrimSpace(*name)

	// Generate registry
	fs := token.NewFileSet()
	reg := helper.NewRegistry()
	if *genGetter {
		reg.Register(helper.NewGetterGenerator(*name, fs))
	}
	if *genConstructor {
		reg.Register(helper.NewConstructorGenerator(*name))
	}

	// Open file
	var inFile *os.File = nil
	if *path != "" {
		var err error
		inFile, err = os.OpenFile(*path, os.O_RDWR, 0o644)
		if err != nil {
			panic(err)
		}
		defer inFile.Close()
	}

	// Choose input
	in := os.Stdin
	if !*stdin && *path != "" {
		in = inFile
	}
	s := helper.NewSource(
		helper.SourceWithReader(in),
		helper.SourceWithPath(*path),
	)

	// Choose output
	out := os.Stdout
	if *shouldWrite {
		out = inFile
	}
	defer out.Close()

	helper.Gen(reg, s, fs, out)
}
