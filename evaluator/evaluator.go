package evaluator

import (
	"go-interpreter/ast"
	"go-interpreter/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements)
	case *ast.BlockStatement:
		return evalBlockStatement(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return object.NewInteger(node.Value)
	case *ast.Boolean:
		return object.NewBoolean(node.Value)
	case *ast.ReturnStatement:
		return object.NewReturnValue(Eval(node.ReturnValue))
	case *ast.PrefixExpression:
		return evalPrefixExpression(node.Operator, Eval(node.Right))
	case *ast.InfixExpression:
		return evalInfixExpression(node.Operator, Eval(node.Left), Eval(node.Right))
	case *ast.IfExpression:
		return evalIfExpression(node)
	default:
		return nil
	}
}

func evalProgram(stmts []ast.Statement) (res object.Object) {
	for _, statement := range stmts {
		res = Eval(statement)

		if returnValue, ok := res.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}

	return
}

func evalBlockStatement(stmts []ast.Statement) (res object.Object) {
	for _, statement := range stmts {
		res = Eval(statement)

		if res != nil && res.Type() == object.RETURN_VALUE_OBJ {
			return
		}
	}

	return
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return object.NULL
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch t := right.(type) {
	case *object.Boolean:
		return object.NewBoolean(!t.Value)
	case *object.Integer:
		return object.NewBoolean(t.Value == 0)
	case *object.Null:
		return object.TRUE
	default:
		return object.FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() == object.INTEGER_OBJ || right.Type() == object.BOOLEAN_OBJ {
		return object.NewInteger(-objectToInteger(right))
	}

	return object.NULL
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	if left.Type() == object.NULL_OBJ && right.Type() == object.NULL_OBJ {
		switch operator {
		case "==":
			return object.TRUE
		case "!=":
			return object.FALSE
		default:
			return object.NULL
		}
	} else if (left.Type() == object.INTEGER_OBJ || left.Type() == object.BOOLEAN_OBJ) &&
		(right.Type() == object.INTEGER_OBJ || right.Type() == object.BOOLEAN_OBJ) {
		return evalIntegerInfixExpression(operator, objectToInteger(left), objectToInteger(right))
	}

	return object.NULL
}

func evalIntegerInfixExpression(operator string, lv, rv int64) object.Object {
	switch operator {
	case "+":
		return object.NewInteger(lv + rv)
	case "-":
		return object.NewInteger(lv - rv)
	case "*":
		return object.NewInteger(lv * rv)
	case "/":
		return object.NewInteger(lv / rv)
	case "<":
		return object.NewBoolean(lv < rv)
	case ">":
		return object.NewBoolean(lv > rv)
	case "==":
		return object.NewBoolean(lv == rv)
	case "!=":
		return object.NewBoolean(lv != rv)
	default:
		return object.NULL
	}
}

func objectToInteger(obj object.Object) int64 {
	switch t := obj.(type) {
	case *object.Integer:
		return t.Value
	case *object.Boolean:
		if t.Value {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	}

	return object.NULL
}

func isTruthy(obj object.Object) bool {
	switch t := obj.(type) {
	case *object.Boolean:
		return t.Value
	case *object.Integer:
		return t.Value != 0
	case *object.Null:
		return false
	default:
		return true
	}
}
