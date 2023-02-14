package ast

import (
	"go-interpreter/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var b strings.Builder

	for _, s := range p.Statements {
		b.WriteString(s.String())
	}

	return b.String()
}

// Statements
type LetStatement struct {
	Token *token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var b strings.Builder

	b.WriteString(ls.TokenLiteral() + " ")
	b.WriteString(ls.Name.String())
	b.WriteString(" = ")

	if ls.Value != nil {
		b.WriteString(ls.Value.String())
	}

	b.WriteString(";")

	return b.String()
}

type ReturnStatement struct {
	Token       *token.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var b strings.Builder

	b.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		b.WriteString(rs.ReturnValue.String())
	}

	b.WriteString(";")

	return b.String()
}

type ExpressionStatement struct {
	Token      *token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type BlockStatement struct {
	Token      *token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var b strings.Builder

	stmts := []string{}
	for _, stmt := range bs.Statements {
		stmts = append(stmts, stmt.String())
	}

	b.WriteString(strings.Join(stmts, "\n"))

	return b.String()
}

// Expressions
type Identifier struct {
	Token *token.Token // the token.IDENT token
	Value string
}

func NewIdentifier(tok *token.Token, v string) *Identifier {
	return &Identifier{Token: tok, Value: v}
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type IntegerLiteral struct {
	Token *token.Token
	Value int64
}

func NewIntegerLiteral(tok *token.Token, v int64) *IntegerLiteral {
	return &IntegerLiteral{Token: tok, Value: v}
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type Boolean struct {
	Token *token.Token
	Value bool
}

func NewBoolean(tok *token.Token, v bool) *Boolean {
	return &Boolean{Token: tok, Value: v}
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type PrefixExpression struct {
	Token    *token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var b strings.Builder

	b.WriteString("(")
	b.WriteString(pe.Operator)
	b.WriteString(pe.Right.String())
	b.WriteString(")")

	return b.String()
}

type InfixExpression struct {
	Token       *token.Token // The operator token, e.g. +
	Left, Right Expression
	Operator    string
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var b strings.Builder

	b.WriteString("(")
	b.WriteString(ie.Left.String())
	b.WriteString(" " + ie.Operator + " ")
	b.WriteString(ie.Right.String())
	b.WriteString(")")

	return b.String()
}

type IfExpression struct {
	Token                    *token.Token // The 'if' token
	Condition                Expression
	Consequence, Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var b strings.Builder

	b.WriteString("if")
	b.WriteString(ie.Condition.String())
	b.WriteString(" ")
	b.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		b.WriteString("else ")
		b.WriteString(ie.Alternative.String())
	}

	return b.String()
}

type FunctionLiteral struct {
	Token      *token.Token // The 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var b strings.Builder

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	b.WriteString(fl.TokenLiteral())
	b.WriteString("(")
	b.WriteString(strings.Join(params, ", "))
	b.WriteString(") ")
	b.WriteString(fl.Body.String())

	return b.String()
}

type CallExpression struct {
	Token     *token.Token // The '(' token
	Function  Expression   // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var b strings.Builder

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	b.WriteString(ce.Function.String())
	b.WriteString("(")
	b.WriteString(strings.Join(args, ", "))
	b.WriteString(")")

	return b.String()
}
