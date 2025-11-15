package parser

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/ast"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/lexer"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/token"
)

type Parser struct {
	l       *lexer.Lexer
	trace   bool
	current int

	tok    token.Token
	errors []string

	rules map[token.TokenType]parseRule
}

type precedence int

const (
	_ precedence = iota
	NONE
	ASSIGNMENT
	OR
	AND
	EQUALITY
	COMPARISON
	TERM
	FACTOR
	UNARY
	CALL
	SUBSCRIPT
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type parseRule struct {
	PrefixParseFn prefixParseFn
	InfixParseFn  infixParseFn
	Precedence    precedence
}

func NewParser(input string, trace bool) *Parser {
	l := lexer.NewLexer(input)
	p := &Parser{
		l:      l,
		errors: []string{},
		trace:  trace,
	}
	p.tok = p.peek0()

	p.rules = map[token.TokenType]parseRule{
		token.ILLEGAL: {p.reportIllegalToken, nil, NONE},
		token.EOF:     {nil, nil, NONE},
		token.IDENT:   {p.parseIdentifier, nil, NONE},
		token.NUMBER:  {p.parseNumberLiteral, nil, NONE},
		token.STRING:  {p.parseStringLiteral, nil, NONE},
		token.ASSIGN:  {nil, p.parseAssignmentExpression, ASSIGNMENT},
		token.PLUS:    {nil, p.parseBinaryExpression, TERM},
		token.MINUS:   {p.parseUnaryExpression, p.parseBinaryExpression, TERM},
		token.BANG:    {p.parseUnaryExpression, nil, NONE},
		token.STAR:    {nil, p.parseBinaryExpression, FACTOR},
		token.SLASH:   {nil, p.parseBinaryExpression, FACTOR},
		token.LT:      {nil, p.parseBinaryExpression, COMPARISON},
		token.GT:      {nil, p.parseBinaryExpression, COMPARISON},
		token.EQ:      {nil, p.parseBinaryExpression, EQUALITY},
		token.NOT_EQ:  {nil, p.parseBinaryExpression, EQUALITY},
		token.COMMA:   {nil, nil, NONE},
		token.SEMI:    {nil, nil, NONE},
		token.COLON:   {nil, nil, NONE},
		token.LPAREN:  {p.parseGroupingExpression, p.parseCallExpression, CALL},
		token.RPAREN:  {nil, nil, NONE},
		token.LBRACE:  {p.parseHashLiteral, nil, NONE},
		token.RBRACE:  {nil, nil, NONE},
		token.LBRACK:  {p.parseArrayLiteral, p.parseSubscriptExpression, SUBSCRIPT},
		token.RBRACK:  {nil, nil, NONE},
		token.FN:      {p.parseFunctionLiteral, nil, NONE},
		token.LET:     {nil, nil, NONE},
		token.TRUE:    {p.parseBooleanLiteral, nil, NONE},
		token.FALSE:   {p.parseBooleanLiteral, nil, NONE},
		token.IF:      {p.parseIfExpression, nil, NONE},
		token.ELSE:    {nil, nil, NONE},
		token.RETURN:  {nil, nil, NONE},
		token.NULL:    {p.parseNullLiteral, nil, NONE},
		token.OR:      {nil, p.parseLogicalExpression, OR},
		token.AND:     {nil, p.parseLogicalExpression, AND},
	}

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	if p.trace {
		defer un(trace("ParseProgram"))
	}

	statements := []ast.Statement{}
	for ; p.tok.Type != token.EOF; p.next() {
		p.tok = p.peek0()

		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		}
	}
	return &ast.Program{Statements: statements}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) parseStatement() ast.Statement {
	if p.trace {
		defer un(trace("ParseStatement"))
	}

	switch p.tok.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.SEMI:
		p.skip(token.SEMI)
		return nil
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetDeclaration {
	if p.trace {
		defer un(trace("ParseLetStatement"))
	}

	stmt := &ast.LetDeclaration{Token: p.tok}
	if !p.expect(token.IDENT) {
		return nil
	}

	stmt.Name = p.parseIdentifier().(*ast.Identifier)
	if !p.expect(token.ASSIGN) {
		return nil
	}

	p.next()

	stmt.Value = p.parseExpression(ASSIGNMENT - 1) // right associativity

	if p.peek1().Type == token.SEMI {
		p.next()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	if p.trace {
		defer un(trace("ParseReturnStatement"))
	}

	stmt := &ast.ReturnStatement{Token: p.tok}

	p.next()

	stmt.ReturnValue = p.parseExpression(ASSIGNMENT - 1) // right associativity

	if p.peek1().Type == token.SEMI {
		p.next()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	if p.trace {
		defer un(trace("ParseExpressionStatement"))
	}

	stmt := &ast.ExpressionStatement{Token: p.tok}
	stmt.Expression = p.parseExpression(ASSIGNMENT - 1) // right associativity

	if p.peek1().Type == token.SEMI {
		p.next()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	if p.trace {
		defer un(trace("ParseBlockStatement"))
	}

	block := &ast.BlockStatement{Token: p.tok}

	p.next()

	statements := []ast.Statement{}
	for ; p.tok.Type != token.RBRACE && p.tok.Type != token.EOF; p.next() {
		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		}
	}
	block.Statements = statements

	if p.tok.Type != token.RBRACE {
		p.error("expected token to be <RBRACE> but was <%s>", p.tok.Type)
		return nil
	}

	return block
}

func (p *Parser) parseExpression(rightPrecedence precedence) ast.Expression {
	if p.trace {
		defer un(trace("ParseExpression"))
	}

	prefix := p.getRule(p.tok.Type).PrefixParseFn
	if prefix == nil {
		p.error("no prefix parse function defined for %s", p.tok.Type)
		return nil
	}
	left := prefix()

	for rightPrecedence < p.getRule(p.peek1().Type).Precedence {
		p.next()

		infix := p.getRule(p.tok.Type).InfixParseFn
		if infix == nil {
			return left
		}
		left = infix(left)
	}

	return left
}

func (p *Parser) parseGroupingExpression() ast.Expression {
	p.next()

	expr := p.parseExpression(ASSIGNMENT - 1) // right associativity
	if !p.expect(token.RPAREN) {
		return nil
	}

	return expr
}

func (p *Parser) parseUnaryExpression() ast.Expression {
	if p.trace {
		defer un(trace("ParseUnaryExpression"))
	}

	expr := &ast.UnaryExpression{
		Token:    p.tok,
		Operator: p.tok.Literal,
	}

	p.next()

	expr.Right = p.parseExpression(UNARY)

	return expr
}

func (p *Parser) parseBinaryExpression(left ast.Expression) ast.Expression {
	if p.trace {
		defer un(trace("ParseBinaryExpression"))
	}

	expr := &ast.BinaryExpression{
		Token:    p.tok,
		Left:     left,
		Operator: p.tok.Literal,
	}
	rule := p.getRule(p.tok.Type)

	p.next()

	expr.Right = p.parseExpression(rule.Precedence)

	return expr
}

func (p *Parser) parseLogicalExpression(left ast.Expression) ast.Expression {
	if p.trace {
		defer un(trace("ParseLogicalExpression"))
	}

	expr := &ast.BinaryExpression{
		Token:    p.tok,
		Left:     left,
		Operator: p.tok.Literal,
	}
	rule := p.getRule(p.tok.Type)

	p.next()

	expr.Right = p.parseExpression(rule.Precedence)

	return expr
}

func (p *Parser) parseIfExpression() ast.Expression {
	if p.trace {
		defer un(trace("ParseIfExpression"))
	}

	expr := &ast.ConditionalExpression{Token: p.tok}
	if !p.expect(token.LPAREN) {
		return nil
	}

	p.next()

	expr.Condition = p.parseExpression(ASSIGNMENT - 1) // right associativity

	if !p.expect(token.RPAREN) {
		return nil
	}

	if !p.expect(token.LBRACE) {
		return nil
	}

	expr.Consequence = p.parseBlockStatement()
	if p.peek1().Type == token.ELSE {
		p.next()

		if !p.expect(token.LBRACE) {
			return nil
		}
		expr.Alternative = p.parseBlockStatement()
	}

	return expr
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	if p.trace {
		defer un(trace("ParseFunctionLiteral"))
	}

	expr := &ast.FunctionLiteral{Token: p.tok}
	if !p.expect(token.LPAREN) {
		return nil
	}

	expr.Parameters = p.parseIdentifierList(token.RPAREN)

	if !p.expect(token.RPAREN) {
		return nil
	}

	if !p.expect(token.LBRACE) {
		return nil
	}

	expr.Body = p.parseBlockStatement()

	return expr
}

func (p *Parser) parseCallExpression(left ast.Expression) ast.Expression {
	if p.trace {
		defer un(trace("ParseCallExpression"))
	}

	expr := &ast.CallExpression{
		Token:  p.tok,
		Callee: left,
	}
	expr.Arguments = p.parseExpressionList(token.RPAREN)

	if !p.expect(token.RPAREN) {
		return nil
	}

	return expr
}

func (p *Parser) parseSubscriptExpression(left ast.Expression) ast.Expression {
	if p.trace {
		defer un(trace("ParseSubscriptExpression"))
	}

	expr := &ast.SubscriptExpression{
		Token: p.tok,
		Base:  left,
	}

	p.next()

	expr.Subscript = p.parseExpression(ASSIGNMENT - 1) // right associativity

	if !p.expect(token.RBRACK) {
		return nil
	}

	return expr
}

var lvalues = []ast.NodeType{
	ast.IDENTIFIER,
	ast.SUBSCRIPT_EXPRESSION,
}

func (p *Parser) parseAssignmentExpression(left ast.Expression) ast.Expression {
	if p.trace {
		defer un(trace("ParseAssignmentExpression"))
	}

	if !slices.Contains(lvalues, left.Type()) {
		p.error("unexpected lvalue type <%s>", left.Type())
		return nil
	}

	expr := &ast.AssignmentExpression{
		Token:  p.tok,
		LValue: left,
	}

	p.next()

	expr.RValue = p.parseExpression(ASSIGNMENT - 1) // right associativity

	return expr
}

func (p *Parser) parseIdentifier() ast.Expression {
	if p.trace {
		defer un(trace("ParseIdentifier"))
	}

	return &ast.Identifier{Token: p.tok, Value: p.tok.Literal}
}

func (p *Parser) parseNumberLiteral() ast.Expression {
	if p.trace {
		defer un(trace("ParseNumberLiteral"))
	}

	value, err := strconv.ParseFloat(p.tok.Literal, 64)
	if err != nil {
		p.error("cannot parse float %q", p.tok.Literal)
		return nil
	}
	return &ast.NumberLiteral{Token: p.tok, Value: value}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	if p.trace {
		defer un(trace("ParseStringLiteral"))
	}

	return &ast.StringLiteral{
		Token: p.tok,
		Value: p.tok.Literal[1 : len(p.tok.Literal)-1],
	}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	if p.trace {
		defer un(trace("ParseBooleanLiteral"))
	}

	return &ast.BooleanLiteral{Token: p.tok, Value: p.tok.Type == token.TRUE}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	if p.trace {
		defer un(trace("ParseArrayLiteral"))
	}

	expr := &ast.ArrayLiteral{Token: p.tok}
	expr.Elements = p.parseExpressionList(token.RBRACK)

	if !p.expect(token.RBRACK) {
		return nil
	}

	return expr
}

func (p *Parser) parseHashLiteral() ast.Expression {
	if p.trace {
		defer un(trace("ParseHashLiteral"))
	}

	expr := &ast.HashLiteral{Token: p.tok}
	expr.Keys = []ast.Expression{}
	expr.Pairs = map[ast.Expression]ast.Expression{}

	if p.peek1().Type != token.RBRACE {
		p.next()

		key := p.parseExpression(ASSIGNMENT - 1) // right associativity
		if !p.expect(token.COLON) {
			return nil
		}

		p.next()

		value := p.parseExpression(ASSIGNMENT - 1) // right associativity

		expr.Keys = append(expr.Keys, key)
		expr.Pairs[key] = value

		for p.peek1().Type != token.RBRACE {
			if !p.expect(token.COMMA) {
				return nil
			}

			p.next()

			key := p.parseExpression(ASSIGNMENT - 1) // right associativity
			if !p.expect(token.COLON) {
				return nil
			}

			p.next()

			value := p.parseExpression(ASSIGNMENT - 1) // right associativity

			expr.Keys = append(expr.Keys, key)
			expr.Pairs[key] = value
		}
	}

	if !p.expect(token.RBRACE) {
		return nil
	}

	return expr
}

func (p *Parser) parseNullLiteral() ast.Expression {
	if p.trace {
		defer un(trace("ParseNullLiteral"))
	}

	return &ast.NullLiteral{Token: p.tok}
}

func (p *Parser) reportIllegalToken() ast.Expression {
	p.error("illegal token: <%s>", p.tok.Literal)
	return nil
}

func (p *Parser) parseIdentifierList(terminator token.TokenType) []*ast.Identifier {
	idents := []*ast.Identifier{}
	if p.peek1().Type == terminator {
		return idents
	}

	p.next()

	ident := p.parseIdentifier().(*ast.Identifier)
	if ident != nil {
		idents = append(idents, ident)
	}

	for p.peek1().Type != token.RPAREN {
		if !p.expect(token.COMMA) {
			return nil
		}

		p.next()

		ident := p.parseIdentifier().(*ast.Identifier)
		if ident != nil {
			idents = append(idents, ident)
		}
	}

	return idents
}

func (p *Parser) parseExpressionList(terminator token.TokenType) []ast.Expression {
	exprs := []ast.Expression{}
	if p.peek1().Type == terminator {
		return exprs
	}

	p.next()

	expr := p.parseExpression(ASSIGNMENT - 1) // right associativity
	if expr != nil {
		exprs = append(exprs, expr)
	}

	for p.peek1().Type != terminator {
		if !p.expect(token.COMMA) {
			return nil
		}

		p.next()

		expr := p.parseExpression(ASSIGNMENT - 1) // right associativity
		if expr != nil {
			exprs = append(exprs, expr)
		}
	}

	return exprs
}

func (p *Parser) skip(toks ...token.TokenType) {
	for slices.Contains(toks, p.peek1().Type) {
		p.next()
	}
}

func (p *Parser) expect(ttype token.TokenType) bool {
	if p.peek1().Type == ttype {
		p.next()
		return true
	}
	p.error(
		"expected next token to be <%s> but was <%s>",
		ttype,
		p.peek1().Type,
	)
	return false
}

func (p *Parser) peek0() token.Token {
	if p.tok.Type == token.EOF {
		return p.tok
	}
	return p.l.Token(p.current)
}

func (p *Parser) peek1() token.Token {
	if p.tok.Type == token.EOF {
		return p.tok
	}
	return p.l.Token(p.current + 1)
}

func (p *Parser) next() {
	if p.tok.Type == token.EOF {
		return
	}
	p.current += 1
	p.tok = p.l.Token(p.current)
}

func (p *Parser) getRule(ttype token.TokenType) parseRule {
	return p.rules[ttype]
}

func (p *Parser) error(format string, args ...any) {
	p.errors = append(p.errors, fmt.Sprintf(format, args...))
}
