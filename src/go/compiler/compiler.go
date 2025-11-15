package compiler

import (
	"fmt"

	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/ast"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/code"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/object"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/util"
)

type Compiler struct {
	nodes *util.Stack[CompilerStackItem]

	chunk *Chunk
}

type StackItemType int

const (
	NODE_ITEM StackItemType = iota
	OPERATOR_ITEM
)

type CompilerStackItem interface {
	Type() StackItemType
}

func (n *NodeItem) Type() StackItemType {
	return NODE_ITEM
}
func (o OperatorItem) Type() StackItemType {
	return OPERATOR_ITEM
}

type NodeItem struct {
	node ast.Node
}

type OperatorItem string

type Chunk struct {
	Bytecode  code.Bytecode
	Constants []object.Object
}

func NewCompiler(node ast.Node) *Compiler {
	nodes := util.NewStack[CompilerStackItem](util.INITIAL_STACK_CAPACITY)
	nodes.Push(&NodeItem{node})

	return &Compiler{
		nodes: nodes,
		chunk: &Chunk{
			Bytecode:  code.Bytecode{},
			Constants: []object.Object{},
		},
	}
}

func (c *Compiler) Compile() (*Chunk, error) {
	for c.nodes.Size() != 0 {
		item := c.nodes.Pop()

		switch item.Type() {
		case NODE_ITEM:
			node := item.(*NodeItem).node

			switch node.Type() {
			case ast.PROGRAM:
				node := node.(*ast.Program)
				for i := len(node.Statements) - 1; i >= 0; i-- {
					c.nodes.Push(&NodeItem{node.Statements[i]})
				}
			case ast.ERROR:
				return nil, fmt.Errorf("ERROR: %s\n", node.(*ast.Error).Message)
			case ast.LET_DECLARATION:
			case ast.RETURN_STATEMENT:
			case ast.EXPRESSION_STATEMENT:
				c.nodes.Push(&NodeItem{node.(*ast.ExpressionStatement).Expression})
			case ast.BLOCK_STATEMENT:
			case ast.UNARY_EXPRESSION:
			case ast.BINARY_EXPRESSION:
				node := node.(*ast.BinaryExpression)
				c.nodes.Push(OperatorItem(node.Operator))
				c.nodes.Push(&NodeItem{node.Right})
				c.nodes.Push(&NodeItem{node.Left})
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
				return nil, fmt.Errorf("unexpected node type: %T", node.Type())
			}
		case OPERATOR_ITEM:
			switch item.(OperatorItem) {
			case "+":
				c.emit(code.OP_ADD)
			default:
				return nil, fmt.Errorf("unexpected operator: %s\n", item)
			}
		default:
			return nil, fmt.Errorf("unexpected item type: %T", item.Type())
		}
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
