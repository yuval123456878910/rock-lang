package interpeter

import (
	"rocks/parser"
	"strconv"
)

var convLibEnv Environment = Environment{
	Keyfuncs: map[string]Keyfunc{
		"to_string": func(args ...any) (any, []string) {
			Data := args[0]
			switch d := Data.(type) {
			case byte:
				return string(d), []string{"string"}
			case int:
				return strconv.Itoa(d), []string{"string"}
			case float64:
				return strconv.FormatFloat(d, 'f', -1, 64), []string{"string"}
			}
			return "", []string{"string"}
		},
		"to_char": func(args ...any) (any, []string) {
			stringText, ok := args[0].(string)
			if !ok {
				parser.Panic("Runtime error", "Counl't not translate between 'type' to byte")
			}
			if len(stringText) <= 0 {
				parser.Panic("Runtime error", "Byte length is less than 1")
			}
			return stringText[0], []string{"char"}
		},
		"to_int": func(args ...any) (any, []string) {
			switch Typed := args[0].(type) {
			case string:
				i, err := strconv.Atoi(Typed)
				if err != nil {
					parser.Panic("Runtime error", "couldn't translate between string to int")
				}
				return i, []string{"int"}
			case float64:
				return int(Typed), []string{"int"}
			}
			return 0, []string{"int"}
		},
	},
}
