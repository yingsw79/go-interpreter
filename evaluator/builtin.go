package evaluator

import (
	"fmt"
	"go-interpreter/object"
	"os"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) (object.Object, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("wrong number of arguments: got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return object.NewInteger(int64(len(arg.Value))), nil
			case *object.Array:
				return object.NewInteger(int64(len(arg.Elements))), nil
			default:
				return nil, fmt.Errorf("argument to 'len' not supported, got '%s'", args[0].Type())
			}
		},
	},
	"exit": {
		Fn: func(args ...object.Object) (object.Object, error) {
			if len(args) != 0 {
				return nil, fmt.Errorf("wrong number of arguments: got=%d, want=0", len(args))
			}

			os.Exit(0)

			return object.NULL, nil
		},
	},
}
