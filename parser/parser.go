package parser

import (
	"fmt"
	"os"
	"rocks/lexer"
	"slices"
)

var aviableTypes []string = []string{
	"string", "float", "list", "dict", "bool", "char", "int",
}

func FindNexer(StartPos int, Tokens []lexer.Token, startFind lexer.Token, endFind lexer.Token) (int, error) {
	SkipEnds := 0 // this ment to fix when the code find } but it isnt the end
	for idx, token := range Tokens[StartPos+1:] {
		if Equals2Token(token, endFind) {

			if SkipEnds == 0 {
				return StartPos + 1 + idx, nil
			}
			SkipEnds--
		} else if Equals2Token(token, startFind) {
			SkipEnds++
		}

	}

	return 0, fmt.Errorf("The system coudnt find end for '%s', %d", endFind.Value, SkipEnds)
}

type NewIdent struct {
	Name    string
	Type    string
	IsConst bool
	Content any
}

type RefactIdent struct {
	Name    string
	Content any
}

type Parser struct {
	Input  []lexer.Token
	Output []any
}

func (p *Parser) Parsing() {
	p.Output = Parse(p.Input)
}

type Function struct {
	Name       string
	Perameters []Parimiter  // type and name
	Returns    []TypeReturn // types
	Body       []any        // all the function code
}

type Parimiter struct {
	Type string
	Name string
}

type TypeReturn struct {
	Type string
}

type Reach struct {
	Path string
}

type Moudle struct {
	Body []any
}

type Return struct {
	Exprs []any
}

type House struct {
	Names    []string
	Contents []any
}

type IfStm struct {
	Condition any
	Body      []any
	Else      *IfStm
}

func SearchStartToken(Array []lexer.Token, Start int, funcApply func(item any) any, Item any) (int, error) {
	for i := Start; i < len(Array); i++ {
		if funcApply(Array[i]) == Item {
			return i, nil
		}
	}
	return 0, fmt.Errorf("Coudnt find the item: %V", Item)
}

func SearchEnd(Tokens []lexer.Token, pos int) int {
	End := pos
	for End+1 < len(Tokens) && Tokens[End].Type != lexer.NEWLINE && Tokens[End].Type != lexer.EOF {
		End++
	}
	return End
}

func ReturnsSepDouble(Tokens []lexer.Token, sep lexer.Token) [][]lexer.Token {
	TokenSeperation := [][]lexer.Token{}
	sepPos := 0
	Level := 0
	for pos := 0; pos < len(Tokens); pos++ {
		if Tokens[pos] == sep && Level <= 0 {
			TokenSeperation = append(TokenSeperation, Tokens[sepPos:pos])
			sepPos = pos + 1
			continue
		}
		switch Tokens[pos].Value {
		case "(", "{", "[":
			Level++
			continue
		case ")", "}", "]":
			Level--
			continue
		}
	}
	if len(Tokens) >= sepPos {
		TokenSeperation = append(TokenSeperation, Tokens[sepPos:])
	}
	return TokenSeperation
}

func Equals2Token(FirstToken lexer.Token, SecondToken lexer.Token) bool {
	return (FirstToken.Value == SecondToken.Value) && (FirstToken.Type == SecondToken.Type)
}

