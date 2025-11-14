package token

type TokenType int

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL TokenType = iota

	EOF

	// Identifiers + literals
	IDENT  // add, foobar, x, y, ...
	NUMBER // 1343456
	STRING

	// Operators
	ASSIGN
	PLUS
	MINUS
	BANG
	STAR
	SLASH
	LT
	GT

	EQ
	NOT_EQ

	// Delimiters
	COMMA
	SEMI
	COLON
	LPAREN
	RPAREN
	LBRACE
	RBRACE
	LBRACK
	RBRACK

	// Keywords
	FN
	LET
	TRUE
	FALSE
	IF
	ELSE
	RETURN
	NULL
	AND
	OR
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	IDENT:   "IDENT",
	NUMBER:  "NUMBER",
	STRING:  "STRING",
	ASSIGN:  "ASSIGN",
	PLUS:    "PLUS",
	MINUS:   "MINUS",
	BANG:    "BANG",
	STAR:    "STAR",
	SLASH:   "SLASH",
	LT:      "LT",
	GT:      "GT",
	EQ:      "EQ",
	NOT_EQ:  "NOT_EQ",
	COMMA:   "COMMA",
	SEMI:    "SEMI",
	COLON:   "COLON",
	LPAREN:  "LPAREN",
	RPAREN:  "RPAREN",
	LBRACE:  "LBRACE",
	RBRACE:  "RBRACE",
	LBRACK:  "LBRACK",
	RBRACK:  "RBRACK",
	FN:      "FN",
	LET:     "LET",
	TRUE:    "TRUE",
	FALSE:   "FALSE",
	IF:      "IF",
	ELSE:    "ELSE",
	RETURN:  "RETURN",
	NULL:    "NULL",
	OR:      "OR",
	AND:     "AND",
}

var keywords = map[string]TokenType{
	"fn":     FN,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"null":   NULL,
	"or":     OR,
	"and":    AND,
}

func LookupIdent(ident string) TokenType {
	if tt, ok := keywords[ident]; ok {
		return tt
	}
	return IDENT
}

func (tt TokenType) String() string {
	return tokens[tt]
}
