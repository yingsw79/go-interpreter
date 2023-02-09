package evaluator

import (
	"go-interpreter/ast"
	"go-interpreter/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return object.NewInteger(node.Value)
	case *ast.Boolean:
		return object.NewBoolean(node.Value)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node.Operator, Eval(node.Right))
	case *ast.InfixExpression:
		return evalInfixExpression(node.Operator, Eval(node.Left), Eval(node.Right))
	}
	return nil
}

func evalStatements(stmts []ast.Statement) (res object.Object) {
	for _, statement := range stmts {
		res = Eval(statement)
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
	switch r := right.(type) {
	case *object.Boolean:
		return object.NewBoolean(!r.Value)
	case *object.Null:
		return object.TRUE
	case *object.Integer:
		if r.Value == 0 {
			return object.TRUE
		}
		return object.FALSE
	default:
		return object.FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return object.NULL
	}
	return object.NewInteger(-right.(*object.Integer).Value)
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return object.NewBoolean(left == right)
	case operator == "!=":
		return object.NewBoolean(left != right)
	default:
		return object.NULL
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return object.NewInteger(leftVal + rightVal)
	case "-":
		return object.NewInteger(leftVal - rightVal)
	case "*":
		return object.NewInteger(leftVal * rightVal)
	case "/":
		return object.NewInteger(leftVal / rightVal)
	case "<":
		return object.NewBoolean(leftVal < rightVal)
	case ">":
		return object.NewBoolean(leftVal > rightVal)
	case "==":
		return object.NewBoolean(leftVal == rightVal)
	case "!=":
		return object.NewBoolean(leftVal != rightVal)
	default:
		return object.NULL
	}
}
