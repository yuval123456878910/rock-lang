package interpeter

/* Add:
add two set
*/

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"maps"
	"math"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"

	"rocks/lexer"
	"rocks/parser"
)

func IsLink(text string) bool {
	data, err := url.ParseRequestURI(text)
	if err != nil {
		return false
	}
	return (data.Scheme == "https" || data.Scheme == "http") && data.Host != ""
}

type ThreadChannel struct {
	Values any
	Types  []string
}

var (
	FuncType         string   = ReturnType(parser.Function{})
	ReachType        string   = ReturnType(parser.Reach{})
	NewIdentType     string   = ReturnType(parser.NewIdent{})
	OporationType    string   = ReturnType(parser.Oporation{})
	LiteralsLisr     []string = []string{lexer.LITERAL, lexer.CHAR, lexer.STRING}
	ApproveSideToOp  []string = []string{ReturnType(0), ReturnType(0.0), ReturnType(""), ReturnType(byte(0))}
	IntType          string   = ReturnType(0)
	RefactType       string   = ReturnType(parser.RefactIdent{})
	HouseType        string   = ReturnType(parser.House{})
	IfType           string   = ReturnType(parser.IfStm{})
	WhileType        string   = ReturnType(parser.WhileLoop{})
	UnaryOpType      string   = ReturnType(parser.UnaryOp{})
	ThreadType       string   = ReturnType(parser.Thread{})
	AwaitType        string   = ReturnType(parser.Await{})
	BreakType        string   = ReturnType(parser.Break{})
	PipeType         string   = ReturnType(parser.Pipe{})
	DictType         string   = ReturnType(parser.Dictenary{})
	AccessMethodType string   = ReturnType(parser.AccessMethod{})
	SelectList       string   = ReturnType(parser.SectionList{})
)

type StructInstance struct {
	TypeName string
	Fields   map[string]*Ident
}

var ReachImported []string = []string{}

type Keyfunc func(args ...any) (any, []string)

type Environment struct {
	FuncMap     map[string]parser.Function
	VariableMap map[string]Ident
	ParseDate   []any
	Output      []any
	Keyfuncs    map[string]Keyfunc
	Returned    bool
	Breaked     bool
	StructMap   map[string]parser.Struct
}

func BoolToInt(Bool bool) int {
	if Bool {
		return 1
	}
	return 0
}

func GetDataFromSelectList(Slice []any, Start any, End any, Long bool) []any {
	if !Long {
		return []any{Slice[Start.(int)]}
	}

	if Start == nil {
		Start = 0
	}
	if End == nil {
		End = len(Slice)
	}

	return Slice[Start.(int):End.(int)]
}

func GetStartValueType(Env Environment, TypeGot string) (any, error) {
	if !slices.Contains(parser.AviableTypes, TypeGot) {
		return nil, fmt.Errorf("Undefined typed: %s", TypeGot)
	}
	var DefualValue any
	switch TypeGot {
	case "string":
		DefualValue = ""
	case "int":
		DefualValue = 0
	case "float":
		DefualValue = 0.0
	case "list":
		DefualValue = []any{}
	case "dict":
		DefualValue = map[any]any{}
	case "char":
		DefualValue = ' '
	}

	if structDec, ok := Env.StructMap[TypeGot]; ok && DefualValue == nil {
		DefualValue = structDec
	}

	return DefualValue, nil
}

func NewEnvironment(parseDate []any, funcMap map[string]parser.Function, variableMap map[string]Ident) Environment {
	TempEnv := Environment{}
	TempEnv.FuncMap = funcMap
	TempEnv.VariableMap = variableMap
	TempEnv.ParseDate = parseDate
	TempEnv.Keyfuncs = map[string]Keyfunc{}
	TempEnv.StructMap = map[string]parser.Struct{}
	TempEnv.Returned = false
	TempEnv.Breaked = false
	TempEnv.Keyfuncs["print"] = Keyfunc(func(args ...any) (any, []string) {
		for _, arg := range args {
			switch Targ := arg.(type) {
			case []any:
				fmt.Print(Targ...)
				fmt.Print(" ")
			case *Ident:
				fmt.Print(Targ.Value)
				fmt.Print(" ")
			default:
				fmt.Print(arg, " ")
			}
		}
		fmt.Println()
		return args, []string{"any"}
	})
	TempEnv.Keyfuncs["append"] = func(args ...any) (any, []string) {
		Data := args[0].([]any)[0].([]any)
		Data = append(Data, args[1])
		return Data, []string{"list"}
	}
	TempEnv.Keyfuncs["look"] = func(args ...any) (any, []string) {
		Data := args[0].([]any)[0].([]any)

		TypeReturn := "any"
		switch DataReturn := Data[args[1].(int)].(type) {
		case int:
			return DataReturn, []string{"int"}
		case float64:
			return DataReturn, []string{"float"}
		case string:
			return DataReturn, []string{"string"}
		}
		return Data[args[1].(int)], []string{TypeReturn}
	}
	TempEnv.Keyfuncs["scan"] = func(args ...any) (any, []string) {
		Text := args[0].(string)
		fmt.Print(Text)
		Scanner := bufio.NewScanner(os.Stdin)
		if Scanner.Scan() {
		}
		return Scanner.Text(), []string{"string"}
	}
	TempEnv.Keyfuncs["at"] = func(args ...any) (any, []string) {
		dict := args[0].(map[any]any)
		return dict[args[1]], []string{"any"}
	}
	return TempEnv
}

