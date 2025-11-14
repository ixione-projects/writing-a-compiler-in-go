package lexer

import (
	"slices"

	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/token"
)

type Lexer struct {
	input   string
	start   int
	current int

	ch  byte
	eof bool

	tokens []token.Token
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: input, tokens: []token.Token{}}
}

func (l *Lexer) BufferLength() int {
	return len(l.tokens)
}

func (l *Lexer) Token(index int) token.Token {
	var ok bool = false
	if l.BufferLength() < index+1 {
		ok = l.ensure(index + 1 - len(l.tokens))
	}

	if ok {
		return l.tokens[index]
	}
	return l.tokens[len(l.tokens)-1]
}

func (l *Lexer) ensure(n int) bool {
	for range n {
		tok := l.NextToken()
		if tok.Type == token.EOF {
			return false
		}
	}
	return true
}

func (l *Lexer) Tokens() []token.Token {
	for l.tokens[len(l.tokens)-1].Type != token.EOF {
		l.NextToken()
	}
	return l.tokens
}

func (l *Lexer) NextToken() token.Token {
	if l.eof {
		return l.tokens[len(l.tokens)-1]
	}

	for !l.isEOF() {
		l.start = l.current

		l.ch = l.peek0()
		switch l.ch {
		case ' ', '\t', '\r', '\n':
			l.skip(' ', '\t', '\r', '\n')
		case '=':
			l.next()
			if l.match('=') {
				return l.emit(token.EQ)
			}
			return l.emit(token.ASSIGN)
		case '+':
			l.next()
			return l.emit(token.PLUS)
		case '-':
			l.next()
			return l.emit(token.MINUS)
		case '!':
			l.next()
			if l.match('=') {
				return l.emit(token.NOT_EQ)
			}
			return l.emit(token.BANG)
		case '/':
			l.next()
			return l.emit(token.SLASH)
		case '*':
			l.next()
			return l.emit(token.STAR)
		case '<':
			l.next()
			return l.emit(token.LT)
		case '>':
			l.next()
			return l.emit(token.GT)
		case ';':
			l.next()
			return l.emit(token.SEMI)
		case ':':
			l.next()
			return l.emit(token.COLON)
		case ',':
			l.next()
			return l.emit(token.COMMA)
		case '(':
			l.next()
			return l.emit(token.LPAREN)
		case ')':
			l.next()
			return l.emit(token.RPAREN)
		case '{':
			l.next()
			return l.emit(token.LBRACE)
		case '}':
			l.next()
			return l.emit(token.RBRACE)
		case '[':
			l.next()
			return l.emit(token.LBRACK)
		case ']':
			l.next()
			return l.emit(token.RBRACK)
		case '"':
			return l.string()
		default:
			if isAlpha(l.ch) {
				return l.ident()
			} else if isNumber(l.ch) {
				return l.number()
			} else {
				l.next()
				return l.emit(token.ILLEGAL)
			}
		}
	}

	l.start = l.current
	l.eof = true
	return l.emit(token.EOF)
}

func (l *Lexer) ident() token.Token {
	for isAlphaNumeric(l.ch) {
		l.next()
	}
	ident := l.input[l.start:l.current]
	tok := token.Token{
		Type:    token.LookupIdent(ident),
		Literal: ident,
	}
	l.tokens = append(l.tokens, tok)
	return tok
}

func (l *Lexer) number() token.Token {
	for isNumber(l.ch) {
		l.next()
	}

	if l.ch == '.' {
		if !isNumber(l.peek1()) {
			return token.Token{
				Type:    token.ILLEGAL,
				Literal: l.input[l.start:l.current],
			}
		}

		l.next()
		for isNumber(l.ch) {
			l.next()
		}
	}

	tok := token.Token{
		Type:    token.NUMBER,
		Literal: l.input[l.start:l.current],
	}
	l.tokens = append(l.tokens, tok)
	return tok
}

func (l *Lexer) string() token.Token {
	l.next()
	for l.ch != '"' && l.ch != 0 {
		l.next()
	}

	if l.ch != '"' {
		return token.Token{
			Type:    token.ILLEGAL,
			Literal: l.input[l.start:l.current],
		}
	}

	l.next()

	tok := token.Token{
		Type:    token.STRING,
		Literal: l.input[l.start:l.current],
	}
	l.tokens = append(l.tokens, tok)
	return tok
}

func (l *Lexer) emit(ttype token.TokenType) token.Token {
	tok := token.Token{
		Type:    ttype,
		Literal: l.input[l.start:l.current],
	}
	l.tokens = append(l.tokens, tok)
	return tok
}

func (l *Lexer) skip(chars ...byte) {
	for slices.Contains(chars, l.ch) {
		l.next()
	}
}

func (l *Lexer) match(char byte) bool {
	if l.ch == char {
		l.next()
		return true
	}
	return false
}

func (l *Lexer) peek0() byte {
	if l.isEOF() {
		return 0
	}
	return l.input[l.current]
}

func (l *Lexer) peek1() byte {
	if l.current+1 >= len(l.input) {
		return 0
	}
	return l.input[l.current+1]
}

func (l *Lexer) next() {
	if l.isEOF() {
		return
	}

	l.current += 1
	if l.isEOF() {
		l.ch = 0
	} else {
		l.ch = l.input[l.current]
	}
}

func (l *Lexer) isEOF() bool {
	return l.current >= len(l.input)
}

func isAlphaNumeric(ch byte) bool {
	return isAlpha(ch) || isNumber(ch)
}

func isAlpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isNumber(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
