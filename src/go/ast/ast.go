package ast

import (
	"bytes"
	"strings"

	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/token"
)

type NodeType int

const (
	PROGRAM NodeType = iota
	ERROR
	LET_DECLARATION
	RETURN_STATEMENT
	EXPRESSION_STATEMENT
	BLOCK_STATEMENT
	UNARY_EXPRESSION
	BINARY_EXPRESSION
	LOGICAL_EXPRESSION
	CONDITIONAL_EXPRESSION
	FUNCTION_LITERAL
	CALL_EXPRESSION
	ASSIGNMENT_EXPRESSION
	SUBSCRIPT_EXPRESSION
	IDENTIFIER
	NUMBER_LITERAL
	BOOLEAN_LITERAL
	STRING_LITERAL
	ARRAY_LITERAL
	HASH_LITERAL
	NULL_LITERAL
)

type Node interface {
	TokenLiteral() string

	Type() NodeType
	String() string
}

type Program struct {
	Statements []Statement
}

func (e *Program) TokenLiteral() string {
	if len(e.Statements) <= 0 {
		return ""
	}
	return e.Statements[0].TokenLiteral()
}

func (e *Program) Type() NodeType {
	return PROGRAM
}

func (e *Program) String() string {
	var out bytes.Buffer
	for _, s := range e.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type Error struct {
	Token   token.Token
	Message string
}

func (e *Error) TokenLiteral() string {
	return e.Token.Literal
}

func (e *Error) Type() NodeType {
	return ERROR
}

func (e *Error) String() string {
	return e.Message
}

type Statement interface {
	Node
	statementNode()
}

func (ls *LetDeclaration) statementNode()      {}
func (rs *ReturnStatement) statementNode()     {}
func (es *ExpressionStatement) statementNode() {}
func (bs *BlockStatement) statementNode()      {}

func (e *Error) statementNode() {}

type LetDeclaration struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetDeclaration) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetDeclaration) Type() NodeType {
	return LET_DECLARATION
}

func (ls *LetDeclaration) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(ls.Name.String())
	out.WriteString("=")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")

	return out.String()
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) Type() NodeType {
	return RETURN_STATEMENT
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral())
	out.WriteString(" ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) Type() NodeType {
	return EXPRESSION_STATEMENT
}

func (es *ExpressionStatement) String() string {
	var out bytes.Buffer

	if es.Expression != nil {
		out.WriteString(es.Expression.String())
	}
	out.WriteString(";")

	return out.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStatement) Type() NodeType {
	return BLOCK_STATEMENT
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	out.WriteString("}")
	return out.String()
}

type Expression interface {
	Node
	expressionNode()
}

func (ue *UnaryExpression) expressionNode()       {}
func (be *BinaryExpression) expressionNode()      {}
func (ce *ConditionalExpression) expressionNode() {}
func (fl *FunctionLiteral) expressionNode()       {}
func (ce *CallExpression) expressionNode()        {}
func (ae *AssignmentExpression) expressionNode()  {}
func (ie *SubscriptExpression) expressionNode()   {}
func (i *Identifier) expressionNode()             {}
func (nl *NumberLiteral) expressionNode()         {}
func (bl *BooleanLiteral) expressionNode()        {}
func (sl *StringLiteral) expressionNode()         {}
func (al *ArrayLiteral) expressionNode()          {}
func (hl *HashLiteral) expressionNode()           {}
func (nl *NullLiteral) expressionNode()           {}

func (e *Error) expressionNode() {}

type UnaryExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (ue *UnaryExpression) TokenLiteral() string {
	return ue.Token.Literal
}

func (ue *UnaryExpression) Type() NodeType {
	return UNARY_EXPRESSION
}

func (ue *UnaryExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ue.Operator)
	out.WriteString(ue.Right.String())
	out.WriteString(")")

	return out.String()
}

type BinaryExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (be *BinaryExpression) TokenLiteral() string {
	return be.Token.Literal
}

func (be *BinaryExpression) Type() NodeType {
	return BINARY_EXPRESSION
}

func (be *BinaryExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(be.Left.String())
	out.WriteString(be.Operator)
	out.WriteString(be.Right.String())
	out.WriteString(")")

	return out.String()
}

type LogicalExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (le *LogicalExpression) TokenLiteral() string {
	return le.Token.Literal
}

func (le *LogicalExpression) Type() NodeType {
	return LOGICAL_EXPRESSION
}

func (le *LogicalExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(le.Left.String())
	out.WriteString(le.Operator)
	out.WriteString(le.Right.String())
	out.WriteString(")")

	return out.String()
}

type ConditionalExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ce *ConditionalExpression) TokenLiteral() string {
	return ce.Token.Literal
}

func (ce *ConditionalExpression) Type() NodeType {
	return CONDITIONAL_EXPRESSION
}

