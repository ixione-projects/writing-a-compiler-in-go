package compiler

import (
	"fmt"

	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/ast"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/code"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/object"
)

type Compiler struct {
	node ast.Node

	symbols *SymbolTable

	chunk *Chunk
}

type Chunk struct {
	Bytecode  code.Bytecode
	Constants []object.Object
}

func NewCompiler(node ast.Node) *Compiler {
	return &Compiler{
		node:    node,
		symbols: NewSymbolTable(),
		chunk: &Chunk{
			Bytecode:  code.Bytecode{},
			Constants: []object.Object{},
		},
	}
}

func NewCompilerWithState(program *ast.Program, symbols *SymbolTable, constants []object.Object) *Compiler {
	return &Compiler{
		node:    program,
		symbols: symbols,
		chunk: &Chunk{
			Bytecode:  code.Bytecode{},
			Constants: constants,
		},
	}
}

func (c *Compiler) CompileProgram() (*Chunk, error) {
	if c.node == nil {
		return c.chunk, nil
	}

	program := c.node.(*ast.Program)
	for _, stmt := range program.Statements {
		c.node = stmt
		err := c.compileStatement()
		if err != nil {
			return nil, err
		}
	}

	c.node = nil
	return c.chunk, nil
}

func (c *Compiler) compileStatement() error {
	switch c.node.Type() {
	case ast.LET_DECLARATION:
		err := c.compileLetStatement()
		if err != nil {
			return err
		}
	case ast.RETURN_STATEMENT:
	case ast.EXPRESSION_STATEMENT:
		err := c.compileExpressionStatement()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unexpected node type: %T", c.node.Type())
	}
	return nil
}

func (c *Compiler) compileLetStatement() error {
	node := c.node.(*ast.LetDeclaration)
	c.node = node.Value
	err := c.compileExpression()
	if err != nil {
		return err
	}
	symbol := c.symbols.Define(node.Name.Value)
	c.emit(code.OP_SET_GLOBAL, symbol.Index)
	return nil
}

func (c *Compiler) compileExpressionStatement() error {
	c.node = c.node.(*ast.ExpressionStatement).Expression
	err := c.compileExpression()
	if err != nil {
		return err
	}
	c.emit(code.OP_POP)
	return nil
}

