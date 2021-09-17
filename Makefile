testdata:
	cat ./exampledatas/main.go | go run ./cmds/gostructhelper -stdin -file ./exampledatas/main.go -type MyType -constructor -getter
