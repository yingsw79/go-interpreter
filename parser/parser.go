package parser

import (
	"fmt"
	"go-interpreter/ast"
	"go-interpreter/lexer"
	"go-interpreter/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

type (
	prefixParseFn func() (ast.Expression, error)
	infixParseFn  func(ast.Expression) (ast.Expression, error)
)

type Parser struct {
	l                   *lexer.Lexer
	curToken, peekToken *token.Token
	errors              []error
	prefixParseFns      map[token.TokenType]prefixParseFn
	infixParseFns       map[token.TokenType]infixParseFn
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		prefixParseFns: map[token.TokenType]prefixParseFn{},
		infixParseFns:  map[token.TokenType]infixParseFn{},
	}

	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	for ; !p.curTokenIs(token.EOF); p.nextToken() {
		if stmt, err := p.parseStatement(); err != nil {
			p.errors = append(p.errors, err)
		} else {
			program.Statements = append(program.Statements, stmt)
		}
	}

	return program
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() (*ast.LetStatement, error) {
	stmt := &ast.LetStatement{Token: p.curToken}

	if err := p.expectPeek(token.IDENT); err != nil {
		return nil, err
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if err := p.expectPeek(token.ASSIGN); err != nil {
		return nil, err
	}

	// TODO: 跳过对表达式的处理，直到遇见分号
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt, nil
}

func (p *Parser) parseReturnStatement() (*ast.ReturnStatement, error) {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// TODO: 跳过对表达式的处理，直到遇见分号
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt, nil
}

func (p *Parser) parseExpressionStatement() (*ast.ExpressionStatement, error) {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	var err error
	stmt.Expression, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt, nil
}

func (p *Parser) parseIdentifier() (ast.Expression, error) {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}, nil
}

func (p *Parser) parseIntegerLiteral() (ast.Expression, error) {
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse %q as integer", p.curToken.Literal)
	}

	return &ast.IntegerLiteral{Token: p.curToken, Value: value}, nil
}

func (p *Parser) parsePrefixExpression() (ast.Expression, error) {
	exp := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	var err error
	exp.Right, err = p.parseExpression(PREFIX)

	return exp, err
}

func (p *Parser) parseInfixExpression(left ast.Expression) (ast.Expression, error) {
	exp := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()

	var err error
	exp.Right, err = p.parseExpression(precedence)

	return exp, err
}

func (p *Parser) parseExpression(precedence int) (exp ast.Expression, err error) {
	prefix, ok := p.prefixParseFns[p.curToken.Type]
	if !ok {
		return nil, fmt.Errorf("no prefix parse function for %s found", p.curToken.Type)
	}

	exp, err = prefix()
	if err != nil {
		return
	}

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix, ok := p.infixParseFns[p.peekToken.Type]
		if !ok {
			return nil, fmt.Errorf("no infix parse function for %s found", p.peekToken.Type)
		}

		p.nextToken()

		exp, err = infix(exp)
		if err != nil {
			return
		}
	}

	return
}

func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(tp token.TokenType) bool  { return p.curToken.Type == tp }
func (p *Parser) peekTokenIs(tp token.TokenType) bool { return p.peekToken.Type == tp }

func (p *Parser) expectPeek(tp token.TokenType) error {
	if p.peekTokenIs(tp) {
		p.nextToken()
		return nil
	}

	return fmt.Errorf("expected next token to be %s, got %s instead", tp, p.peekToken.Type)
}

func (p *Parser) registerPrefix(tp token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tp] = fn
}

func (p *Parser) registerInfix(tp token.TokenType, fn infixParseFn) {
	p.infixParseFns[tp] = fn
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}
