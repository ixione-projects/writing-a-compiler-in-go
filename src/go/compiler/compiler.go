package compiler

import (
	"fmt"

	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/ast"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/code"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/object"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/util"
)

type Compiler struct {
	items *util.Stack[CompilerStackItem]

	chunk *Chunk
}

type StackItemType int

const (
	NODE_ITEM StackItemType = iota
	OP_CODE_ITEM
)

type CompilerStackItem interface {
	Type() StackItemType
}

type NodeItem struct {
	node ast.Node
}

func (n *NodeItem) Type() StackItemType {
	return NODE_ITEM
}

type OpCodeItem code.OpCode

func (o OpCodeItem) Type() StackItemType {
	return OP_CODE_ITEM
}

type Chunk struct {
	Bytecode  code.Bytecode
	Constants []object.Object
}

func NewCompiler(node ast.Node) *Compiler {
	nodes := util.NewStack[CompilerStackItem](util.INITIAL_STACK_CAPACITY)
	nodes.Push(&NodeItem{node})

	return &Compiler{
		items: nodes,
		chunk: &Chunk{
			Bytecode:  code.Bytecode{},
			Constants: []object.Object{},
		},
	}
}

func (c *Compiler) Compile() (*Chunk, error) {
	for c.items.Size() != 0 {
		item := c.items.Pop()

		switch item.Type() {
		case NODE_ITEM:
			node := item.(*NodeItem).node

			switch node.Type() {
			case ast.PROGRAM:
				node := node.(*ast.Program)
				for i := len(node.Statements) - 1; i >= 0; i-- {
					c.items.Push(&NodeItem{node.Statements[i]})
				}
			case ast.ERROR:
				return nil, fmt.Errorf("ERROR: %s\n", node.(*ast.Error).Message)
			case ast.LET_DECLARATION:
			case ast.RETURN_STATEMENT:
			case ast.EXPRESSION_STATEMENT:
				c.items.Push(OpCodeItem(code.OP_POP))
				c.items.Push(&NodeItem{node.(*ast.ExpressionStatement).Expression})
			case ast.BLOCK_STATEMENT:
			case ast.UNARY_EXPRESSION:
				node := node.(*ast.UnaryExpression)
				switch node.Operator {
				case "-":
					c.items.Push(OpCodeItem(code.OP_MINUS))
				case "!":
					c.items.Push(OpCodeItem(code.OP_BANG))
				default:
					return nil, fmt.Errorf("unexpected unary operator: %s\n", node.Operator)
				}
				c.items.Push(&NodeItem{node.Right})
			case ast.BINARY_EXPRESSION:
				node := node.(*ast.BinaryExpression)
				if node.Operator == "<" {
					c.items.Push(OpCodeItem(code.OP_GREATER))
					c.items.Push(&NodeItem{node.Left})
					c.items.Push(&NodeItem{node.Right})
				} else {
					switch node.Operator {
					case "+":
						c.items.Push(OpCodeItem(code.OP_ADD))
					case "-":
						c.items.Push(OpCodeItem(code.OP_SUB))
					case "*":
						c.items.Push(OpCodeItem(code.OP_MUL))
					case "/":
						c.items.Push(OpCodeItem(code.OP_DIV))
					case "==":
						c.items.Push(OpCodeItem(code.OP_EQUAL))
					case "!=":
						c.items.Push(OpCodeItem(code.OP_NOT_EQUAL))
					case ">":
						c.items.Push(OpCodeItem(code.OP_GREATER))
					default:
						return nil, fmt.Errorf("unexpected binary operator: %s\n", node.Operator)
					}
					c.items.Push(&NodeItem{node.Right})
					c.items.Push(&NodeItem{node.Left})
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
				if node.(*ast.BooleanLiteral).Value {
					c.emit(code.OP_TRUE)
				} else {
					c.emit(code.OP_FALSE)
				}
			case ast.STRING_LITERAL:
			case ast.ARRAY_LITERAL:
			case ast.HASH_LITERAL:
			case ast.NULL_LITERAL:
			default:
				return nil, fmt.Errorf("unexpected node type: %T", node.Type())
			}
		case OP_CODE_ITEM:
			c.emit(code.OpCode(item.(OpCodeItem)))
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
