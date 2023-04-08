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

	case *ast.ForLoopStatement:
		return evalForLoopStatement(node, env)

	case *ast.ExpressionList:
		return evalExpressionList(node, env)

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

	case *ast.ShortCircuitExpression:
		return evalShortCircuitExpression(node, env)

	case *ast.PrefixIncAndDec:
		return evalPrefixIncAndDec(node, env)

	case *ast.AssignmentConverter:
		return evalAssignmentConverter(node, env)

	case *ast.Assignment:
		return evalAssignment(node, env, false)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.FunctionLiteral:
		return evalFunctionLiteral(node, env)

	case *ast.CallExpression:
		return evalCallExpression(node, env)

	case *ast.ArrayLiteral:
		return evalArrayLiteral(node, env)

	case *ast.IndexExpression:
		return evalIndexExpression(node, env, false)

	default:
		return nil, errors.New("invalid syntax")
	}
}

func evalProgram(p *ast.Program, env *object.Environment) (res object.Object, err error) {
	for _, statement := range p.Stmts {
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

// TODO scope
func evalBlockStatement(bs *ast.BlockStatement, env *object.Environment) (res object.Object, err error) {
	for _, statement := range bs.Stmts {
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
	return Eval(es.Expr, env)
}

func evalLetStatement(ls *ast.LetStatement, env *object.Environment) (object.Object, error) {
	switch e := ls.Value.(type) {
	case *ast.Assignment:
		if _, err := evalAssignment(e, env, true); err != nil {
			return nil, err
		}

	case *ast.Identifier, *ast.ExpressionList:
		idents, err := checkAssignable(e, env, true)
		if err != nil {
			return nil, err
		}
		for _, id := range idents {
			id.Set(object.NULL)
		}

	default:
		return nil, errors.New("invalid syntax")
	}

	return nil, nil
}

func evalReturnStatement(rs *ast.ReturnStatement, env *object.Environment) (object.Object, error) {
	if rs.ReturnValue == nil {
		return object.NewReturnValue(object.NULL), nil
	}

	val, err := Eval(rs.ReturnValue, env)
	if err != nil {
		return nil, err
	}
	return object.NewReturnValue(val), nil
}

// TODO break continue
func evalForLoopStatement(fs *ast.ForLoopStatement, env *object.Environment) (object.Object, error) {
	env = object.NewEnclosedEnvironment(env)
	if fs.Init != nil {
		if _, err := Eval(fs.Init, env); err != nil {
			return nil, err
		}
	}

	for {
		if fs.Condition != nil {
			cond, err := Eval(fs.Condition, env)
			if err != nil {
				return nil, err
			}
			if !isTruthy(cond) {
				break
			}
		}
		if _, err := Eval(fs.Body, env); err != nil {
			return nil, err
		}
		if fs.Update != nil {
			if _, err := Eval(fs.Update, env); err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}

func evalIdentifier(ident *ast.Identifier, env *object.Environment) (object.Object, error) {
	name := ident.Value
	if val, _ := env.Get(name); val != nil {
		return val, nil
	}

	if builtin, ok := builtins[name]; ok {
		return builtin, nil
	}

	return nil, fmt.Errorf("name '%s' is not defined", name)
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

	return evalUnaryOperator(pe.Operator, val)
}

func evalUnaryOperator(operator string, obj object.Object) (object.Object, error) {
	fn, ok := unaryOperatorFns[operator]
	if !ok {
		return nil, fmt.Errorf("unknown operator: '%s'", operator)
	}

	return fn(obj)
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

	return evalBinaryOperator(ie.Operator, l, r)
}

func evalBinaryOperator(operator string, l, r object.Object) (object.Object, error) {
	fn, ok := binaryOperatorFns[operator]
	if !ok {
		return nil, fmt.Errorf("unknown operator: '%s'", operator)
	}

	return fn(l, r)
}

func evalShortCircuitExpression(sc *ast.ShortCircuitExpression, env *object.Environment) (object.Object, error) {
	l, err := Eval(sc.Left, env)
	if err != nil {
		return nil, err
	}

	switch sc.Operator {
	case "&&":
		if !isTruthy(l) {
			return l, nil
		}
	case "||":
		if isTruthy(l) {
			return l, nil
		}
	}

	return Eval(sc.Right, env)
}

func evalPrefixIncAndDec(p *ast.PrefixIncAndDec, env *object.Environment) (object.Object, error) {
	return evalAssignment(p.Expr, env, false)
}

func evalAssignmentConverter(ac *ast.AssignmentConverter, env *object.Environment) (object.Object, error) {
	if _, err := evalAssignment(ac.Expr, env, false); err != nil {
		return nil, err
	}
	return nil, nil
}

func evalAssignment(ae *ast.Assignment, env *object.Environment, isDeclaration bool) (object.Object, error) {
	as, err := checkAssignable(ae.Left, env, isDeclaration)
	if err != nil {
		return nil, err
	}

	var r object.Object
	if e, ok := ae.Right.(*ast.Assignment); ok {
		r, err = evalAssignment(e, env, isDeclaration)
	} else {
		r, err = Eval(ae.Right, env)
	}
	if err != nil {
		return nil, err
	}

	return assign(as, r)
}

func checkAssignable(e ast.Expression, env *object.Environment, isDeclaration bool) ([]object.Assignable, error) {
	var res []object.Assignable

	switch e := e.(type) {
	case *ast.Identifier:
		ident, err := newIdentifier(e.Value, env, isDeclaration)
		if err != nil {
			return nil, err
		}
		res = append(res, ident)
		return res, nil

	case *ast.ExpressionList:
		for _, v := range e.Exprs {
			tmp, err := checkAssignable(v, env, isDeclaration)
			if err != nil {
				return nil, err
			}
			res = append(res, tmp...)
		}
		return res, nil

	case *ast.IndexExpression:
		if !isDeclaration {
			ie, err := evalIndexExpression(e, env, true)
			if err != nil {
				return nil, err
			}
			res = append(res, ie.(object.Assignable))
			return res, nil
		}
	}

	return nil, errors.New("invalid syntax")
}

func newIdentifier(name string, env *object.Environment, isDeclaration bool) (*object.Identifier, error) {
	if isDeclaration {
		if env.IsExist(name) {
			return nil, fmt.Errorf("identifier '%s' has already been declared", name)
		}
		return object.NewIdentifier(name, env), nil
	}

	if val, identEnv := env.Get(name); val != nil {
		return object.NewIdentifier(name, identEnv), nil
	}
	return nil, fmt.Errorf("name '%s' is not defined", name)
}

func assign(as []object.Assignable, r object.Object) (object.Object, error) {
	if r == nil {
		return nil, errors.New("the right value of '=' cannot be empty")
	}

	if expList, ok := r.(*object.ExpressionList); ok {
		if len(as) != len(expList.Elements) {
			return nil, errors.New("the lengths of the left and right sides of '=' are not equal")
		}
		for i, o := range expList.Elements {
			as[i].Set(o)
		}
	} else {
		if len(as) != 1 {
			return nil, errors.New("the lengths of the left and right sides of '=' are not equal")
		}
		as[0].Set(r)
	}
	return r, nil
}

func evalExpressionList(el *ast.ExpressionList, env *object.Environment) (object.Object, error) {
	elements, err := evalExpressions(el.Exprs, env)
	if err != nil {
		return nil, err
	}
	return object.NewExpressionList(elements), nil
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

func evalFunctionLiteral(fl *ast.FunctionLiteral, env *object.Environment) (object.Object, error) {
	return object.NewFunction(fl.Params, fl.Body, env), nil
}

func evalCallExpression(ce *ast.CallExpression, env *object.Environment) (object.Object, error) {
	fn, err := Eval(ce.Func, env)
	if err != nil {
		return nil, err
	}

	args, err := evalExpressions(ce.Args, env)
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
		return unwrapReturnValue(res), nil

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

func evalIndexExpression(ie *ast.IndexExpression, env *object.Environment, isAssignment bool) (object.Object, error) {
	l, err := Eval(ie.Left, env)
	if err != nil {
		return nil, err
	}

	idx, err := Eval(ie.Indices, env)
	if err != nil {
		return nil, err
	}

	switch l.Type() {
	case object.ARRAY_OBJ:
		if idx.Type() != object.INTEGER_OBJ {
			return nil, fmt.Errorf("array indices must be integers, not '%s'", idx.Type())
		}

		return evalArrayIndexExpression(l, idx, isAssignment)
	default:
		return nil, fmt.Errorf("index operator not supported: '%s'", l.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object, isAssignment bool) (object.Object, error) {
	arr, idx := array.(*object.Array).Elements, index.(*object.Integer).Value
	if idx < 0 || idx >= int64(len(arr)) {
		return nil, fmt.Errorf("array index out of range")
	}
	if isAssignment {
		return object.NewArrayIndex(arr, idx), nil
	}
	return arr[idx], nil
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}
