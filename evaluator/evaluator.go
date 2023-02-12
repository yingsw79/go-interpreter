package evaluator

import (
	"errors"
	"fmt"
	"go-interpreter/ast"
	"go-interpreter/object"
)

func Eval(node ast.Node) (object.Object, error) {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)

	case *ast.BlockStatement:
		return evalBlockStatement(node)

	case *ast.ExpressionStatement:
		return evalExpressionStatement(node)

	case *ast.LetStatement:
		return evalLetStatement(node)

	case *ast.ReturnStatement:
		return evalReturnStatement(node)

	case *ast.Identifier:
		return evalIdentifier(node)

	case *ast.IntegerLiteral:
		return evalIntegerLiteral(node)

	case *ast.Boolean:
		return evalBoolean(node)

	case *ast.PrefixExpression:
		return evalPrefixExpression(node)

	case *ast.InfixExpression:
		return evalInfixExpression(node)

	case *ast.IfExpression:
		return evalIfExpression(node)

	default:
		return nil, errors.New("invalid syntax")
	}
}

func evalProgram(p *ast.Program) (res object.Object, err error) {
	for _, statement := range p.Statements {
		if res, err = Eval(statement); err != nil {
			return
		}

		if returnValue, ok := res.(*object.ReturnValue); ok {
			res = returnValue.Value
			return
		}
	}

	return
}

func evalBlockStatement(bs *ast.BlockStatement) (res object.Object, err error) {
	for _, statement := range bs.Statements {
		if res, err = Eval(statement); err != nil || res.Type() == object.RETURN_VALUE_OBJ {
			return
		}
	}

	return
}

func evalExpressionStatement(es *ast.ExpressionStatement) (object.Object, error) {
	return Eval(es.Expression)
}

func evalLetStatement(ls *ast.LetStatement) (object.Object, error) {
	val, err := Eval(ls.Value)
	if err != nil {
		return nil, err
	}

	object.ENV.Set(ls.Name.Value, val)

	return nil, nil
}

func evalReturnStatement(rs *ast.ReturnStatement) (object.Object, error) {
	val, err := Eval(rs.ReturnValue)
	if err != nil {
		return nil, err
	}

	return object.NewReturnValue(val), nil
}

func evalIdentifier(ident *ast.Identifier) (object.Object, error) {
	val, ok := object.ENV.Get(ident.Value)
	if !ok {
		return nil, fmt.Errorf("name %q is not defined", ident.Value)
	}

	return val, nil
}

func evalIntegerLiteral(il *ast.IntegerLiteral) (object.Object, error) {
	return object.NewInteger(il.Value), nil
}

func evalBoolean(b *ast.Boolean) (object.Object, error) {
	return object.NewBoolean(b.Value), nil
}

func evalPrefixExpression(pe *ast.PrefixExpression) (object.Object, error) {
	val, err := Eval(pe.Right)
	if err != nil {
		return nil, err
	}

	switch pe.Operator {
	case "!":
		return evalBangOperatorExpression(val)
	case "-":
		return evalMinusPrefixOperatorExpression(val)
	default:
		return nil, fmt.Errorf("unknown operator: %q", pe.Operator)
	}
}

func evalBangOperatorExpression(obj object.Object) (object.Object, error) {
	switch obj := obj.(type) {
	case *object.Boolean:
		return object.NewBoolean(!obj.Value), nil
	case *object.Integer:
		return object.NewBoolean(obj.Value == 0), nil
	case *object.Null:
		return object.TRUE, nil
	case *object.ReturnValue:
		return nil, fmt.Errorf("bad operand type for unary !: %q", obj.Type())
	default:
		return object.FALSE, nil
	}
}

func evalMinusPrefixOperatorExpression(obj object.Object) (object.Object, error) {
	if obj.Type() == object.INTEGER_OBJ || obj.Type() == object.BOOLEAN_OBJ {
		return object.NewInteger(-objectToInteger(obj)), nil
	}

	return nil, fmt.Errorf("bad operand type for unary -: %q", obj.Type())
}

func evalInfixExpression(ie *ast.InfixExpression) (object.Object, error) {
	l, err := Eval(ie.Left)
	if err != nil {
		return nil, err
	}

	r, err := Eval(ie.Right)
	if err != nil {
		return nil, err
	}

	if l.Type() == object.NULL_OBJ && r.Type() == object.NULL_OBJ {
		switch ie.Operator {
		case "==":
			return object.TRUE, nil
		case "!=":
			return object.FALSE, nil
		}
	} else if (l.Type() == object.INTEGER_OBJ || l.Type() == object.BOOLEAN_OBJ) &&
		(r.Type() == object.INTEGER_OBJ || r.Type() == object.BOOLEAN_OBJ) {
		return evalIntegerInfixExpression(ie.Operator, l, r)
	}

	return nil, fmt.Errorf("%q not supported between %q and %q", ie.Operator, l.Type(), r.Type())
}

func evalIntegerInfixExpression(operator string, l, r object.Object) (object.Object, error) {
	lv, rv := objectToInteger(l), objectToInteger(r)

	switch operator {
	case "+":
		return object.NewInteger(lv + rv), nil
	case "-":
		return object.NewInteger(lv - rv), nil
	case "*":
		return object.NewInteger(lv * rv), nil
	case "/":
		if rv == 0 {
			return nil, errors.New("division by zero")
		}

		return object.NewInteger(lv / rv), nil
	case "<":
		return object.NewBoolean(lv < rv), nil
	case ">":
		return object.NewBoolean(lv > rv), nil
	case "==":
		return object.NewBoolean(lv == rv), nil
	case "!=":
		return object.NewBoolean(lv != rv), nil
	default:
		return nil, fmt.Errorf("unknown operator: %q", operator)
	}
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

func evalIfExpression(ie *ast.IfExpression) (object.Object, error) {
	condition, err := Eval(ie.Condition)
	if err != nil {
		return nil, err
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	}

	return object.NULL, nil
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Integer:
		return obj.Value != 0
	case *object.Null:
		return false
	default:
		return true
	}
}
