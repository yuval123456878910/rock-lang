/*
TODO: LEXER
add  UNKOWN (ilegal): Done
*/

package lexer

import (
	"fmt"
	"slices"
	"strings"
)

const (
	OneComment        = "#"
	LinesCommentStart = "#!"
	LinesCommentEnd   = "!#"
)
const (
	LITERAL    = "literal"
	STRING     = "string"
	OPERATOR   = "operator"
	PUNCTUATOR = "punctuator"
	IDENTIFIER = "identifier"
	KEYWORD    = "keyword"
	UNKOWN     = "unknown"
	CHAR       = "char"
	NEWLINE    = "newline"
	EOF        = "EOF"
)

var Keywords []string = []string{"struct", "func", "reach", "int", "string", "float", "any", "var", "const", "list", "dict", "if", "else", "elseif", "break", "continue", "thread", "await", "return", "char", "house", "while"}
var Operators []byte = []byte{'+', '-', '/', '*', '=', '<', '>', '|', '!'}
var LegalOperators []string = []string{"+", "-", "/", "=", "+=", "-=", "*=", "/=", "<=", ">=", "|>", "==", "!="}
var Punctuators []byte = []byte{'(', ')', '[', ']', '{', '}', ':', ',', '.'}

type Lexer struct {
	Input  string
	ch     byte
	pos    int
	Tokens []Token
}

type Token struct {
	Value string
	Type  string
}

func (l *Lexer) Next() {
	l.pos++
	if l.pos >= len(l.Input) {
		l.ch = 0
		return
	}
	l.ch = l.Input[l.pos]
}

func IsLetter(digit byte) bool {
	return digit >= 'a' && digit <= 'z' || digit >= 'A' && digit <= 'Z' || digit == '_'
}

func IsNumber(digit byte) bool {
	return digit >= '0' && digit <= '9' && digit != ' '
}

func IsSpace(digit byte) bool {
	return digit == ' ' || digit == '\n' || digit == '\t' || digit == '\r'
}

func IsUnknowen(digit byte) bool {
	return !IsLetter(digit) && !IsNumber(digit) && !slices.Contains(Punctuators, digit) && !IsOporator(digit) && !IsSpace(digit) && string(digit) != `'` && string(digit) != `"`
}

func (l *Lexer) SkipSpaces() {
	if len(l.Input) <= 0 {
		l.pos = 0
		return
	}
	if l.pos == 0 {
		l.ch = l.Input[0]
	}

	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.Next()
	}
}

func (l *Lexer) SneekPeak() byte {
	if l.pos+1 >= len(l.Input) {
		return 0
	}
	return l.Input[l.pos+1]
}

func (l *Lexer) BackPeak() byte {
	if l.pos-1 < 0 {
		return 0
	}

	return l.Input[l.pos-1]
}

func (l *Lexer) DoubleSneekPeak() byte {
	if l.pos+2 >= len(l.Input) {
		return 0
	}
	return l.Input[l.pos+2]
}

func IsOporator(digit byte) bool {
	return slices.Contains(Operators, digit)
}

func FindLastNotSpace(Text string, startPos int) int {
	for startPos >= 0 {
		if Text[startPos] == ' ' {
			startPos--
		} else {
			break
		}
	}
	return startPos
}

func CanBeNegetive(Text string, pos int) bool {
	if pos == 0 {
		return true
	}
	if FindLastNotSpace(Text, pos-1) < 0 {
		return true
	}
	last := Text[FindLastNotSpace(Text, pos-1)]
	return ((IsOporator(last)) || slices.Contains(Punctuators, last))

}

