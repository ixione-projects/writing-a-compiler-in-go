package compiler

import (
	"fmt"

	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/ast"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/code"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/object"
)

type Chunk struct {
	Bytecode  code.Bytecode
	Constants []object.Object
}

type Compiler struct {
	chunk *Chunk
}

func New() *Compiler {
	return &Compiler{
		chunk: &Chunk{
			Bytecode:  code.Bytecode{},
			Constants: []object.Object{},
		},
	}
}

func (c *Compiler) Compile(node ast.Node) (*Chunk, error) {
	switch node.Type() {
	case ast.PROGRAM:
		for _, stmt := range node.(*ast.Program).Statements {
			_, err := c.Compile(stmt)
			if err != nil {
				return nil, err
			}
		}
	case ast.ERROR:
		return nil, fmt.Errorf("ERROR: %s\n", node.(*ast.Error).Message)
	case ast.LET_DECLARATION:
	case ast.RETURN_STATEMENT:
	case ast.EXPRESSION_STATEMENT:
		return c.Compile(node.(*ast.ExpressionStatement).Expression)
	case ast.BLOCK_STATEMENT:
	case ast.UNARY_EXPRESSION:
	case ast.BINARY_EXPRESSION:
		node := node.(*ast.BinaryExpression)
		_, err := c.Compile(node.Left)
		if err != nil {
			return nil, err
		}
		_, err = c.Compile(node.Right)
		if err != nil {
			return nil, err
		}

		switch node.Operator {
		case "+":
			c.emit(code.OP_ADD)
		default:
			return nil, fmt.Errorf("unexpected binary operator: %s\n", node.Operator)
		}
	case ast.LOGICAL_EXPRESSION:
	case ast.CONDITIONAL_EXPRESSION:
	case ast.FUNCTION_LITERAL:
	case ast.ASSIGNMENT_EXPRESSION:
	case ast.CALL_EXPRESSION:
	case ast.SUBSCRIPT_EXPRESSION:
	case ast.IDENTIFIER:
	case ast.NUMBER_LITERAL:
		constant := object.Number(node.(*ast.NumberLiteral).Value)
		c.emit(code.OP_CONSTANT, c.makeConstant(constant))
	case ast.BOOLEAN_LITERAL:
	case ast.STRING_LITERAL:
	case ast.ARRAY_LITERAL:
	case ast.HASH_LITERAL:
	case ast.NULL_LITERAL:
	default:
		return nil, fmt.Errorf("unexpected node type: %T", node)
	}

	return c.chunk, nil
}

func (c *Compiler) makeConstant(constant object.Object) int {
	c.chunk.Constants = append(c.chunk.Constants, constant)
	return len(c.chunk.Constants) - 1
}

func (c *Compiler) emit(op code.OpCode, operands ...int) int {
	pos := len(c.chunk.Bytecode)
	c.chunk.Bytecode = append(c.chunk.Bytecode, code.Make(op, operands...)...)
	return pos
}
