package main

import (
	"flag"
	"os"

	helper "github.com/vanhtuan0409/gostructhelper"
)

func main() {
	name := flag.String("type", "", "Struct name to generate")
	path := flag.String("file", "", "Source code")
	shouldWrite := flag.Bool("write", false, "Overwrite source file")
	stdin := flag.Bool("stdin", false, "Read from stdin")
	disableConstructor := flag.Bool("no-constructor", false, "Generate constructor")
	disableGetter := flag.Bool("no-getter", false, "Generate getter")
	flag.Parse()
	if *name == "" {
		flag.Usage()
		return
	}

	// Generate registry
	reg := helper.NewRegistry()
	if !*disableConstructor {
		reg.Register(helper.NewConstructorGenerator(*name))
	}
	if !*disableGetter {
		reg.Register(helper.NewGetterGenerator(*name))
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

	helper.Gen(reg, s, out)
}
