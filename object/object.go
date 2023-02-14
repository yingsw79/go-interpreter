package object

import (
	"go-interpreter/ast"
	"strconv"
	"strings"
)

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	FUNCTION_OBJ     = "FUNCTION"
	NULL_OBJ         = "NULL"
)

var (
	NULL  = &Null{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func NewInteger(v int64) *Integer   { return &Integer{Value: v} }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return strconv.FormatInt(i.Value, 10) }

type Boolean struct {
	Value bool
}

func NewBoolean(v bool) *Boolean {
	if v {
		return TRUE
	}
	return FALSE
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return strconv.FormatBool(b.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func NewReturnValue(v Object) *ReturnValue { return &ReturnValue{Value: v} }
func (rv *ReturnValue) Type() ObjectType   { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string    { return rv.Value.Inspect() }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func NewFunction(params []*ast.Identifier, body *ast.BlockStatement, env *Environment) *Function {
	return &Function{Parameters: params, Body: body, Env: env}
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var b strings.Builder

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	b.WriteString("fn")
	b.WriteString("(")
	b.WriteString(strings.Join(params, ", "))
	b.WriteString(") {\n")
	b.WriteString(f.Body.String())
	b.WriteString("\n}")

	return b.String()
}
