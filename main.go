package main

import (
	"fmt"
	"os"
	"rocks/interpeter"
	lexer "rocks/lexer"
	parserIm "rocks/parser"
)

func main() {
	fmt.Println("START: ")
	var L1 lexer.Lexer
	Data, _ := os.ReadFile("testFile.ro")
	L1.Input = string(Data)

	L1.CorrectBackSlash()
	L1.LexerAll()
	L1.AddEOF()
	parser := parserIm.Parser{}
	parser.Input = L1.Tokens
	parser.Parsing()
	fmt.Println(parser.Output)

	NewEnv := interpeter.NewEnvironment(parser.Output, map[string]parserIm.Function{}, map[string]interpeter.Ident{})
	NewEnv.Interpeter()

}