func ReturnsSepOnceTokens(Tokens []lexer.Token, sep lexer.Token) []lexer.Token {
	TokenSeperation := []lexer.Token{}
	Level := 0
	LastToken := Tokens[0]
	for pos := 0; pos < len(Tokens); pos++ {
		if Tokens[pos] == sep && Level <= 0 {
			TokenSeperation = append(TokenSeperation, LastToken)
			LastToken.Value = ""
			LastToken.Value = ""
			continue
		}
		switch Tokens[pos].Value {
		case "(", "{", "[":
			Level++
			continue
		case ")", "}", "]":
			Level--
			continue
		}

		LastToken = Tokens[pos]
	}
	if !Equals2Token(LastToken, lexer.Token{Value: "", Type: ""}) {
		TokenSeperation = append(TokenSeperation, LastToken)
	}
	return TokenSeperation
}
func Parse(Tokens []lexer.Token) []any {
	Global_Result := []any{}
	for pos := 0; pos < len(Tokens); pos++ {
		Token := Tokens[pos]
		if pos > 0 && Tokens[pos-1].Type == lexer.IDENTIFIER && Token.Value == "(" && Token.Type == lexer.PUNCTUATOR {

			End, err := FindClose(Tokens, pos+1, "(", ")")
			if err != nil {
				fmt.Println("Error: End of the call funcion wasnt found!")
			}
			Call := CallFunction{Name: Tokens[pos-1].Value}
			InPrunctiotion := false

			var CurrentArg []lexer.Token
			for _, token := range Tokens[pos+1 : End] {
				if token.Type == lexer.PUNCTUATOR && slices.Contains([]string{"{", "}", "(", ")", "{", "}"}, token.Value) {
					InPrunctiotion = true
					pos++
				}

				if token.Type == lexer.PUNCTUATOR && token.Value == "," && !InPrunctiotion {

					TempExpress := Expretion{Tokens: CurrentArg}
					Call.ParimitersInput = append(Call.ParimitersInput, ParseBinding(&TempExpress, 0))
					CurrentArg = []lexer.Token{}
					continue
				}
				CurrentArg = append(CurrentArg, token)

			}
			if len(CurrentArg) > 0 {
				TempExpress := Expretion{Tokens: CurrentArg}
				Call.ParimitersInput = append(Call.ParimitersInput, ParseBinding(&TempExpress, 0))
			}
			Global_Result = append(Global_Result, Call)
			pos = End + 1
			continue
		}
		switch Token.Value {
		case "func":
			// syntax: func main(int main) (int) {}

			var Current_Result Function
			Current_Result.Name = Tokens[pos+1].Value // get name of the function like
			if Tokens[pos+2].Value != "(" {
				fmt.Println("Where the decleration '()'")
				return []any{}
			}
			StartToken := lexer.Token{Value: "(", Type: lexer.PUNCTUATOR}
			EndToken := lexer.Token{Value: ")", Type: lexer.PUNCTUATOR}
			index, err := FindNexer(pos+3, Tokens, StartToken, EndToken)

			if err != nil {
				fmt.Println("Error: The perameters didnt end!")
				return nil
			}

			if index+1 > pos {
				for idx, perameters := range Tokens[pos+3 : index+1] {
					if perameters.Type == lexer.KEYWORD && slices.Contains(aviableTypes, perameters.Value) {
						TempPar := Parimiter{Type: perameters.Value, Name: Tokens[pos+4:][idx].Value}
						Current_Result.Perameters = append(Current_Result.Perameters, TempPar)
					}
				}
			}

			if Tokens[index+2].Value == "(" && Tokens[index+2].Type == lexer.PUNCTUATOR {

				StartToken2 := lexer.Token{Value: "(", Type: lexer.PUNCTUATOR}
				EndToken2 := lexer.Token{Value: ")", Type: lexer.PUNCTUATOR}
				index2, err := FindNexer(index+2, Tokens, StartToken2, EndToken2)
				if err != nil {
					fmt.Println("An error had accourd, mayby () didnt closed!")
					return []any{}
				}
				for _, returns := range Tokens[index+3 : index2+1] {
					if returns.Type == lexer.KEYWORD && slices.Contains(aviableTypes, returns.Value) {
						Current_Result.Returns = append(Current_Result.Returns, TypeReturn{Type: returns.Value})
					} else {
						fmt.Println("You included untype as a return parameter")
						return []any{}
					}

				}
			}

			StartBodyIdx, errBefor := SearchStartToken(Tokens, pos, func(item any) any { return item.(lexer.Token).Value }, "{")

			if errBefor != nil {
				fmt.Println("Codnt find start!")
			}
			StartToken2 := lexer.Token{Value: "{", Type: lexer.PUNCTUATOR}
			EndToken2 := lexer.Token{Value: "}", Type: lexer.PUNCTUATOR}
			EndIdx, err := FindNexer(StartBodyIdx, Tokens, StartToken2, EndToken2)

			if err != nil {
				fmt.Println("Error 105, The function didn't end")
			}

			Current_Result.Body = Parse(Tokens[StartBodyIdx+1 : EndIdx])

			Global_Result = append(Global_Result, Current_Result)
			pos = EndIdx + 1
			continue

		case "var", "const":

			NewIdentefire := NewIdent{}
			NameToken := Tokens[pos+1]
			if NameToken.Type != lexer.IDENTIFIER {
				fmt.Println("Token isn't an identefire,", NameToken.Value)
				return []any{}
			}
			NewIdentefire.Name = NameToken.Value

			TypeToken := Tokens[pos+2]
			if !slices.Contains(aviableTypes, TypeToken.Value) {
				fmt.Printf("'%s' isn't a type in decleration of value '%s'\n", TypeToken.Value, NameToken.Value)
				return []any{}
			}
			NewIdentefire.Type = TypeToken.Value
			if pos+3 >= len(Tokens) || !(Tokens[pos+3].Value == "=" && Tokens[pos+3].Type == lexer.OPERATOR) {
				Global_Result = append(Global_Result, NewIdentefire)
				continue
			}
			EndOfLine := SearchEnd(Tokens, pos)

			NewExpretion := Expretion{Tokens: Tokens[pos+4 : EndOfLine+1]}
			NewIdentefire.IsConst = Tokens[pos].Value == "const"
			AST := ParseBinding(&NewExpretion, 0)
			NewIdentefire.Content = AST
			pos = EndOfLine

			Global_Result = append(Global_Result, NewIdentefire)
			continue
		case "reach":
			NewReach := Reach{}
			path := Tokens[pos+1]
			if path.Type != lexer.STRING {
				fmt.Println("The reach path isnt a string")
				return []any{}
			}
			NewReach.Path = path.Value
			Global_Result = append(Global_Result, NewReach)
			pos = SearchEnd(Tokens, pos)
			continue
		case "return":
			EndLine := SearchEnd(Tokens, pos)

			Contexts := ReturnsSepDouble(Tokens[pos+1:EndLine+1], lexer.Token{Value: ",", Type: lexer.PUNCTUATOR})
			exprs := []any{}
			for _, value := range Contexts {
				tempExp := Expretion{Tokens: value}
				exprs = append(exprs, ParseBinding(&tempExp, 0))
			}

			NewReturn := Return{Exprs: exprs}
			Global_Result = append(Global_Result, NewReturn)
			pos = EndLine
			continue
		case "house":
			house := House{}

			EndLine := SearchEnd(Tokens, pos)
			EqlLocation := pos

			for !Equals2Token(Tokens[EqlLocation], lexer.Token{Value: "=", Type: lexer.OPERATOR}) {
				EqlLocation++
			}
			VarsTokens := ReturnsSepOnceTokens(Tokens[pos+1:EqlLocation], lexer.Token{Value: ",", Type: lexer.PUNCTUATOR})

			ContentSlice := ReturnsSepDouble(Tokens[EqlLocation+1:EndLine+1], lexer.Token{Value: ",", Type: lexer.PUNCTUATOR})
			TempContextBinding := []any{}
			for _, context := range ContentSlice {
				NewExpr := Expretion{Tokens: context}
				TempContextBinding = append(TempContextBinding, ParseBinding(&NewExpr, 0))
			}
			house.Contents = TempContextBinding
			Vars := []string{}
			for _, token := range VarsTokens {
				Vars = append(Vars, token.Value)
			}
			house.Names = Vars
			pos = EndLine
			Global_Result = append(Global_Result, house)
			continue
		case "=":
			if Tokens[pos].Type != lexer.OPERATOR {
				continue
			}
			name := Tokens[pos-1]
			end := SearchEnd(Tokens, pos)
			NewEnv := Expretion{Tokens: Tokens[pos+1 : end]}
			contect := ParseBinding(&NewEnv, 0)
			NewRefactIdent := RefactIdent{Name: name.Value, Content: contect}
			Global_Result = append(Global_Result, NewRefactIdent)
			pos = end
			continue

		case "if":
			var NewIf IfStm
			StartToken := lexer.Token{Value: "{", Type: lexer.PUNCTUATOR}
			EndToken := lexer.Token{Value: "}", Type: lexer.PUNCTUATOR}

			BlockStart, err := SearchStartToken(Tokens, pos, func(item any) any { return item.(lexer.Token).Value }, "{")
			if err != nil {
				fmt.Println("Coudnt find a start to { !")
				os.Exit(0)
			}

			EndBody, err2 := FindNexer(BlockStart, Tokens, StartToken, EndToken)
			if err2 != nil {
				fmt.Println(err2.Error())
				os.Exit(0)
			}

			// Condition between 'if' and '{'
			NewExpr := Expretion{Tokens: Tokens[pos+1 : BlockStart]}
			NewIf.Condition = ParseBinding(&NewExpr, 0)

			// Body tokens strictly inside '{' and '}'
			NewIf.Body = Parse(Tokens[BlockStart+1 : EndBody])

			// Set pos directly to EndBody ('}').
			// The main for-loop's pos++ automatically steps to the next token.
			pos = EndBody
			var PointerToIf *IfStm = &NewIf

		loop:
			for pos+1 < len(Tokens) && (Equals2Token(Tokens[pos+1], lexer.Token{Value: "elseif", Type: lexer.KEYWORD}) || Equals2Token(Tokens[pos+1], lexer.Token{Value: "else", Type: lexer.KEYWORD})) {
				pos++ // Move to 'elseif' or 'else'
				switch Tokens[pos].Value {
				case "elseif":
					var NewIfElse IfStm
					BlockStart2, err3 := SearchStartToken(Tokens, pos, func(item any) any { return item.(lexer.Token).Value }, "{")
					if err3 != nil {
						fmt.Println("Coudnt find a start to { !")
						os.Exit(0)
					}

					EndBody2, err4 := FindNexer(BlockStart2, Tokens, StartToken, EndToken)
					if err4 != nil {
						fmt.Println(err4.Error())
						os.Exit(0)
					}

					NewExpr2 := Expretion{Tokens: Tokens[pos+1 : BlockStart2]}
					NewIfElse.Condition = ParseBinding(&NewExpr2, 0)
					NewIfElse.Body = Parse(Tokens[BlockStart2+1 : EndBody2])

					PointerToIf.Else = &NewIfElse
					PointerToIf = &NewIfElse
					pos = EndBody2

				case "else":
					var NewElse IfStm
					BlockStart3, err := SearchStartToken(Tokens, pos, func(item any) any { return item.(lexer.Token).Value }, "{")
					if err != nil {
						fmt.Println("Coudnt find a start to { !")
						os.Exit(0)
					}

					EndBody3, err5 := FindNexer(BlockStart3, Tokens, StartToken, EndToken)
					if err5 != nil {
						fmt.Println("Coudnt find an end to else!")
						os.Exit(1)
					}

					NewElse.Body = Parse(Tokens[BlockStart3+1 : EndBody3])
					NewElse.Condition = 1
					PointerToIf.Else = &NewElse

					PointerToIf = &NewElse
					pos = EndBody3
					break loop
				}
			}

			Global_Result = append(Global_Result, NewIf)
			continue
		}

	}
	if len(Global_Result) == 0 {
		return []any{}
	}
	return Global_Result
}

