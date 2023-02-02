package parser

import (
	"fmt"
	"go-interpreter/ast"
	"go-interpreter/lexer"
	"go-interpreter/token"
)

type Parser struct {
	l                   *lexer.Lexer
	curToken, peekToken *token.Token
	errors              []error
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	for ; !p.curTokenIs(token.EOF); p.nextToken() {
		if stmt, err := p.parseStatement(); err != nil {
			p.errors = append(p.errors, err)
		} else if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
	}

	return program
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil, nil
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

func (p *Parser) curTokenIs(tp token.TokenType) bool  { return p.curToken.Type == tp }
func (p *Parser) peekTokenIs(tp token.TokenType) bool { return p.peekToken.Type == tp }

func (p *Parser) expectPeek(tp token.TokenType) error {
	if p.peekTokenIs(tp) {
		p.nextToken()
		return nil
	}

	return fmt.Errorf("expected next token to be %s, got %s instead", tp, p.peekToken.Type)
}