func ReturnType(obj any) string {
	return fmt.Sprintf("%T", obj)
}

type Ident struct {
	Value   any
	Name    string
	Type    string
	IsConst bool
}
type SaveDataEnv struct {
	FuncMapNames    []string
	VariableMapName []string
}

func LoadSeqKeys(seq iter.Seq[string]) []string {
	slice := []string{}
	for key := range seq {
		slice = append(slice, key)
	}
	return slice
}

func (s *SaveDataEnv) Save(env Environment) {
	s.VariableMapName = LoadSeqKeys(maps.Keys(env.VariableMap))
	s.FuncMapNames = LoadSeqKeys(maps.Keys(env.FuncMap))
}

func LoadNewData[T any](keys []string, data map[string]T) map[string]T {
	NewVars := map[string]T{}
	for _, key := range keys {
		NewVars[key] = data[key]
	}
	return NewVars
}

func (s *SaveDataEnv) LoadTo(env *Environment) {
	env.VariableMap = LoadNewData(s.VariableMapName, env.VariableMap)
	env.FuncMap = LoadNewData(s.FuncMapNames, env.FuncMap)
}

func GetToSelection(Env *Environment, Select *parser.SectionList) *any {
	var CurrentPointer any = Select
	_, ok := CurrentPointer.(*parser.SectionList).List.(parser.SectionList)
	for ok {
		return GetToSelection(Env, CurrentPointer.(*parser.SectionList).List.(*parser.SectionList))
	}
	var V any
	switch Typed := (CurrentPointer.(*parser.SectionList).List).(type) {
	case lexer.Token:
		V = Env.VariableMap[Typed.Value].Value

	default:
		V, _ = Evaluate(CurrentPointer.(*parser.SectionList).List, Env.VariableMap, Env.FuncMap, Env.Keyfuncs, false)
	}

	var S any
	var E any

	switch TypeV := V.(type) {
	case []any:
		S, _ = Evaluate(CurrentPointer.(*parser.SectionList).Start, Env.VariableMap, Env.FuncMap, Env.Keyfuncs, false)
		E, _ = Evaluate(CurrentPointer.(*parser.SectionList).End, Env.VariableMap, Env.FuncMap, Env.Keyfuncs, false)
		var Wrapper any = TypeV[S.(int):E.(int)]
		return &Wrapper
	case map[any]any:
		S, _ = Evaluate(CurrentPointer.(*parser.SectionList).Start, Env.VariableMap, Env.FuncMap, Env.Keyfuncs, false)
		var Wrapper any = TypeV[S]
		return &Wrapper
	}
	return &V
}

