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
		"!": not,
		"~": bitNOT,
	}

	binaryOperatorFns = map[string]doubleOperandFn{
		"+":  add,
		"-":  sub,
		"*":  mul,
		"/":  div,
		"%":  mod,
		"<":  lt,
		">":  gt,
		"<=": le,
		">=": ge,
		"==": eq,
		"!=": neq,
		"&":  bitAND,
		"|":  bitOR,
		"^":  bitXOR,
		"<<": shl,
		">>": shr,
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

func not(obj object.Object) (object.Object, error) {
	return object.NewBoolean(!isTruthy(obj)), nil
}

func bitNOT(obj object.Object) (object.Object, error) {
	return object.NewInteger(^objectToInteger(obj)), nil
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
	res, err := lt(r, l)
	if err != nil {
		return nil, fmt.Errorf("'>' not supported between '%s' and '%s'", l.Type(), r.Type())
	}
	return res, nil
}

func le(l, r object.Object) (object.Object, error) {
	res, err := gt(l, r)
	if err != nil {
		return nil, fmt.Errorf("'<=' not supported between '%s' and '%s'", l.Type(), r.Type())
	}
	return not(res)
}

func ge(l, r object.Object) (object.Object, error) {
	res, err := lt(l, r)
	if err != nil {
		return nil, fmt.Errorf("'>=' not supported between '%s' and '%s'", l.Type(), r.Type())
	}
	return not(res)
}

func eq(l, r object.Object) (object.Object, error) {
	if l.Type() == object.NULL_OBJ && r.Type() == object.NULL_OBJ {
		return object.TRUE, nil
	}

	defaultErr := fmt.Errorf("'==' not supported between '%s' and '%s'", l.Type(), r.Type())
	if res, err := lt(l, r); err != nil {
		return nil, defaultErr
	} else if res == object.TRUE {
		return object.FALSE, nil
	}

	if res, err := gt(l, r); err != nil {
		return nil, defaultErr
	} else {
		return not(res)
	}
}

func neq(l, r object.Object) (object.Object, error) {
	res, err := eq(l, r)
	if err != nil {
		return nil, fmt.Errorf("'!=' not supported between '%s' and '%s'", l.Type(), r.Type())
	}

	return not(res)
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

func mod(l, r object.Object) (object.Object, error) {
	switch l.Type() {
	case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
		switch r.Type() {
		case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
			modulus := objectToInteger(r)
			if modulus == 0 {
				return nil, errors.New("integer division or modulo by zero")
			}
			return object.NewInteger(((objectToInteger(l) % modulus) + modulus) % modulus), nil
		}
	}

	return nil, fmt.Errorf("'mod' not supported between '%s' and '%s'", l.Type(), r.Type())
}

func bitAND(l, r object.Object) (object.Object, error) {
	switch l.Type() {
	case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
		switch r.Type() {
		case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
			return object.NewInteger(objectToInteger(l) & objectToInteger(r)), nil
		}
	}

	return nil, fmt.Errorf("'&' not supported between '%s' and '%s'", l.Type(), r.Type())
}

func bitOR(l, r object.Object) (object.Object, error) {
	switch l.Type() {
	case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
		switch r.Type() {
		case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
			return object.NewInteger(objectToInteger(l) | objectToInteger(r)), nil
		}
	}

	return nil, fmt.Errorf("'|' not supported between '%s' and '%s'", l.Type(), r.Type())
}

func bitXOR(l, r object.Object) (object.Object, error) {
	switch l.Type() {
	case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
		switch r.Type() {
		case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
			return object.NewInteger(objectToInteger(l) ^ objectToInteger(r)), nil
		}
	}

	return nil, fmt.Errorf("'^' not supported between '%s' and '%s'", l.Type(), r.Type())
}

func shl(l, r object.Object) (object.Object, error) {
	switch l.Type() {
	case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
		switch r.Type() {
		case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
			shift := objectToInteger(r)
			if shift < 0 {
				return nil, errors.New("negative shift amount")
			}
			return object.NewInteger(objectToInteger(l) << shift), nil
		}
	}

	return nil, fmt.Errorf("'<<' not supported between '%s' and '%s'", l.Type(), r.Type())
}

func shr(l, r object.Object) (object.Object, error) {
	switch l.Type() {
	case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
		switch r.Type() {
		case object.INTEGER_OBJ, object.BOOLEAN_OBJ:
			shift := objectToInteger(r)
			if shift < 0 {
				return nil, errors.New("negative shift amount")
			}
			return object.NewInteger(objectToInteger(l) >> shift), nil
		}
	}

	return nil, fmt.Errorf("'>>' not supported between '%s' and '%s'", l.Type(), r.Type())
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
