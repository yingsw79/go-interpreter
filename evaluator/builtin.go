package evaluator

import (
	"errors"
	"fmt"
	"go-interpreter/object"
	"os"
	"sort"
)

type (
	singleOperandFn func(object.Object) (object.Object, error)
	doubleOperandFn func(object.Object, object.Object) (object.Object, error)
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
	"append": {
		Fn: func(args ...object.Object) (object.Object, error) {
			if len(args) == 0 {
				return nil, errors.New("wrong number of arguments: got=0, want>0")
			}

			arr, err := checkIsArray("append", args[0])
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
				for _, o := range arr.Elements[1:] {
					if _, err := lt(arr.Elements[0], o); err != nil {
						return nil, err
					}
				}

				sort.SliceStable(arr.Elements, func(i, j int) bool {
					res, _ := lt(arr.Elements[i], arr.Elements[j])
					return res.(*object.Boolean).Value
				})
			}

			return arr, nil
		},
	},
	"sum": {
		Fn: func(args ...object.Object) (object.Object, error) {
			if len(args) == 0 || len(args) > 2 {
				return nil, fmt.Errorf("wrong number of arguments: got=%d, want=1 or 2", len(args))
			}

			arr, err := checkIsArray("sum", args[0])
			if err != nil {
				return nil, err
			}

			var initializer object.Object
			if len(args) == 2 {
				initializer = args[1]
			}

			if len(arr.Elements) == 0 && initializer == nil {
				return object.NewInteger(0), nil
			}

			res, err := reduce(add, arr.Elements, initializer)
			if err != nil {
				return nil, err
			}

			return res, nil
		},
	},
}

func reduce(fn doubleOperandFn, arr []object.Object, initializer object.Object) (res object.Object, err error) {
	res = initializer

	if res == nil {
		if len(arr) == 0 {
			return nil, fmt.Errorf("reduce() of empty sequence with no initial value")
		}
		res = arr[0]
	} else if len(arr) == 0 {
		return
	} else {
		res, err = fn(res, arr[0])
		if err != nil {
			return nil, err
		}
	}

	for _, o := range arr[1:] {
		res, err = fn(res, o)
		if err != nil {
			return nil, err
		}
	}
	return
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