func (ce *ConditionalExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ce.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(ce.Condition.String())
	out.WriteString(" ")
	out.WriteString(ce.Consequence.String())
	if ce.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(ce.Alternative.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FunctionLiteral) Type() NodeType {
	return FUNCTION_LITERAL
}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, param := range fl.Parameters {
		params = append(params, param.Value)
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     token.Token
	Callee    Expression
	Arguments []Expression
}

func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

func (ce *CallExpression) Type() NodeType {
	return CALL_EXPRESSION
}

func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, arg := range ce.Arguments {
		args = append(args, arg.String())
	}

	out.WriteString(ce.Callee.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ","))
	out.WriteString(")")

	return out.String()
}

type AssignmentExpression struct {
	Token  token.Token
	LValue Expression
	RValue Expression
}

func (ae *AssignmentExpression) TokenLiteral() string {
	return ae.Token.Literal
}

func (ae *AssignmentExpression) Type() NodeType {
	return ASSIGNMENT_EXPRESSION
}

func (ae *AssignmentExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ae.LValue.String())
	out.WriteString("=")
	out.WriteString(ae.RValue.String())
	out.WriteString(")")

	return out.String()
}

type SubscriptExpression struct {
	Token     token.Token
	Base      Expression
	Subscript Expression
}

func (ie *SubscriptExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *SubscriptExpression) Type() NodeType {
	return SUBSCRIPT_EXPRESSION
}

func (ie *SubscriptExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Base.String())
	out.WriteString("[")
	out.WriteString(ie.Subscript.String())
	out.WriteString("])")

	return out.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) Type() NodeType {
	return IDENTIFIER
}

func (i *Identifier) String() string {
	return i.Value
}

type NumberLiteral struct {
	Token token.Token
	Value float64
}

func (nl *NumberLiteral) TokenLiteral() string {
	return nl.Token.Literal
}

func (nl *NumberLiteral) Type() NodeType {
	return NUMBER_LITERAL
}

func (nl *NumberLiteral) String() string {
	return nl.Token.Literal
}

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) TokenLiteral() string {
	return bl.Token.Literal
}

func (bl *BooleanLiteral) Type() NodeType {
	return BOOLEAN_LITERAL
}

func (bl *BooleanLiteral) String() string {
	return bl.Token.Literal
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

func (sl *StringLiteral) Type() NodeType {
	return STRING_LITERAL
}

func (sl *StringLiteral) String() string {
	return sl.Token.Literal
}

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}

func (al *ArrayLiteral) Type() NodeType {
	return ARRAY_LITERAL
}

func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elems := []string{}
	for _, elem := range al.Elements {
		elems = append(elems, elem.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elems, ","))
	out.WriteString("]")

	return out.String()
}

type HashLiteral struct {
	Token token.Token
	Keys  []Expression
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) TokenLiteral() string {
	return hl.Token.Literal
}

func (hl *HashLiteral) Type() NodeType {
	return HASH_LITERAL
}

func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, key := range hl.Keys {
		pairs = append(pairs, key.String()+":"+hl.Pairs[key].String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ","))
	out.WriteString("}")

	return out.String()
}

type NullLiteral struct {
	Token token.Token
}

func (nl *NullLiteral) TokenLiteral() string {
	return nl.Token.Literal
}

func (nl *NullLiteral) Type() NodeType {
	return NULL_LITERAL
}

func (nl *NullLiteral) String() string {
	return nl.Token.Literal
}

var nodes = [...]string{
	PROGRAM:                "PROGRAM",
	ERROR:                  "ERROR",
	LET_DECLARATION:        "LET_DECLARATION",
	RETURN_STATEMENT:       "RETURN_STATEMENT",
	EXPRESSION_STATEMENT:   "EXPRESSION_STATEMENT",
	BLOCK_STATEMENT:        "BLOCK_STATEMENT",
	UNARY_EXPRESSION:       "UNARY_EXPRESSION",
	BINARY_EXPRESSION:      "BINARY_EXPRESSION",
	LOGICAL_EXPRESSION:     "LOGICAL_EXPRESSION",
	CONDITIONAL_EXPRESSION: "CONDITIONAL_EXPRESSION",
	FUNCTION_LITERAL:       "FUNCTION_LITERAL",
	CALL_EXPRESSION:        "CALL_EXPRESSION",
	ASSIGNMENT_EXPRESSION:  "ASSIGNMENT_EXPRESSION",
	SUBSCRIPT_EXPRESSION:   "SUBSCRIPT_EXPRESSION",
	IDENTIFIER:             "IDENTIFIER",
	NUMBER_LITERAL:         "NUMBER_LITERAL",
	BOOLEAN_LITERAL:        "BOOLEAN_LITERAL",
	STRING_LITERAL:         "STRING_LITERAL",
	ARRAY_LITERAL:          "ARRAY_LITERAL",
	HASH_LITERAL:           "HASH_LITERAL",
	NULL_LITERAL:           "NULL_LITERAL",
}

func (nt NodeType) String() string {
	return nodes[nt]
}
