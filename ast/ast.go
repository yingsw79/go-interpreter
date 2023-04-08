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
	Stmts []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Stmts) > 0 {
		return p.Stmts[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var b strings.Builder
	for _, s := range p.Stmts {
		b.WriteString(s.String())
	}
	return b.String()
}

// Statements
type LetStatement struct {
	Token *token.Token // the token.LET token
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var b strings.Builder
	b.WriteString(ls.TokenLiteral() + " ")
	b.WriteString(ls.Value.String())
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

type ForLoopStatement struct {
	Token             *token.Token // The 'for' token
	Init              Statement
	Condition, Update Expression
	Body              *BlockStatement
}

func (f *ForLoopStatement) statementNode()       {}
func (f *ForLoopStatement) TokenLiteral() string { return f.Token.Literal }
func (f *ForLoopStatement) String() string {
	var b strings.Builder
	b.WriteString(f.TokenLiteral())
	b.WriteString(" (")
	if f.Init != nil {
		b.WriteString(f.Init.String())
	} else {
		b.WriteString(";")
	}
	if f.Condition != nil {
		b.WriteString(f.Condition.String())
	}
	b.WriteString(";")
	if f.Update != nil {
		b.WriteString(f.Update.String())
	}
	b.WriteString(") ")
	b.WriteString(f.Body.String())
	return b.String()
}

type ExpressionStatement struct {
	Token *token.Token // the first token of the expression
	Expr  Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() (res string) {
	if es.Expr != nil {
		res = es.Expr.String()
	}
	return res + ";"
}

type BlockStatement struct {
	Token *token.Token // the '{' token
	Stmts []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var b strings.Builder
	stmts := []string{}
	for _, stmt := range bs.Stmts {
		stmts = append(stmts, stmt.String())
	}
	b.WriteString("{ ")
	b.WriteString(strings.Join(stmts, " "))
	b.WriteString(" }")
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

type StringLiteral struct {
	Token *token.Token
	Value string
}

func NewStringLiteral(tok *token.Token, v string) *StringLiteral {
	return &StringLiteral{Token: tok, Value: v}
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

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
	b.WriteString("if ")
	b.WriteString("(" + ie.Condition.String() + ")")
	b.WriteString(" ")
	b.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		b.WriteString(" else ")
		b.WriteString(ie.Alternative.String())
	}
	return b.String()
}

type FunctionLiteral struct {
	Token  *token.Token // The 'fn' token
	Params []*Identifier
	Body   *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var b strings.Builder
	params := []string{}
	for _, p := range fl.Params {
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
	Token *token.Token // The '(' token
	Func  Expression   // Identifier or FunctionLiteral
	Args  []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var b strings.Builder
	args := []string{}
	for _, a := range ce.Args {
		args = append(args, a.String())
	}
	b.WriteString(ce.Func.String())
	b.WriteString("(")
	b.WriteString(strings.Join(args, ", "))
	b.WriteString(")")
	return b.String()
}

type ArrayLiteral struct {
	Token    *token.Token // The '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var b strings.Builder
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	b.WriteString("[")
	b.WriteString(strings.Join(elements, ", "))
	b.WriteString("]")
	return b.String()
}

type IndexExpression struct {
	Token         *token.Token // The '[' token
	Left, Indices Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var b strings.Builder
	b.WriteString("(")
	b.WriteString(ie.Left.String())
	b.WriteString("[")
	b.WriteString(ie.Indices.String())
	b.WriteString("])")
	return b.String()
}

type Assignment struct {
	Token       *token.Token // The '=' token
	Left, Right Expression
}

func (as *Assignment) expressionNode()      {}
func (as *Assignment) TokenLiteral() string { return as.Token.Literal }
func (as *Assignment) String() string {
	var b strings.Builder
	b.WriteString("(")
	b.WriteString(as.Left.String())
	b.WriteString(" = ")
	b.WriteString(as.Right.String())
	b.WriteString(")")
	return b.String()
}

type PrefixIncAndDec struct {
	Token *token.Token // '++' '--' '+='
	Expr  *Assignment
}

func (p *PrefixIncAndDec) expressionNode()      {}
func (p *PrefixIncAndDec) TokenLiteral() string { return p.Token.Literal }
func (p *PrefixIncAndDec) String() string       { return p.Expr.String() }

type AssignmentConverter struct {
	Token *token.Token // '++' '--' '+='
	Expr  *Assignment
}

func (ac *AssignmentConverter) expressionNode()      {}
func (ac *AssignmentConverter) TokenLiteral() string { return ac.Token.Literal }
func (ac *AssignmentConverter) String() string       { return ac.Expr.String() }

type ShortCircuitExpression struct {
	Token       *token.Token
	Left, Right Expression
	Operator    string
}

func (sc *ShortCircuitExpression) expressionNode()      {}
func (sc *ShortCircuitExpression) TokenLiteral() string { return sc.Token.Literal }
func (sc *ShortCircuitExpression) String() string {
	var b strings.Builder
	b.WriteString("(")
	b.WriteString(sc.Left.String())
	b.WriteString(" " + sc.Operator + " ")
	b.WriteString(sc.Right.String())
	b.WriteString(")")
	return b.String()
}

type ExpressionList struct {
	Token *token.Token
	Exprs []Expression
}

func (el *ExpressionList) expressionNode()      {}
func (el *ExpressionList) TokenLiteral() string { return el.Token.Literal }
func (el *ExpressionList) String() string {
	var b strings.Builder
	exprs := []string{}
	for _, el := range el.Exprs {
		exprs = append(exprs, el.String())
	}
	b.WriteString("(")
	b.WriteString(strings.Join(exprs, ", "))
	b.WriteString(")")
	return b.String()
}