func (c *Compiler) compileBlockStatement() error {
	block := c.node.(*ast.BlockStatement)
	for _, stmt := range block.Statements {
		c.node = stmt
		err := c.compileStatement()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Compiler) compileExpression() error {
	switch c.node.Type() {
	case ast.UNARY_EXPRESSION:
		err := c.compileUnaryExpression()
		if err != nil {
			return err
		}
	case ast.BINARY_EXPRESSION:
		err := c.compileBinaryExpression()
		if err != nil {
			return err
		}
	case ast.LOGICAL_EXPRESSION:
	case ast.CONDITIONAL_EXPRESSION:
		err := c.compileConditionalExpression()
		if err != nil {
			return err
		}
	case ast.FUNCTION_LITERAL:
	case ast.ASSIGNMENT_EXPRESSION:
	case ast.CALL_EXPRESSION:
	case ast.SUBSCRIPT_EXPRESSION:
	case ast.IDENTIFIER:
		err := c.compileIdentifier()
		if err != nil {
			return err
		}
	case ast.NUMBER_LITERAL:
		constant := object.Number(c.node.(*ast.NumberLiteral).Value)
		c.emit(code.OP_CONSTANT, c.makeConstant(constant))
	case ast.BOOLEAN_LITERAL:
		if c.node.(*ast.BooleanLiteral).Value {
			c.emit(code.OP_TRUE)
		} else {
			c.emit(code.OP_FALSE)
		}
	case ast.STRING_LITERAL:
	case ast.ARRAY_LITERAL:
	case ast.HASH_LITERAL:
	case ast.NULL_LITERAL:
		c.emit(code.OP_NULL)
	default:
		return fmt.Errorf("unexpected node type: %T", c.node.Type())
	}
	return nil
}

func (c *Compiler) compileUnaryExpression() error {
	expr := c.node.(*ast.UnaryExpression)
	c.node = expr.Right
	err := c.compileExpression()
	if err != nil {
		return err
	}
	switch expr.Operator {
	case "-":
		c.emit(code.OP_MINUS)
	case "!":
		c.emit(code.OP_BANG)
	default:
		return fmt.Errorf("unexpected unary operator: %s\n", expr.Operator)
	}
	return nil
}

func (c *Compiler) compileBinaryExpression() error {
	expr := c.node.(*ast.BinaryExpression)
	if expr.Operator == "<" {
		c.node = expr.Right
		err := c.compileExpression()
		if err != nil {
			return err
		}
		c.node = expr.Left
		err = c.compileExpression()
		if err != nil {
			return err
		}
		c.emit(code.OP_GREATER)
	} else {
		c.node = expr.Left
		err := c.compileExpression()
		if err != nil {
			return err
		}
		c.node = expr.Right
		err = c.compileExpression()
		if err != nil {
			return err
		}
		switch expr.Operator {
		case "+":
			c.emit(code.OP_ADD)
		case "-":
			c.emit(code.OP_SUB)
		case "*":
			c.emit(code.OP_MUL)
		case "/":
			c.emit(code.OP_DIV)
		case "==":
			c.emit(code.OP_EQUAL)
		case "!=":
			c.emit(code.OP_NOT_EQUAL)
		case ">":
			c.emit(code.OP_GREATER)
		default:
			return fmt.Errorf("unexpected binary operator: %s\n", expr.Operator)
		}
	}
	return nil
}

func (c *Compiler) compileConditionalExpression() error {
	expr := c.node.(*ast.ConditionalExpression)
	c.node = expr.Condition
	err := c.compileExpression()
	if err != nil {
		return err
	}

	jumpNotTruthyOffset := c.emit(code.OP_JUMP_NOT_TRUTHY, 0xFFFF)
	c.node = expr.Consequence
	err = c.compileBlockStatement()
	if err != nil {
		return err
	}

	if c.chunk.Bytecode[len(c.chunk.Bytecode)-1] == byte(code.OP_POP) {
		c.chunk.Bytecode = c.chunk.Bytecode[:len(c.chunk.Bytecode)-1]
	}

	jumpOffset := c.emit(code.OP_JUMP, 0xFFFF)

	jump := len(c.chunk.Bytecode) - jumpNotTruthyOffset - 2
	c.patchJumpInstruction(jumpNotTruthyOffset, code.Make(code.OP_JUMP_NOT_TRUTHY, jump-1))

	if expr.Alternative == nil {
		c.emit(code.OP_NULL)
	} else {
		c.node = expr.Alternative
		err = c.compileBlockStatement()
		if err != nil {
			return err
		}

		if c.chunk.Bytecode[len(c.chunk.Bytecode)-1] == byte(code.OP_POP) {
			c.chunk.Bytecode = c.chunk.Bytecode[:len(c.chunk.Bytecode)-1]
		}
	}

	jump = len(c.chunk.Bytecode) - jumpOffset - 2
	c.patchJumpInstruction(jumpOffset, code.Make(code.OP_JUMP, jump-1))
	return nil
}

func (c *Compiler) compileIdentifier() error {
	ident := c.node.(*ast.Identifier)
	symbol, ok := c.symbols.Lookup(ident.Value)
	if !ok {
		return fmt.Errorf("unknown identifier: %s", ident.Value)
	}
	c.emit(code.OP_GET_GLOBAL, symbol.Index)
	return nil
}

func (c *Compiler) patchJumpInstruction(pos int, ins code.Instruction) {
	for i, b := range ins {
		c.chunk.Bytecode[pos+i] = b
	}
}

func (c *Compiler) makeConstant(constant object.Object) int {
	c.chunk.Constants = append(c.chunk.Constants, constant)
	return len(c.chunk.Constants) - 1
}

func (c *Compiler) emit(op code.OpCode, operands ...int) int {
	offset := len(c.chunk.Bytecode)
	c.chunk.Bytecode = append(c.chunk.Bytecode, code.Make(op, operands...)...)
	return offset
}
