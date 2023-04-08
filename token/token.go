package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT"
	STRING = "STRING"
	INT    = "INT"

	// Operators
	ASSIGN             = "="
	BANG               = "!"
	INC                = "++"
	DEC                = "--"
	PLUS               = "+"
	PLUS_ASSIGN        = "+="
	MINUS              = "-"
	MINUS_ASSIGN       = "-="
	ASTERISK           = "*"
	ASTERISK_ASSIGN    = "*="
	SLASH              = "/"
	SLASH_ASSIGN       = "/="
	MOD                = "%"
	MOD_ASSIGN         = "%="
	LT                 = "<"
	GT                 = ">"
	LE                 = "<="
	GE                 = ">="
	EQ                 = "=="
	NOT_EQ             = "!="
	LOGICAL_AND        = "&&"
	LOGICAL_OR         = "||"
	BITWISE_AND        = "&"
	BITWISE_AND_ASSIGN = "&="
	BITWISE_OR         = "|"
	BITWISE_OR_ASSIGN  = "|="
	BITWISE_XOR        = "^"
	BITWISE_XOR_ASSIGN = "^="
	BITWISE_NOT        = "~"
	SHL                = "<<"
	SHL_ASSIGN         = "<<="
	SHR                = ">>"
	SHR_ASSIGN         = ">>="

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
	FORLOOP  = "FOR"
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
	"for":    FORLOOP,
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