type Expretion struct {
	Tokens []lexer.Token
	Pos    int
}
type Oporation struct {
	Op    string
	Left  any
	Right any
}

type CallFunction struct {
	Name            string
	ParimitersInput []any
}

func (e *Expretion) Next() (lexer.Token, error) {
	e.Pos++
	if e.Pos >= len(e.Tokens) {
		return lexer.Token{}, fmt.Errorf("Limit reached")
	}
	return e.Tokens[e.Pos], nil

}

func IsExpress(token lexer.Token) bool {
	return token.Type == lexer.OPERATOR || token.Type == lexer.IDENTIFIER || token.Type == lexer.LITERAL
}

var Binding map[string]float32 = map[string]float32{
	"==": 0,
	"<=": 0,
	">=": 0,
	"+":  1,
	"-":  1,
	"*":  2,
	"/":  2,
}

func RetriveBinding(Value string) float32 {
	v, _ := Binding[Value]
	return v
}

func IsOporator(token lexer.Token) bool {
	return token.Type == lexer.OPERATOR
}

func FindClose(tokens []lexer.Token, startPos int, start string, end string) (int, error) {
	DistanceOfFinish := 0
	for idx, token := range tokens[startPos:] {
		if token.Type == lexer.PUNCTUATOR {
			switch token.Value {
			case start:
				DistanceOfFinish++
			case end:
				if DistanceOfFinish == 0 {
					return startPos + idx, nil
				}
				DistanceOfFinish--

			}
		}
	}
	return 0, fmt.Errorf("No end had bean found!")
}

