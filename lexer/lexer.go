package lexer

import "go-interpreter/token"

type Lexer struct {
	input        string
	position     int  // points to current char
	readPosition int  // after current char
	ch           byte // current char
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() *token.Token {
	tok := &token.Token{}

	l.skipWhitespace()

	switch s := string(l.ch); l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.NewToken(token.EQ, s+string(l.ch))
		} else {
			tok = token.NewToken(token.ASSIGN, s)
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.NewToken(token.NOT_EQ, s+string(l.ch))
		} else {
			tok = token.NewToken(token.BANG, s)
		}
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			tok = token.NewToken(token.LOGICAL_AND, s+string(l.ch))
		} else {
			tok = token.NewToken(token.BITWISE_AND, s)
		}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			tok = token.NewToken(token.LOGICAL_OR, s+string(l.ch))
		} else {
			tok = token.NewToken(token.BITWISE_OR, s)
		}
	case '^':
		tok = token.NewToken(token.BITWISE_XOR, s)
	case '~':
		tok = token.NewToken(token.BITWISE_NOT, s)
	case '+':
		tok = token.NewToken(token.PLUS, s)
	case '-':
		tok = token.NewToken(token.MINUS, s)
	case '*':
		tok = token.NewToken(token.ASTERISK, s)
	case '/':
		tok = token.NewToken(token.SLASH, s)
	case '%':
		tok = token.NewToken(token.MOD, s)
	case '<':
		switch l.peekChar() {
		case '<':
			l.readChar()
			tok = token.NewToken(token.SHL, s+string(l.ch))
		case '=':
			l.readChar()
			tok = token.NewToken(token.LE, s+string(l.ch))
		default:
			tok = token.NewToken(token.LT, s)
		}
	case '>':
		switch l.peekChar() {
		case '>':
			l.readChar()
			tok = token.NewToken(token.SHR, s+string(l.ch))
		case '=':
			l.readChar()
			tok = token.NewToken(token.GE, s+string(l.ch))
		default:
			tok = token.NewToken(token.GT, s)
		}
	case ',':
		tok = token.NewToken(token.COMMA, s)
	case ':':
		tok = token.NewToken(token.COLON, s)
	case ';':
		tok = token.NewToken(token.SEMICOLON, s)
	case '(':
		tok = token.NewToken(token.LPAREN, s)
	case ')':
		tok = token.NewToken(token.RPAREN, s)
	case '{':
		tok = token.NewToken(token.LBRACE, s)
	case '}':
		tok = token.NewToken(token.RBRACE, s)
	case '[':
		tok = token.NewToken(token.LBRACKET, s)
	case ']':
		tok = token.NewToken(token.RBRACKET, s)
	case '"':
		tok = token.NewToken(token.STRING, l.readString())
	case 0:
		tok = token.NewToken(token.EOF, "")
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok = token.NewToken(token.INT, l.readNumber())
			return tok
		} else {
			tok = token.NewToken(token.ILLEGAL, s)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	p := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[p:l.position]
}

func (l *Lexer) readNumber() string {
	p := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[p:l.position]
}

func (l *Lexer) readString() string {
	p := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[p:l.position]
}

func isLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
