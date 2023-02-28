package evaluator

import (
	"errors"
	"fmt"
	"go-interpreter/object"
	"os"
	"sort"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) (object.Object, error) {
			if err := checkArgsLen(len(args), 1); err != nil {
				return nil, err
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
			if err := checkArgsLen(len(args), 0); err != nil {
				return nil, err
			}

			os.Exit(0)

			return object.NULL, nil
		},
	},
	"push": {
		Fn: func(args ...object.Object) (object.Object, error) {
			if len(args) == 0 {
				return nil, errors.New("wrong number of arguments: got=0, want>0")
			}

			arr, err := checkIsArray("push", args[0])
			if err != nil {
				return nil, err
			}

			arr.Elements = append(arr.Elements, args[1:]...)
			return arr, nil
		},
	},
	"pop": {
		Fn: func(args ...object.Object) (object.Object, error) {
			if err := checkArgsLen(len(args), 1); err != nil {
				return nil, err
			}

			arr, err := checkIsArray("pop", args[0])
			if err != nil {
				return nil, err
			}

			n := len(arr.Elements)
			if n == 0 {
				return nil, errors.New("pop from empty array")
			}

			res := arr.Elements[n-1]
			arr.Elements = arr.Elements[:n-1]
			return res, nil
		},
	},
	"reverse": {
		Fn: func(args ...object.Object) (object.Object, error) {
			if err := checkArgsLen(len(args), 1); err != nil {
				return nil, err
			}

			arr, err := checkIsArray("reverse", args[0])
			if err != nil {
				return nil, err
			}

			for i, j := 0, len(arr.Elements)-1; i < j; i, j = i+1, j-1 {
				arr.Elements[i], arr.Elements[j] = arr.Elements[j], arr.Elements[i]
			}

			return arr, nil
		},
	},
	"sort": {
		Fn: func(args ...object.Object) (object.Object, error) {
			if err := checkArgsLen(len(args), 1); err != nil {
				return nil, err
			}

			arr, err := checkIsArray("sort", args[0])
			if err != nil {
				return nil, err
			}

			if len(arr.Elements) > 1 {
				t := arr.Elements[0].Type()

				switch t {
				case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
					for _, o := range arr.Elements {
						if o.Type() != object.INTEGER_OBJ && o.Type() != object.BOOLEAN_OBJ {
							return nil, fmt.Errorf("'<' not supported between '%s' and '%s'", t, o.Type())
						}
					}

					sort.Slice(arr.Elements, func(i, j int) bool {
						return objectToInteger(arr.Elements[i]) < objectToInteger(arr.Elements[j])
					})

				case object.STRING_OBJ:
					for _, o := range arr.Elements {
						if o.Type() != t {
							return nil, fmt.Errorf("'<' not supported between '%s' and '%s'", t, o.Type())
						}
					}

					sort.Slice(arr.Elements, func(i, j int) bool {
						return arr.Elements[i].(*object.String).Value < arr.Elements[j].(*object.String).Value
					})

				default:
					return nil, fmt.Errorf("'%s' is incomparable", t)
				}
			}

			return arr, nil
		},
	},
}

func checkArgsLen(got, want int) error {
	if got != want {
		return fmt.Errorf("wrong number of arguments: got=%d, want=%d", got, want)
	}
	return nil
}

func checkIsArray(fn string, obj object.Object) (*object.Array, error) {
	arr, ok := obj.(*object.Array)
	if !ok {
		return nil, fmt.Errorf("argument to '%s' must be 'ARRAY', got '%s'", fn, obj.Type())
	}
	return arr, nil
}
