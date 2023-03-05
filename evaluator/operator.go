package evaluator

import (
	"errors"
	"fmt"
	"go-interpreter/object"
	"strings"
)

var (
	unaryOperatorFns = map[string]singleOperandFn{
		"-": neg,
		"!": bang,
	}

	binaryOperatorFns = map[string]doubleOperandFn{
		"+":  add,
		"-":  sub,
		"*":  mul,
		"/":  div,
		"<":  lt,
		">":  gt,
		"==": eq,
		"!=": neq,
		"=":  assign,
	}
)

// unary operator
func neg(obj object.Object) (object.Object, error) {
	switch obj.Type() {
	case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
		return object.NewInteger(-objectToInteger(obj)), nil
	default:
		return nil, fmt.Errorf("bad operand type for unary -: '%s'", obj.Type())
	}
}

func bang(obj object.Object) (object.Object, error) {
	return object.NewBoolean(!isTruthy(obj)), nil
}

// binary operator
func lt(l, r object.Object) (object.Object, error) {
	switch l.Type() {
	case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
		switch r.Type() {
		case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
			return object.NewBoolean(objectToInteger(l) < objectToInteger(r)), nil
		}

	case object.STRING_OBJ:
		switch r.Type() {
		case object.STRING_OBJ:
			return object.NewBoolean(l.(*object.String).Value < r.(*object.String).Value), nil
		}
	}

	return nil, fmt.Errorf("'<' not supported between '%s' and '%s'", l.Type(), r.Type())
}

func gt(l, r object.Object) (object.Object, error) {
	switch l.Type() {
	case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
		switch r.Type() {
		case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
			return object.NewBoolean(objectToInteger(l) > objectToInteger(r)), nil
		}

	case object.STRING_OBJ:
		switch r.Type() {
		case object.STRING_OBJ:
			return object.NewBoolean(l.(*object.String).Value > r.(*object.String).Value), nil
		}
	}

	return nil, fmt.Errorf("'>' not supported between '%s' and '%s'", l.Type(), r.Type())
}

func eq(l, r object.Object) (object.Object, error) {
	switch l.Type() {
	case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
		switch r.Type() {
		case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
			return object.NewBoolean(objectToInteger(l) == objectToInteger(r)), nil
		}

	case object.STRING_OBJ:
		switch r.Type() {
		case object.STRING_OBJ:
			return object.NewBoolean(l.(*object.String).Value == r.(*object.String).Value), nil
		}

	case object.NULL_OBJ:
		switch r.Type() {
		case object.NULL_OBJ:
			return object.TRUE, nil
		}
	}

	return nil, fmt.Errorf("'==' not supported between '%s' and '%s'", l.Type(), r.Type())
}

func neq(l, r object.Object) (object.Object, error) {
	switch l.Type() {
	case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
		switch r.Type() {
		case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
			return object.NewBoolean(objectToInteger(l) != objectToInteger(r)), nil
		}

	case object.STRING_OBJ:
		switch r.Type() {
		case object.STRING_OBJ:
			return object.NewBoolean(l.(*object.String).Value != r.(*object.String).Value), nil
		}

	case object.NULL_OBJ:
		switch r.Type() {
		case object.NULL_OBJ:
			return object.FALSE, nil
		}
	}

	return nil, fmt.Errorf("'!=' not supported between '%s' and '%s'", l.Type(), r.Type())
}

func add(l, r object.Object) (object.Object, error) {
	switch l.Type() {
	case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
		switch r.Type() {
		case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
			return object.NewInteger(objectToInteger(l) + objectToInteger(r)), nil
		}

	case object.STRING_OBJ:
		switch r.Type() {
		case object.STRING_OBJ:
			return object.NewString(l.(*object.String).Value + r.(*object.String).Value), nil
		}

	case object.ARRAY_OBJ:
		newObjs := append([]object.Object(nil), l.(*object.Array).Elements...)

		switch r.Type() {
		case object.ARRAY_OBJ:
			return object.NewArray(append(newObjs, r.(*object.Array).Elements...)), nil
		default:
			return object.NewArray(append(newObjs, r)), nil
		}
	}

	if r.Type() == object.ARRAY_OBJ {
		return object.NewArray(append([]object.Object{l}, r.(*object.Array).Elements...)), nil
	}

	return nil, fmt.Errorf("'+' not supported between '%s' and '%s'", l.Type(), r.Type())
}

func sub(l, r object.Object) (object.Object, error) {
	switch l.Type() {
	case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
		switch r.Type() {
		case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
			return object.NewInteger(objectToInteger(l) - objectToInteger(r)), nil
		}
	}

	return nil, fmt.Errorf("'-' not supported between '%s' and '%s'", l.Type(), r.Type())
}

func mul(l, r object.Object) (object.Object, error) {
	switch l.Type() {
	case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
		switch r.Type() {
		case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
			return object.NewInteger(objectToInteger(l) * objectToInteger(r)), nil
		}

	case object.STRING_OBJ:
		switch r.Type() {
		case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
			var b strings.Builder
			s := l.(*object.String).Value
			for i := 0; i < int(objectToInteger(r)); i++ {
				b.WriteString(s)
			}
			return object.NewString(b.String()), nil
		}

	case object.ARRAY_OBJ:
		switch r.Type() {
		case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
			objs := l.(*object.Array).Elements
			newObjs := []object.Object{}
			for i := 0; i < int(objectToInteger(r)); i++ {
				newObjs = append(newObjs, objs...)
			}
			return object.NewArray(newObjs), nil
		}
	}

	return nil, fmt.Errorf("'*' not supported between '%s' and '%s'", l.Type(), r.Type())
}

func div(l, r object.Object) (object.Object, error) {
	switch l.Type() {
	case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
		switch r.Type() {
		case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
			divisor := objectToInteger(r)
			if divisor == 0 {
				return nil, errors.New("division by zero")
			}
			return object.NewInteger(objectToInteger(l) / divisor), nil
		}
	}

	return nil, fmt.Errorf("'/' not supported between '%s' and '%s'", l.Type(), r.Type())
}

func assign(l, r object.Object) (object.Object, error) {
	ident, ok := l.(*object.Identifier)
	if !ok {
		return nil, errors.New("cannot assign to literal")
	}

	ident.Set(ident.Name, r)

	return nil, nil
}

func objectToInteger(obj object.Object) int64 {
	switch obj := obj.(type) {
	case *object.Integer:
		return obj.Value
	case *object.Boolean:
		if obj.Value {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Integer:
		return obj.Value != 0
	case *object.String:
		return obj.Value != ""
	case *object.Array:
		return len(obj.Elements) != 0
	case *object.Null:
		return false
	default:
		return true
	}
}
