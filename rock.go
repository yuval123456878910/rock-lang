package main

import (
	"os"

	"rocks/interpeter"
	lexer "rocks/lexer"
	parserIm "rocks/parser"
)

func main() {
	Location := os.Args[1]
	var L1 lexer.Lexer
	Data, _ := os.ReadFile(Location)
	L1.Input = string(Data)

	L1.CorrectBackSlash()
	L1.LexerAll()
	L1.AddEOF()
	//fmt.Println(L1.Tokens)
	parser := parserIm.Parser{}
	parser.Input = L1.Tokens
	parser.Parsing()

	//litter.Dump(parser.Output...)
	NewEnv := interpeter.NewEnvironment(parser.Output, map[string]parserIm.Function{}, map[string]interpeter.Ident{})
	NewEnv.Interpeter()
}