func Evaluate(CalData any, indentMap map[string]Ident, funcMap map[string]parser.Function, keyFuncs map[string]Keyfunc, ReturnSelectListPointer bool) (any, []string) {
	var Value any
	CalType := ReturnType(CalData)

	switch CalType {
	case ReturnType(lexer.Token{}):
		switch CalData.(lexer.Token).Type {
		case lexer.LITERAL:
			value, err := strconv.ParseFloat(CalData.(lexer.Token).Value, 64)
			if err != nil {
				fmt.Printf("Error converting string to float. Error: %s", err.Error())
				os.Exit(1)
			}
			if math.Mod(value, 1) == 0 && !strings.Contains(CalData.(lexer.Token).Value, ".") {
				return int(value), []string{"int"}
			}
			return value, []string{"float"}
		case lexer.STRING:
			return CalData.(lexer.Token).Value, []string{"string"}
		case lexer.CHAR:
			return byte(CalData.(lexer.Token).Value[0]), []string{"char"}
		case lexer.IDENTIFIER:

			TempIdent := CalData.(lexer.Token)
			IdentGot, ok := indentMap[TempIdent.Value]

			if !ok {
				fmt.Println("Coudnt find a variable named:", TempIdent.Value)
				os.Exit(1)
			}
			return IdentGot.Value, []string{indentMap[TempIdent.Value].Type}
		case lexer.NEWLINE:

			return []any{}, []string{}
		}
	case ReturnType(parser.List{}):

		TempList := CalData.(parser.List)
		NewList := []any{}
		NewTypes := []string{}
		for _, item := range TempList.Items {
			NewItems, NewType := Evaluate(item, indentMap, funcMap, keyFuncs, false)
			NewList = append(NewList, NewItems)
			NewTypes = append(NewTypes, NewType...)
		}
		return []any{NewList}, NewTypes
	case ReturnType(parser.CallFunction{}):
		TempCall := CalData.(parser.CallFunction)

		CallFunc, ok := funcMap[TempCall.Name]

		_, ok2 := keyFuncs[TempCall.Name]
		if !ok && !ok2 {
			fmt.Println("Error, coudnt find the func:", TempCall.Name)
			os.Exit(1)
		}
		CallVarMap := map[string]Ident{}
		TempValues := []any{}
		for idx, ident := range TempCall.ParimitersInput {
			callEval, _ := Evaluate(ident, indentMap, funcMap, keyFuncs, false)
			if ok2 {
				TempValues = append(TempValues, callEval)
			} else {
				CallVarMap[CallFunc.Perameters[idx].Name] = Ident{Value: callEval, Name: CallFunc.Perameters[idx].Name, Type: CallFunc.Perameters[idx].Type, IsConst: false}
			}

		}
		if funcData, ok3 := keyFuncs[TempCall.Name]; ok3 {
			l, t := funcData(TempValues...)

			return l, t
		}

		NewEnv := Environment{ParseDate: CallFunc.Body, VariableMap: CallVarMap, FuncMap: funcMap, Keyfuncs: keyFuncs}
		NewEnv.Interpeter()
		Types := []string{}

		for _, value := range CallFunc.Returns {
			Types = append(Types, value.Type)
		}
		if len(NewEnv.Output) == 1 {
			return NewEnv.Output[0], Types
		}
		return NewEnv.Output, Types
	case ReturnType(parser.AccessMethod{}):
		NewEnv := Environment{VariableMap: indentMap, FuncMap: funcMap, Keyfuncs: keyFuncs}

		return *ToMethod(CalData.(parser.AccessMethod), NewEnv), []string{"any"}
	case UnaryOpType:
		TempUnaryOp := CalData.(parser.UnaryOp)
		Value, _ := Evaluate(TempUnaryOp.Value, indentMap, funcMap, keyFuncs, false)
		switch ValType := Value.(type) {
		case int:
			return -ValType, []string{"int"}
		case float64:
			return -ValType, []string{"float"}
		default:
			fmt.Println("Cant negetive a none number type!")
			os.Exit(1)
		}
		return Value, []string{}
	case ThreadType:

		TempThread := CalData.(parser.Thread)
		Func := TempThread.Content
		ch := make(chan ThreadChannel)
		go func() {
			vals, types := Evaluate(Func, indentMap, funcMap, keyFuncs, false)
			ch <- ThreadChannel{Values: vals, Types: types}
			close(ch)
		}()
		var ReturnedData <-chan ThreadChannel = ch
		return []any{ReturnedData}, []string{"thread"}
	case AwaitType:
		TempAwait := CalData.(parser.Await)
		Ident := TempAwait.ThreadHandle.(lexer.Token)
		IdentData := indentMap[Ident.Value]
		if IdentData.Type != "thread" {
			fmt.Println("You gave a none thread type variable!", IdentData.Name)
		}
		ch := IdentData.Value.(<-chan ThreadChannel)

		resultData := <-ch
		return resultData.Values, resultData.Types
	case PipeType:
		TempPipe := CalData.(parser.Pipe)
		TempPass := TempPipe.PassTo.(parser.CallFunction)

		TempPass.ParimitersInput = append([]any{TempPipe.Arg}, TempPass.ParimitersInput...)
		EvalPass, types := Evaluate(TempPass, indentMap, funcMap, keyFuncs, false)

		return EvalPass, types
	case DictType:
		TempDict := CalData.(parser.Dictenary)
		Dict := TempDict.Elements
		Keys := maps.Keys(Dict)
		NewEvalMap := map[any]any{}
		for key := range Keys {
			NewRightEval := Dict[key]
			RightEval, _ := Evaluate(NewRightEval, indentMap, funcMap, keyFuncs, false)
			LeftEval, _ := Evaluate(key, indentMap, funcMap, keyFuncs, false)
			NewEvalMap[LeftEval] = RightEval
		}
		return NewEvalMap, []string{"dict"}
	case SelectList:
		var EvalOne any
		var EvalTwo any
		TempSelectList := CalData.(parser.SectionList)
		if TempSelectList.Start != nil {
			EvalOne, _ = Evaluate(TempSelectList.Start, indentMap, funcMap, keyFuncs, false)
		}
		if TempSelectList.End != nil {
			EvalTwo, _ = Evaluate(TempSelectList.End, indentMap, funcMap, keyFuncs, false)
		}
		switch value, _ := Evaluate(TempSelectList.List, indentMap, funcMap, keyFuncs, false); value.(type) {
		case []any:
			Section := GetDataFromSelectList(value.([]any)[0].([]any), EvalOne, EvalTwo, TempSelectList.Long)

			if len(Section) == 1 {
				if ReturnSelectListPointer {
					return &(Section[0]), []string{"any"}
				}
				return Section[0], []string{"any"}
			}
			return Section, []string{"list"}
		case map[any]any:
			if TempSelectList.Long == true {
				parser.Panic("Syntax error", "Cant have a select dict with two elements")
			}
			if TempSelectList.Start == nil {
				parser.Panic("Syntax error", "Dict select start is empty")
			}
			if TempSelectList.End != nil {
				parser.Panic("Syntax error", "Dict select end isn't empty")
			}
			EvalOne, Type := Evaluate(TempSelectList.Start, indentMap, funcMap, keyFuncs, false)
			if Type[0] != "string" {
				parser.Panic("Interpeter error", "Cant have a selection of type not string!")
			}

			var wrapper any = value.(map[any]any)[EvalOne]
			if ReturnSelectListPointer {
				return &wrapper, []string{"any"}
			}
			return wrapper, []string{"any"}

		}
		parser.Panic("Runtime error", "Cound reach the selection")
	}

	op, ok := CalData.(parser.Oporation)
	if !ok {
		fmt.Println("Cant continue because there is not oporation selected!", CalData)
		os.Exit(1)
	}
	switch op.Op {
	case "+":
		LeftSide, _ := Evaluate(op.Left, indentMap, funcMap, keyFuncs, false)
		RightSide, _ := Evaluate(op.Right, indentMap, funcMap, keyFuncs, false)
		if !slices.Contains(ApproveSideToOp, ReturnType(LeftSide)) && !slices.Contains(ApproveSideToOp, ReturnType(RightSide)) {
			fmt.Println("Cant do None type!", ReturnType(LeftSide), ReturnType(RightSide))
		}

		switch left := LeftSide.(type) {
		case int:
			right, ok := RightSide.(int)
			if ok {
				return left + right, []string{"int"}
			}
			right2, ok2 := RightSide.(float64)
			if ok2 {
				return right2 + float64(left), []string{"float"}
			}
			right3, ok3 := RightSide.(byte)
			if ok3 {
				return right3 + byte(left), []string{"char"}
			}
		case float64:
			right, ok := RightSide.(int)
			if ok {
				return left + float64(right), []string{"float"}
			}
			right2, ok2 := RightSide.(float64)
			if ok2 {
				return right2 + left, []string{"float"}
			}
			fmt.Println("Cant add two incompadeple types!")
			os.Exit(1)

		case string:
			right, ok := RightSide.(string)
			if !ok {
				fmt.Println("Cant add to a string, a none string", RightSide)
			}
			return left + right, []string{"string"}
		case byte:
			right, ok := RightSide.(int)
			if ok {
				return right + int(left), []string{"int"}
			}
			right2, ok2 := RightSide.(byte)
			if ok2 {
				return left + right2, []string{"byte"}
			}
			fmt.Println("Can't add to a byte a incompatible type!")
		}
	case "-":
		LeftSide, _ := Evaluate(op.Left, indentMap, funcMap, keyFuncs, false)
		RightSide, _ := Evaluate(op.Right, indentMap, funcMap, keyFuncs, false)
		if !slices.Contains(ApproveSideToOp, ReturnType(LeftSide)) && !slices.Contains(ApproveSideToOp, ReturnType(RightSide)) {
			fmt.Println("Cant do None type!", ReturnType(LeftSide), ReturnType(RightSide))
		}
		switch left := LeftSide.(type) {
		case int:
			right, ok := RightSide.(int)
			if ok {
				return left - right, []string{"int"}
			}
			right2, ok2 := RightSide.(float64)
			if ok2 {
				return float64(left) - right2, []string{"float"}
			}
			right3, ok3 := RightSide.(byte)
			if ok3 {
				return byte(left) - right3, []string{"byte"}
			}
			fmt.Println("Can't added an incompatible type to an integer!")
			os.Exit(1)
		case float64:
			right, ok := RightSide.(int)
			if ok {
				return left - float64(right), []string{"float"}
			}
			right2, ok2 := RightSide.(float64)
			if ok2 {
				return left - right2, []string{"float"}
			}
			fmt.Println("Cant add two incompadeple types!")
			os.Exit(1)

		case byte:

			right, ok := RightSide.(int)
			if ok {
				return int(left) - right, []string{"int"}
			}
			right2, ok2 := RightSide.(byte)
			if ok2 {
				return left - right2, []string{"byte"}
			}
			fmt.Println("Can't substract to a byte a incompatible type!")
			os.Exit(1)
		}
	case "*":
		LeftSide, _ := Evaluate(op.Left, indentMap, funcMap, keyFuncs, false)
		RightSide, _ := Evaluate(op.Right, indentMap, funcMap, keyFuncs, false)
		switch left := LeftSide.(type) {
		case int:
			integerRight, ok := RightSide.(int)
			if ok {
				return integerRight * left, []string{"int"}
			}
			floatRight, ok := RightSide.(float64)
			if ok {
				return floatRight * float64(left), []string{"float"}
			}
			fmt.Println("Cant add two incompadeple types!")
			os.Exit(1)
		case float64:
			integerRight, ok := RightSide.(int)
			if ok {
				return float64(integerRight) * left, []string{"float"}
			}
			floatRight, ok := RightSide.(float64)
			if ok {
				return floatRight * left, []string{"float"}
			}
			fmt.Println("Cant add two incompadeple types!")
			os.Exit(1)
		case string:
			integerRight, ok := RightSide.(int)
			if ok {
				value := ""
				for i := 0; i < integerRight; i++ {
					value += left
				}
				return value, []string{"string"}
			}
			fmt.Println("Cant add two incompadeple types!")
			os.Exit(1)
		}
	case "/":
		LeftSide, _ := Evaluate(op.Left, indentMap, funcMap, keyFuncs, false)
		RightSide, _ := Evaluate(op.Right, indentMap, funcMap, keyFuncs, false)
		if !slices.Contains(ApproveSideToOp, ReturnType(LeftSide)) && !slices.Contains(ApproveSideToOp, ReturnType(RightSide)) {
			fmt.Println("Cant do None type!", ReturnType(LeftSide), ReturnType(RightSide))
		}
		switch left := LeftSide.(type) {
		case int:
			right, ok := RightSide.(int)
			if ok {
				return float64(left) / float64(right), []string{"float"}
			}
			right2, ok2 := RightSide.(float64)
			if ok2 {
				return float64(left) / float64(right2), []string{"float"}
			}
		case float64:
			right, ok := RightSide.(int)
			if ok {
				return left / float64(right), []string{"float"}
			}
			right2, ok2 := RightSide.(float64)
			if ok2 {
				return left / right2, []string{"float"}
			}
			fmt.Println("Cant add two incompadeple types!")
			os.Exit(1)

		}
	case "==":
		LeftSide, _ := Evaluate(op.Left, indentMap, funcMap, keyFuncs, false)
		RightSide, _ := Evaluate(op.Right, indentMap, funcMap, keyFuncs, false)
		return BoolToInt(LeftSide == RightSide), []string{"int"}
	case "!=":
		LeftSide, _ := Evaluate(op.Left, indentMap, funcMap, keyFuncs, false)
		RightSide, _ := Evaluate(op.Right, indentMap, funcMap, keyFuncs, false)
		return BoolToInt(LeftSide != RightSide), []string{"int"}
	case ">":

		LeftSide, _ := Evaluate(op.Left, indentMap, funcMap, keyFuncs, false)
		RightSide, _ := Evaluate(op.Right, indentMap, funcMap, keyFuncs, false)
		if !slices.Contains(ApproveSideToOp, ReturnType(LeftSide)) && !slices.Contains(ApproveSideToOp, ReturnType(RightSide)) {
			fmt.Println("Cant do None type!", ReturnType(LeftSide), ReturnType(RightSide))
		}

		switch left := LeftSide.(type) {
		case int:
			RightValue, ok := RightSide.(int)
			if ok {
				return BoolToInt(left > RightValue), []string{"int"}
			}

			RightValue2, ok2 := RightSide.(float64)
			if ok2 {
				return BoolToInt(float64(left) > RightValue2), []string{"int"}
			}
		case float64:

			RightValue, ok := RightSide.(int)
			if ok {
				return BoolToInt(left > float64(RightValue)), []string{"int"}
			}
			RightValue2, ok2 := RightSide.(float64)
			if ok2 {
				return BoolToInt(left > RightValue2), []string{"int"}
			}
		}

	}
	return Value, []string{"null"}
}

