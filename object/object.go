package object

import (
	"fmt"
	"go-interpreter/ast"
	"strconv"
	"strings"
)

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	STRING_OBJ       = "STRING"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	FUNCTION_OBJ     = "FUNCTION"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
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

type String struct {
	Value string
}

func NewString(v string) *String   { return &String{Value: v} }
func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return fmt.Sprintf("%q", s.Value) }

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

type BuiltinFn func(...Object) (Object, error)

type Builtin struct {
	Fn BuiltinFn
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

type Array struct {
	Elements []Object
}

func NewArray(es []Object) *Array   { return &Array{Elements: es} }
func (arr *Array) Type() ObjectType { return ARRAY_OBJ }
func (arr *Array) Inspect() string {
	var b strings.Builder

	elements := []string{}
	for _, e := range arr.Elements {
		elements = append(elements, e.Inspect())
	}

	b.WriteString("[")
	b.WriteString(strings.Join(elements, ", "))
	b.WriteString("]")

	return b.String()
}
