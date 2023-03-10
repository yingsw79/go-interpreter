package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT"
	STRING = "STRING"
	INT    = "INT"

	// Operators
	ASSIGN      = "="
	BANG        = "!"
	PLUS        = "+"
	MINUS       = "-"
	ASTERISK    = "*"
	SLASH       = "/"
	MOD         = "%"
	LT          = "<"
	GT          = ">"
	LE          = "<="
	GE          = ">="
	EQ          = "=="
	NOT_EQ      = "!="
	LOGICAL_AND = "&&"
	LOGICAL_OR  = "||"
	BITWISE_AND = "&"
	BITWISE_OR  = "|"
	BITWISE_XOR = "^"
	BITWISE_NOT = "~"
	SHL         = "<<"
	SHR         = ">>"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

type TokenType string

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

type Token struct {
	Type    TokenType
	Literal string
}

func NewToken(tp TokenType, s string) *Token {
	return &Token{
		Type:    tp,
		Literal: s,
	}
}