type List struct {
	Items []any
}

// error in the ParseBinding has to start at 2
func ParseBinding(Express *Expretion, min_bind float32) any {
	Leftside := any(Express.Tokens[Express.Pos])
	if Leftside.(lexer.Token).Type == lexer.IDENTIFIER && Express.Pos+1 < len(Express.Tokens) && Express.Tokens[Express.Pos+1].Value == "(" {
		End, err := FindClose(Express.Tokens, Express.Pos+2, "(", ")")
		if err != nil {
			fmt.Println("Error: End of the call funcion wasnt found!")
		}
		Call := CallFunction{Name: Express.Tokens[Express.Pos].Value}
		depth := 0
		var CurrentArg []lexer.Token
		for _, token := range Express.Tokens[Express.Pos+2 : End] {

			if token.Type == lexer.PUNCTUATOR && slices.Contains([]string{"[", "(", "{"}, token.Value) {
				depth++

			} else if token.Type == lexer.PUNCTUATOR && slices.Contains([]string{"]", ")", "}"}, token.Value) {
				depth--

			}

			if token.Type == lexer.PUNCTUATOR && token.Value == "," && depth == 0 {
				TempExpress := Expretion{Tokens: CurrentArg}
				Call.ParimitersInput = append(Call.ParimitersInput, ParseBinding(&TempExpress, 0))
				CurrentArg = []lexer.Token{}
				continue

			}
			CurrentArg = append(CurrentArg, token)

		}
		if len(CurrentArg) > 0 {
			TempExpress := Expretion{Tokens: CurrentArg}
			Call.ParimitersInput = append(Call.ParimitersInput, ParseBinding(&TempExpress, 0))
		}
		Leftside = Call
		Express.Pos = End + 1
	} else if Leftside.(lexer.Token).Value == "(" {
		End, err := FindClose(Express.Tokens, Express.Pos+1, "(", ")")
		if err != nil {
			fmt.Println("No, end to '(' had been found")
			return Leftside
		}
		TempExpress := Expretion{Tokens: Express.Tokens[Express.Pos+1 : End]}
		Leftside = ParseBinding(&TempExpress, 0)
		Express.Pos = End + 1
	} else if Leftside.(lexer.Token).Value == "[" && Leftside.(lexer.Token).Type == lexer.PUNCTUATOR {

		End, err := FindClose(Express.Tokens, Express.Pos+1, "[", "]")
		if err != nil {
			fmt.Println("No, end to '[' had been found")
			return Leftside
		}
		Temped := []lexer.Token{}
		var NewList List
		var InInside int = 0
		for pos := Express.Pos + 1; pos <= End; pos++ {
			if InInside > 0 {
				Temped = append(Temped, Express.Tokens[pos])
				continue
			}
			if Express.Tokens[pos].Value == "[" {
				InInside++
				Temped = append(Temped, Express.Tokens[pos])
				continue
			} else if Express.Tokens[pos].Value == "]" {
				InInside--
				Temped = append(Temped, Express.Tokens[pos])
				continue
			}
			if Express.Tokens[pos].Value == "," && Express.Tokens[pos].Type == lexer.PUNCTUATOR {
				TempExpress := Expretion{Tokens: Temped}
				NewList.Items = append(NewList.Items, ParseBinding(&TempExpress, 0))
				Temped = []lexer.Token{}

			} else {
				Temped = append(Temped, Express.Tokens[pos])
			}
		}
		if len(Temped) > 0 {
			TempExpress := Expretion{Tokens: Temped}
			NewList.Items = append(NewList.Items, ParseBinding(&TempExpress, 0))
		}
		Leftside = NewList
		Express.Pos = End + 1

	} else {
		Express.Next()
	}
	for Express.Pos < len(Express.Tokens) {
		CurrentToken := Express.Tokens[Express.Pos]

		if !IsOporator(CurrentToken) {
			Express.Next()
			continue
		}

		if RetriveBinding(CurrentToken.Value) >= float32(min_bind) {
			Express.Next()
			var Rightside any = ParseBinding(Express, RetriveBinding(CurrentToken.Value)+1)
			Leftside = Oporation{Op: CurrentToken.Value, Right: Rightside, Left: Leftside}
		} else {
			break
		}
	}
	return Leftside
}
