package evaluator

import (
	"errors"
	"fmt"
	"go-interpreter/ast"
	"go-interpreter/object"
)

func Eval(node ast.Node, env *object.Environment) (object.Object, error) {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ExpressionStatement:
		return evalExpressionStatement(node, env)

	case *ast.LetStatement:
		return evalLetStatement(node, env)

	case *ast.ReturnStatement:
		return evalReturnStatement(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.IntegerLiteral:
		return evalIntegerLiteral(node)

	case *ast.Boolean:
		return evalBoolean(node)

	case *ast.StringLiteral:
		return evalStringLiteral(node)

	case *ast.PrefixExpression:
		return evalPrefixExpression(node, env)

	case *ast.InfixExpression:
		return evalInfixExpression(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.FunctionLiteral:
		return evalFunctionLiteral(node, env)

	case *ast.CallExpression:
		return evalCallExpression(node, env)

	case *ast.ArrayLiteral:
		return evalArrayLiteral(node, env)

	case *ast.IndexExpression:
		return evalIndexExpression(node, env)

	default:
		return nil, errors.New("invalid syntax")
	}
}

func evalProgram(p *ast.Program, env *object.Environment) (res object.Object, err error) {
	for _, statement := range p.Statements {
		if res, err = Eval(statement, env); err != nil {
			return
		}

		if returnValue, ok := res.(*object.ReturnValue); ok {
			res = returnValue.Value
			return
		}
	}

	return
}

func evalBlockStatement(bs *ast.BlockStatement, env *object.Environment) (res object.Object, err error) {
	for _, statement := range bs.Statements {
		if res, err = Eval(statement, env); err != nil {
			return
		}

		if res != nil && res.Type() == object.RETURN_VALUE_OBJ {
			return
		}
	}

	return
}

func evalExpressionStatement(es *ast.ExpressionStatement, env *object.Environment) (object.Object, error) {
	return Eval(es.Expression, env)
}

func evalLetStatement(ls *ast.LetStatement, env *object.Environment) (object.Object, error) {
	val, err := Eval(ls.Value, env)
	if err != nil {
		return nil, err
	}

	env.Set(ls.Name.Value, val)

	return nil, nil
}

func evalReturnStatement(rs *ast.ReturnStatement, env *object.Environment) (object.Object, error) {
	val, err := Eval(rs.ReturnValue, env)
	if err != nil {
		return nil, err
	}

	return object.NewReturnValue(val), nil
}

func evalIdentifier(ident *ast.Identifier, env *object.Environment) (object.Object, error) {
	if val, ok := env.Get(ident.Value); ok {
		return val, nil
	}

	if builtin, ok := builtins[ident.Value]; ok {
		return builtin, nil
	}

	return nil, fmt.Errorf("name '%s' is not defined", ident.Value)
}

func evalStringLiteral(s *ast.StringLiteral) (object.Object, error) {
	return object.NewString(s.Value), nil
}

func evalIntegerLiteral(il *ast.IntegerLiteral) (object.Object, error) {
	return object.NewInteger(il.Value), nil
}

func evalBoolean(b *ast.Boolean) (object.Object, error) {
	return object.NewBoolean(b.Value), nil
}

func evalPrefixExpression(pe *ast.PrefixExpression, env *object.Environment) (object.Object, error) {
	val, err := Eval(pe.Right, env)
	if err != nil {
		return nil, err
	}

	switch pe.Operator {
	case "!":
		return evalBangOperatorExpression(val)
	case "-":
		return evalMinusPrefixOperatorExpression(val)
	default:
		return nil, fmt.Errorf("unknown operator: '%s'", pe.Operator)
	}
}

func evalBangOperatorExpression(obj object.Object) (object.Object, error) {
	switch obj := obj.(type) {
	case *object.Boolean:
		return object.NewBoolean(!obj.Value), nil
	case *object.Integer:
		return object.NewBoolean(obj.Value == 0), nil
	case *object.String:
		return object.NewBoolean(obj.Value == ""), nil
	case *object.Null:
		return object.TRUE, nil
	default:
		return object.FALSE, nil
	}
}

func evalMinusPrefixOperatorExpression(obj object.Object) (object.Object, error) {
	if obj.Type() == object.INTEGER_OBJ || obj.Type() == object.BOOLEAN_OBJ {
		return object.NewInteger(-objectToInteger(obj)), nil
	}

	return nil, fmt.Errorf("bad operand type for unary -: '%s'", obj.Type())
}

func evalInfixExpression(ie *ast.InfixExpression, env *object.Environment) (object.Object, error) {
	l, err := Eval(ie.Left, env)
	if err != nil {
		return nil, err
	}

	r, err := Eval(ie.Right, env)
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

	} else if l.Type() == object.STRING_OBJ && r.Type() == object.STRING_OBJ {
		return evalStringInfixExpression(ie.Operator, l, r)
	}

	return nil, fmt.Errorf("'%s' not supported between '%s' and '%s'", ie.Operator, l.Type(), r.Type())
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
		return nil, fmt.Errorf("'%s' not supported between '%s' and '%s'", operator, l.Type(), r.Type())
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

func evalStringInfixExpression(operator string, l, r object.Object) (object.Object, error) {
	lv, rv := l.(*object.String).Value, r.(*object.String).Value

	switch operator {
	case "+":
		return object.NewString(lv + rv), nil
	case "<":
		return object.NewBoolean(lv < rv), nil
	case ">":
		return object.NewBoolean(lv > rv), nil
	case "==":
		return object.NewBoolean(lv == rv), nil
	case "!=":
		return object.NewBoolean(lv != rv), nil
	default:
		return nil, fmt.Errorf("'%s' not supported between '%s' and '%s'", operator, l.Type(), r.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) (object.Object, error) {
	condition, err := Eval(ie.Condition, env)
	if err != nil {
		return nil, err
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}

	return object.NULL, nil
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Integer:
		return obj.Value != 0
	case *object.String:
		return obj.Value != ""
	case *object.Null:
		return false
	default:
		return true
	}
}

func evalFunctionLiteral(fl *ast.FunctionLiteral, env *object.Environment) (object.Object, error) {
	return object.NewFunction(fl.Parameters, fl.Body, env), nil
}

func evalCallExpression(ce *ast.CallExpression, env *object.Environment) (object.Object, error) {
	fn, err := Eval(ce.Function, env)
	if err != nil {
		return nil, err
	}

	args, err := evalExpressions(ce.Arguments, env)
	if err != nil {
		return nil, err
	}

	return applyFunction(fn, args)
}

func evalExpressions(exps []ast.Expression, env *object.Environment) ([]object.Object, error) {
	res := []object.Object{}

	for _, e := range exps {
		evaluated, err := Eval(e, env)
		if err != nil {
			return nil, err
		}
		res = append(res, evaluated)
	}

	return res, nil
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for i, param := range fn.Parameters {
		env.Set(param.Value, args[i])
	}

	return env
}

func applyFunction(obj object.Object, args []object.Object) (object.Object, error) {
	switch fn := obj.(type) {
	case *object.Function:
		extendEnv := extendFunctionEnv(fn, args)
		res, err := Eval(fn.Body, extendEnv)
		if err != nil {
			return nil, err
		}

		if returnValue, ok := res.(*object.ReturnValue); ok {
			res = returnValue.Value
		}

		return res, nil

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return nil, fmt.Errorf("not a function: '%s'", obj.Type())
	}
}

func evalArrayLiteral(al *ast.ArrayLiteral, env *object.Environment) (object.Object, error) {
	elements, err := evalExpressions(al.Elements, env)
	if err != nil {
		return nil, err
	}

	return object.NewArray(elements), nil
}

func evalIndexExpression(ie *ast.IndexExpression, env *object.Environment) (object.Object, error) {
	l, err := Eval(ie.Left, env)
	if err != nil {
		return nil, err
	}

	idx, err := Eval(ie.Index, env)
	if err != nil {
		return nil, err
	}

	switch l.Type() {
	case object.ARRAY_OBJ:
		if idx.Type() != object.INTEGER_OBJ {
			return nil, fmt.Errorf("array indices must be integers, not '%s'", idx.Type())
		}

		return evalArrayIndexExpression(l, idx)
	default:
		return nil, fmt.Errorf("index operator not supported: '%s'", l.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) (object.Object, error) {
	arr := array.(*object.Array).Elements
	idx := index.(*object.Integer).Value

	if idx < 0 || idx >= int64(len(arr)) {
		return nil, fmt.Errorf("array index out of range")
	}

	return arr[idx], nil
}
