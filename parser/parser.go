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
	ASSIGN      // =
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

var precedences = map[token.TokenType]int{
	token.ASSIGN:   ASSIGN,
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
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
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)

	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.ASSIGN, p.parseInfixExpression)

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
	var err error

	if err = p.expectPeek(token.IDENT); err != nil {
		return nil, err
	}

	stmt.Name = p.newIdentifier()

	if err = p.expectPeek(token.ASSIGN); err != nil {
		return nil, err
	}

	p.nextToken()

	stmt.Value, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt, nil
}

func (p *Parser) parseReturnStatement() (*ast.ReturnStatement, error) {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	var err error
	stmt.ReturnValue, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if p.peekTokenIs(token.SEMICOLON) {
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

func (p *Parser) parseBlockStatement() (*ast.BlockStatement, error) {
	block := &ast.BlockStatement{Token: p.curToken}

	p.nextToken()

	for ; !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF); p.nextToken() {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		block.Statements = append(block.Statements, stmt)
	}

	return block, nil
}

func (p *Parser) parseExpression(precedence int) (ast.Expression, error) {
	prefix, ok := p.prefixParseFns[p.curToken.Type]
	if !ok {
		return nil, fmt.Errorf("no prefix parse function for '%s' found", p.curToken.Type)
	}

	exp, err := prefix()
	if err != nil {
		return nil, err
	}

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix, ok := p.infixParseFns[p.peekToken.Type]
		if !ok {
			return nil, fmt.Errorf("no infix parse function for '%s' found", p.peekToken.Type)
		}

		p.nextToken()

		exp, err = infix(exp)
		if err != nil {
			return nil, err
		}
	}

	return exp, nil
}

func (p *Parser) parsePrefixExpression() (ast.Expression, error) {
	exp := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}

	p.nextToken()

	var err error
	exp.Right, err = p.parseExpression(PREFIX)
	if err != nil {
		return nil, err
	}

	return exp, nil
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
	if err != nil {
		return nil, err
	}

	return exp, nil
}

func (p *Parser) parseGroupedExpression() (ast.Expression, error) {
	p.nextToken()

	exp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if err = p.expectPeek(token.RPAREN); err != nil {
		return nil, err
	}

	return exp, nil
}

func (p *Parser) parseIfExpression() (ast.Expression, error) {
	exp := &ast.IfExpression{Token: p.curToken}
	var err error

	if err = p.expectPeek(token.LPAREN); err != nil {
		return nil, err
	}

	exp.Condition, err = p.parseGroupedExpression()
	if err != nil {
		return nil, err
	}

	if err = p.expectPeek(token.LBRACE); err != nil {
		return nil, err
	}

	exp.Consequence, err = p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if err = p.expectPeek(token.LBRACE); err != nil {
			return nil, err
		}

		exp.Alternative, err = p.parseBlockStatement()
		if err != nil {
			return nil, err
		}
	}

	return exp, nil
}

func (p *Parser) parseFunctionLiteral() (ast.Expression, error) {
	lit := &ast.FunctionLiteral{Token: p.curToken}
	var err error

	if err = p.expectPeek(token.LPAREN); err != nil {
		return nil, err
	}

	lit.Parameters, err = p.parseFunctionParameters()
	if err != nil {
		return nil, err
	}

	if err = p.expectPeek(token.LBRACE); err != nil {
		return nil, err
	}

	lit.Body, err = p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	return lit, nil
}

func (p *Parser) parseFunctionParameters() ([]*ast.Identifier, error) {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers, nil
	}

	p.nextToken()
	identifiers = append(identifiers, p.newIdentifier())

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		identifiers = append(identifiers, p.newIdentifier())
	}

	if err := p.expectPeek(token.RPAREN); err != nil {
		return nil, err
	}

	return identifiers, nil
}

func (p *Parser) parseCallExpression(function ast.Expression) (ast.Expression, error) {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}

	var err error
	exp.Arguments, err = p.parseExpressionList(token.COMMA, token.RPAREN)
	if err != nil {
		return nil, err
	}

	return exp, nil
}

func (p *Parser) parseExpressionList(sep, end token.TokenType) ([]ast.Expression, error) {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list, nil
	}

	p.nextToken()

	exp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	list = append(list, exp)

	for !p.peekTokenIs(end) {
		if err = p.expectPeek(sep); err != nil {
			return nil, err
		}

		p.nextToken()

		exp, err = p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}

		list = append(list, exp)
	}

	p.nextToken()

	return list, nil
}

func (p *Parser) parseArrayLiteral() (ast.Expression, error) {
	array := &ast.ArrayLiteral{Token: p.curToken}

	var err error
	array.Elements, err = p.parseExpressionList(token.COMMA, token.RBRACKET)
	if err != nil {
		return nil, err
	}

	return array, nil
}

// TODO: SliceExpression
func (p *Parser) parseIndexExpression(left ast.Expression) (ast.Expression, error) {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()

	var err error
	exp.Indices, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if err = p.expectPeek(token.RBRACKET); err != nil {
		return nil, err
	}

	return exp, nil
}

// TODO a, b = c, d = 1, 2
func (p *Parser) parseIdentifier() (ast.Expression, error) {
	return p.newIdentifier(), nil
}

func (p *Parser) newIdentifier() *ast.Identifier {
	return ast.NewIdentifier(p.curToken, p.curToken.Literal, p.peekTokenIs(token.ASSIGN))
}

func (p *Parser) parseStringLiteral() (ast.Expression, error) {
	return ast.NewStringLiteral(p.curToken, p.curToken.Literal), nil
}

func (p *Parser) parseIntegerLiteral() (ast.Expression, error) {
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse '%s' as integer", p.curToken.Literal)
	}

	return ast.NewIntegerLiteral(p.curToken, value), nil
}

func (p *Parser) parseBoolean() (ast.Expression, error) {
	return ast.NewBoolean(p.curToken, p.curTokenIs(token.TRUE)), nil
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

	return fmt.Errorf("expected next token to be '%s', got '%s' instead", tp, p.peekToken.Type)
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) registerPrefix(tp token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tp] = fn
}

func (p *Parser) registerInfix(tp token.TokenType, fn infixParseFn) {
	p.infixParseFns[tp] = fn
}