func ToMethod(Method parser.AccessMethod, Env Environment) *any {
	var Result any = &Method

	// 1. Get the Ident from the map
	mapValue := Env.VariableMap[(*Result.(*parser.AccessMethod)).Object.(lexer.Token).Value]

	// 2. Put it into an 'any' interface variable first
	var interfaceContainer any = mapValue

	// 3. Now you can safely take the pointer to 'any'
	var IdentLocation *any = &interfaceContainer

	Final := false
	for !Final {
		switch mt := Method.Method.(type) {
		case lexer.Token:
			var Wrapper any = (*IdentLocation).(Ident).Value.(StructInstance).Fields[mt.Value]
			return &Wrapper

		case parser.AccessMethod:
			// Note: This assignment will fail if Fields[...], which is likely a *any,
			// is assigned directly to IdentLocation. Let's make sure it matches.
			var Wrapper any = (*IdentLocation).(Ident).Value.(StructInstance).Fields[mt.Method.(lexer.Token).Value]
			IdentLocation = &Wrapper
		}
	}
	return IdentLocation
}

func (Env *Environment) Interpeter() {
	if Env.VariableMap == nil {
		Env.VariableMap = map[string]Ident{}
	}
	if Env.FuncMap == nil {
		Env.FuncMap = map[string]parser.Function{}
	}
	for i := 0; i < len(Env.ParseDate); i++ {
		ParseToken := Env.ParseDate[i]
		// fmt.Println(ParseToken)
		switch ReturnType(ParseToken) {
		case FuncType:
			FuncTemp := ParseToken.(parser.Function)
			Env.FuncMap[FuncTemp.Name] = FuncTemp

		case NewIdentType:
			IdentTemp := ParseToken.(parser.NewIdent)
			if structDef, isStruct := Env.StructMap[IdentTemp.Type]; isStruct {
				instance := StructInstance{
					TypeName: IdentTemp.Type,
					Fields:   make(map[string]*Ident),
				}
				for fieldName, Value := range structDef.Members_Ident {
					var StartVal Ident

					StartVal.Value, _ = GetStartValueType(*Env, Value.Type)
					StartVal.Type = Value.Type
					StartVal.IsConst = false
					StartVal.Name = Value.Name
					instance.Fields[fieldName] = &StartVal // Set default field initializers
				}
				Env.VariableMap[IdentTemp.Name] = Ident{
					Value: instance, Name: IdentTemp.Name, Type: IdentTemp.Type,
				}
				continue
			}
			Result, _ := Evaluate(IdentTemp.Content, Env.VariableMap, Env.FuncMap, Env.Keyfuncs, false)

			NewIdent := Ident{Value: Result, Name: IdentTemp.Name, Type: IdentTemp.Type, IsConst: IdentTemp.IsConst}
			Env.VariableMap[NewIdent.Name] = NewIdent

		case OporationType:
			OporationTemp := ParseToken.(parser.Oporation)
			t, _ := Evaluate(OporationTemp, Env.VariableMap, Env.FuncMap, Env.Keyfuncs, false)
			Env.Output = append(Env.Output, t)

		case ReturnType(lexer.Token{}):
			token := ParseToken.(lexer.Token)
			if token.Type == lexer.NEWLINE {
				continue
			}
			if token.Type == lexer.IDENTIFIER {
				val, ok := Env.VariableMap[token.Value]
				if !ok {
					fmt.Println("Couldn't find veruble names:", token.Value)
					os.Exit(1)
				}
				Env.Output = append(Env.Output, val.Value)

			} else {
				t, _ := Evaluate(ParseToken, Env.VariableMap, Env.FuncMap, Env.Keyfuncs, false)
				Env.Output = append(Env.Output, t)
			}
		case ReturnType(parser.Struct{}):
			StructReg := ParseToken.(parser.Struct)
			Env.StructMap[StructReg.Name] = StructReg
		case ReturnType(parser.CallFunction{}):
			NewSave := SaveDataEnv{}
			NewSave.Save(*Env)
			TempCall := ParseToken.(parser.CallFunction)

			CallFunc, ok := Env.FuncMap[TempCall.Name]
			_, ok2 := Env.Keyfuncs[TempCall.Name]
			if !ok && !ok2 {
				fmt.Println("Error, coudnt find the func:", TempCall.Name)
				os.Exit(1)
			}
			CallVarMap := map[string]Ident{}
			TempValues := []any{}
			for idx, ident := range TempCall.ParimitersInput {
				callEval, _ := Evaluate(ident, Env.VariableMap, Env.FuncMap, Env.Keyfuncs, false)
				if ok2 {
					TempValues = append(TempValues, callEval)
				} else {
					CallVarMap[CallFunc.Perameters[idx].Name] = Ident{Value: callEval, Name: CallFunc.Perameters[idx].Name, Type: CallFunc.Perameters[idx].Type, IsConst: false}
				}

			}

			if funcData, ok := Env.Keyfuncs[TempCall.Name]; ok {
				Result, _ := funcData(TempValues...)
				Env.Output = append(Env.Output, Result)
				continue
			}

			NewEnv := Environment{ParseDate: CallFunc.Body, VariableMap: CallVarMap, FuncMap: Env.FuncMap, Keyfuncs: Env.Keyfuncs}

			NewEnv.Interpeter()
			NewSave.LoadTo(Env)
			Env.Output = append(Env.Output, NewEnv.Output...)

		case ReturnType(parser.Return{}):
			TempReturn := ParseToken.(parser.Return)
			Env.Output = []any{}
			for _, expr := range TempReturn.Exprs {
				ReturnEval, _ := Evaluate(expr, Env.VariableMap, Env.FuncMap, Env.Keyfuncs, false)
				Env.Output = append(Env.Output, ReturnEval)
			}
			Env.Returned = true
			return

		case RefactType:
			TempRefact := ParseToken.(parser.RefactIdent)
			// dont mess: Env.VariableMap[NewIdent.Name] = NewIdent
			switch objectType := TempRefact.Object.(type) {
			case lexer.Token:

				IdentGet, ok := Env.VariableMap[objectType.Value]
				if !ok {
					fmt.Printf("Coudnt find a veruble named: %s\n", objectType.Value)
					os.Exit(1)
				}
				if IdentGet.IsConst {
					fmt.Println("Cant edit a const veruble: ", objectType.Value)
					os.Exit(1)
				}
				NewEnv, Type := Evaluate(TempRefact.Content, Env.VariableMap, Env.FuncMap, Env.Keyfuncs, false)

				var tempIdent Ident = Env.VariableMap[IdentGet.Name]
				if !slices.Equal(Type, []string{tempIdent.Type}) {
					fmt.Println("Cant change a val with incopadble type!")
					os.Exit(1)
				}
				tempIdent.Value = NewEnv
				Env.VariableMap[IdentGet.Name] = tempIdent
			case parser.AccessMethod:

				NewEnv, _ := Evaluate(TempRefact.Content, Env.VariableMap, Env.FuncMap, Env.Keyfuncs, false)
				switch methodGot := (*ToMethod(objectType, *Env)).(type) {
				case *Ident:
					methodGot.Value = NewEnv
				}
			case parser.SectionList:
				newVal, _ := Evaluate(TempRefact.Content, Env.VariableMap, Env.FuncMap, Env.Keyfuncs, false)

				// 2. Resolve the targeted list or dict expression
				targetVal, _ := Evaluate(objectType.List, Env.VariableMap, Env.FuncMap, Env.Keyfuncs, false)

				// 3. Resolve index / key values
				var startIdx any
				if objectType.Start != nil {
					startIdx, _ = Evaluate(objectType.Start, Env.VariableMap, Env.FuncMap, Env.Keyfuncs, false)
				}

				// 4. Update the target collection container
				switch container := targetVal.(type) {
				case map[any]any:
					// Update key-value pair in dictionary
					container[startIdx] = newVal

				case []any:
					// Unpack double-wrapped list structure if nested
					actualSlice := container
					if len(container) == 1 {
						if innerSlice, ok := container[0].([]any); ok {
							actualSlice = innerSlice
						}
					}

					idx, ok := startIdx.(int)
					if !ok {
						// Handle float to int conversions if numbers evaluate as float64
						if fIdx, isFloat := startIdx.(float64); isFloat {
							idx = int(fIdx)
						} else {
							parser.Panic("Runtime Error", "List index must be an integer")
						}
					}

					if idx < 0 || idx >= len(actualSlice) {
						parser.Panic("Runtime Error", "Index out of range")
					}

					// Mutate element in place
					actualSlice[idx] = newVal
				}

			}
		case HouseType:
			TempHouse := ParseToken.(parser.House)
			Names := TempHouse.Names
			Types := []string{}
			Contents := []any{}
			flattenContexts := []any{}
			for _, content := range TempHouse.Contents {

				context, typed := Evaluate(content, Env.VariableMap, Env.FuncMap, Env.Keyfuncs, false)

				Contents = append(Contents, context)

				switch contextType := context.(type) {
				case []any:
					if slices.Compare(typed, []string{"list"}) != 0 {
						flattenContexts = append(flattenContexts, contextType...)
					} else {
						flattenContexts = append(flattenContexts, contextType)
					}

				default:
					flattenContexts = append(flattenContexts, contextType)
				}

				Types = append(Types, typed...)
			}
			if len(flattenContexts) != len(Names) || len(Types) != len(Names) {
				fmt.Println("Not every arg has a corospond arg! The length are:", len(Names), len(flattenContexts), Names, flattenContexts)
				os.Exit(1)
			}
			for index := range len(Names) {
				Env.VariableMap[Names[index]] = Ident{Value: flattenContexts[index], Type: Types[index], Name: Names[index], IsConst: false}
			}
		case IfType:
			var TempIfPointer parser.IfStm = ParseToken.(parser.IfStm)
			for true {
				Values, types := Evaluate(TempIfPointer.Condition, Env.VariableMap, Env.FuncMap, Env.Keyfuncs, false)
				if len(types) != 1 {
					fmt.Println("The condition isnt one typed!")
					os.Exit(1)
				}
				if types[0] != "int" {
					fmt.Println("The condition type isnt an int!")
					os.Exit(1)
				}

				if Values.(int) == 1 {
					NewSave := SaveDataEnv{}
					NewSave.Save(*Env)
					NewEnv := NewEnvironment(TempIfPointer.Body, Env.FuncMap, Env.VariableMap)
					NewEnv.Interpeter()
					if NewEnv.Returned {
						Env.Output = append(Env.Output, NewEnv.Output...)
						return
					}
					NewSave.LoadTo(Env)
					break

				}
				if TempIfPointer.Else == nil {
					break
				}
				TempIfPointer = *TempIfPointer.Else

			}
		case ReachType:

			TempReach := ParseToken.(parser.Reach)
			Path := TempReach.Path
			if IsLink(Path) {
				Data, err := http.Get(Path)
				if err != nil {
					fmt.Println("Coudnt load online file!")
					os.Exit(1)
				}
				defer Data.Body.Close()
				if Data.StatusCode != http.StatusOK {
					fmt.Println("Couldn't load file from the internet! Error code:", Data.StatusCode)
					os.Exit(1)
				}
				data, err2 := io.ReadAll(Data.Body)
				if err2 != nil {
					fmt.Println("Couldn't load file content!", Path)
				}
				NewText := string(data)
				// here is the parsing and tokenizing
				LexerProsses := lexer.Lexer{Input: NewText}
				LexerProsses.LexerAll()
				LexerProsses.AddEOF()
				ParserProsses := parser.Parse(LexerProsses.Tokens)
				Env.ParseDate = slices.Insert(Env.ParseDate, i+1, ParserProsses...)
				ReachImported = append(ReachImported, Path)
				continue
			}
			Cpath, err := os.Getwd()
			if err != nil {
				fmt.Println("Error: coudnt load the current dirrectory!", Path)
				os.Exit(1)
			}
			Path = strings.ReplaceAll(Path, "$", Cpath)
			if slices.Contains(ReachImported, Path) {
				fmt.Println("You dont need to reach (import) the same file!", Path)
				continue
			}
			file, err2 := os.Open(Path)
			if err2 != nil {
				fmt.Println("Codunt load the file!")
				os.Exit(1)
			}
			defer file.Close()
			TextByte, err3 := io.ReadAll(file)
			if err3 != nil {
				fmt.Println("Coudnt load the context from the file!")
				os.Exit(1)
			}
			NewText := string(TextByte)
			// here is the parsing and tokenizing
			LexerProsses := lexer.Lexer{Input: NewText}
			LexerProsses.LexerAll()
			LexerProsses.AddEOF()

			ParserProsses := parser.Parse(LexerProsses.Tokens)
			Env.ParseDate = slices.Insert(Env.ParseDate, i+1, ParserProsses...)
			ReachImported = append(ReachImported, Path)
			continue
		case WhileType:
			NewSave := SaveDataEnv{}
			NewSave.Save(*Env)
			TempWhile := ParseToken.(parser.WhileLoop)
			Data, Types := Evaluate(TempWhile.Condition, Env.VariableMap, Env.FuncMap, Env.Keyfuncs, false)
			if len(Types) > 1 || len(Types) <= 0 || len(Types) == 1 && Types[0] != "int" {
				parser.Panic("Syntax error", fmt.Sprint("Coudnt load condition!", TempWhile.Condition))
			}
			for Data.(int) == 1 {
				// FIX: Pass Env.VariableMap directly so mutated variables persist inside the loop

				NewEnv := NewEnvironment(TempWhile.Body, Env.FuncMap, Env.VariableMap)
				NewEnv.Interpeter()

				// FIX: If return triggered inside while loop, exit the entire block
				if NewEnv.Returned {
					Env.Output = append(Env.Output, NewEnv.Output...)
					Env.Returned = true
					return
				}
				if NewEnv.Breaked {
					Env.Output = append(Env.Output, NewEnv.Output...)
					break
				}
				Data, Types = Evaluate(TempWhile.Condition, Env.VariableMap, Env.FuncMap, Env.Keyfuncs, false)
			}
			NewSave.LoadTo(Env)
		case BreakType:
			Env.Breaked = true
		case SelectList:
			TempSelectList := ParseToken.(parser.SectionList)
			Eval, _ := Evaluate(TempSelectList, Env.VariableMap, Env.FuncMap, Env.Keyfuncs, false)
			Env.Output = append(Env.Output, Eval)
		}
	}
}