func (l *Lexer) IdentifyWord() {
	l.SkipSpaces()

	for l.ch == '\n' && l.pos < len(l.Input) {
		var NewToken Token
		NewToken.Type = NEWLINE
		NewToken.Value = ""
		l.Tokens = append(l.Tokens, NewToken)
		l.Next()
	}

	if l.pos >= len(l.Input) {
		return
	}
	if IsLetter(l.ch) {
		var NewToken Token

		StartingPos := l.pos
		for (IsLetter(l.ch) || IsNumber(l.ch)) && l.pos < len(l.Input) {
			l.Next()
		}
		NewToken.Value = string(l.Input[StartingPos:l.pos])
		if slices.Contains(Keywords, NewToken.Value) {
			NewToken.Type = KEYWORD
		} else {
			NewToken.Type = IDENTIFIER
		}
		l.Tokens = append(l.Tokens, NewToken)

	} else if (IsNumber(l.SneekPeak()) && l.ch == '-' && (len(l.Tokens) == 0 || l.Tokens[len(l.Tokens)-1].Type == OPERATOR || (l.Tokens[len(l.Tokens)-1].Type != LITERAL && l.Tokens[len(l.Tokens)-1].Type != IDENTIFIER))) || IsNumber(l.ch) {
		var NewToken Token

		StartingPos := l.pos
		if l.ch == '-' {
			l.Next()
		}
		for (IsNumber(l.ch) || l.ch == '.') && l.pos < len(l.Input) {
			l.Next()
		}
		NewToken.Value = string(l.Input[StartingPos:l.pos])
		NewToken.Type = LITERAL
		l.Tokens = append(l.Tokens, NewToken)
	} else if string(l.ch) == `'` {
		var NewToken Token

		StartingPos := l.pos + 1
		l.Next()
		for string(l.ch) != `'` && l.pos < len(l.Input) {
			l.Next()
		}

		NewToken.Value = string(l.Input[StartingPos:l.pos])
		if len(NewToken.Value) > 3 {
			fmt.Println("Char cant be more than one charecter!")
			return
		}
		NewToken.Type = CHAR
		l.Tokens = append(l.Tokens, NewToken)
		l.Next()
	} else if string(l.ch) == `"` {
		var NewToken Token

		StartingPos := l.pos + 1
		l.Next()
		for string(l.ch) != `"` && l.pos < len(l.Input) {
			l.Next()
		}

		NewToken.Value = string(l.Input[StartingPos:l.pos])
		NewToken.Value = strings.ReplaceAll(NewToken.Value, `\n`, "\n")
		NewToken.Value = strings.ReplaceAll(NewToken.Value, `\t`, "\t")
		NewToken.Value = strings.ReplaceAll(NewToken.Value, `\r`, "\r")
		NewToken.Type = STRING
		l.Tokens = append(l.Tokens, NewToken)
		l.Next()

	} else if slices.Contains(Punctuators, l.ch) {
		var NewToken Token
		NewToken.Type = PUNCTUATOR
		NewToken.Value = string(l.ch)
		l.Next()
		l.Tokens = append(l.Tokens, NewToken)

	} else if IsUnknowen(l.ch) {
		var NewToken Token

		StartingPos := l.pos
		l.Next()
		for IsUnknowen(l.ch) && l.pos < len(l.Input) {
			l.Next()
		}

		if string(l.Input[StartingPos:l.pos]) == OneComment && l.pos < len(l.Input) {
			for l.ch != '\n' && l.pos < len(l.Input) {
				l.Next()
			}
			return
		}
		if string(l.Input[StartingPos:l.pos]) == LinesCommentStart && l.pos < len(l.Input) {
			if !(string(l.ch)+string(l.SneekPeak()) == LinesCommentEnd) && l.pos < len(l.Input) {
				for !(string(l.ch)+string(l.SneekPeak()) == LinesCommentEnd) && l.pos < len(l.Input) {
					l.Next()
				}
			}
			l.Next() // skip !
			l.Next() // skip #
			return
		}
		NewToken.Value = string(l.Input[StartingPos:l.pos])

		NewToken.Type = UNKOWN
		l.Tokens = append(l.Tokens, NewToken)
		l.Next()
	} else if IsOporator(l.ch) {
		var NewToken Token

		StartingPos := l.pos

		for (IsOporator(l.ch) && slices.Contains(LegalOperators, string(l.Input[StartingPos:l.pos+1])) || IsOporator(l.ch) && l.pos == StartingPos) && l.pos < len(l.Input) {
			l.Next()

		}
		NewToken.Value = string(l.Input[StartingPos:l.pos])
		NewToken.Type = OPERATOR
		l.Tokens = append(l.Tokens, NewToken)
	}
}
func (l *Lexer) LexerAll() {
	for l.pos < len(l.Input) {
		l.IdentifyWord()
	}
}

func (l *Lexer) CorrectBackSlash() {
	l.Input = strings.ReplaceAll(l.Input, `\n`, "\n")
	l.Input = strings.ReplaceAll(l.Input, `\t`, "\t")
	l.Input = strings.ReplaceAll(l.Input, `\r`, "\r")
}

func (l *Lexer) AddEOF() {
	TempEOF := Token{Value: "", Type: "EOF"}
	l.Tokens = append(l.Tokens, TempEOF)
}
