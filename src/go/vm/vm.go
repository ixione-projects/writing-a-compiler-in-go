package vm

import (
	"fmt"

	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/code"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/compiler"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/object"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/util"
)

const STACK_SIZE = 2045

type VM struct {
	chunk *compiler.Chunk
	trace bool

	stack *util.Stack[object.Object]
	ip    int

	popped object.Object
}

func New(chunk *compiler.Chunk, trace bool) *VM {
	return &VM{
		chunk: chunk,
		trace: trace,

		stack: util.NewStack[object.Object](STACK_SIZE),
	}
}

const (
	TRUE  = object.Boolean(true)
	FALSE = object.Boolean(false)
)

var NULL = &object.Null{}

func (vm *VM) Run() error {
	vm.ip = 0
	for vm.ip < len(vm.chunk.Bytecode) {
		op := code.OpCode(vm.chunk.Bytecode[vm.ip])
		switch op {
		case code.OP_CONSTANT:
			index := code.ReadUint16(vm.chunk.Bytecode[vm.ip+1:])
			vm.ip += 2
			err := vm.push(vm.chunk.Constants[index])
			if err != nil {
				return err
			}
		case code.OP_MINUS, code.OP_BANG:
			err := vm.runUnaryExpressionInstruction(op)
			if err != nil {
				return err
			}
		case code.OP_ADD, code.OP_SUB, code.OP_MUL, code.OP_DIV, code.OP_EQUAL, code.OP_NOT_EQUAL, code.OP_GREATER:
			err := vm.runBinaryExpressionInstruction(op)
			if err != nil {
				return err
			}
		case code.OP_TRUE:
			err := vm.push(TRUE)
			if err != nil {
				return err
			}
		case code.OP_FALSE:
			err := vm.push(FALSE)
			if err != nil {
				return err
			}
		case code.OP_POP:
			vm.popped = vm.pop()
		default:
			return fmt.Errorf("unexpected instruction: %s\n", op)
		}

		vm.ip += 1
	}

	return nil
}

func (vm *VM) runUnaryExpressionInstruction(op code.OpCode) error {
	right := vm.pop()

	switch op {
	case code.OP_BANG:
		if isTruthy(right) {
			return vm.push(FALSE)
		}
		return vm.push(TRUE)
	case code.OP_MINUS:
		switch right.Type() {
		case object.NUMBER:
			return vm.push(object.Number(-right.(object.Number)))
		}
	}

	return fmt.Errorf("unknown operation: %s%s", op, right.Type())
}

func (vm *VM) runBinaryExpressionInstruction(op code.OpCode) error {
	right := vm.pop()
	left := vm.pop()

	switch {
	case left.Type() == object.NUMBER && right.Type() == object.NUMBER:
		switch op {
		case code.OP_ADD:
			return vm.push(left.(object.Number) + right.(object.Number))
		case code.OP_SUB:
			return vm.push(left.(object.Number) - right.(object.Number))
		case code.OP_MUL:
			return vm.push(left.(object.Number) * right.(object.Number))
		case code.OP_DIV:
			return vm.push(left.(object.Number) / right.(object.Number))
		case code.OP_EQUAL:
			return vm.push(toBoolean(left.(object.Number) == right.(object.Number)))
		case code.OP_NOT_EQUAL:
			return vm.push(toBoolean(left.(object.Number) != right.(object.Number)))
		case code.OP_GREATER:
			return vm.push(toBoolean(left.(object.Number) > right.(object.Number)))
		}
	case left.Type() != right.Type():
		return fmt.Errorf("unknown operation: %s %s %s", left.Type(), op, right.Type())
	case op == code.OP_EQUAL:
		return vm.push(toBoolean(left == right))
	case op == code.OP_NOT_EQUAL:
		return vm.push(toBoolean(left != right))
	}

	return fmt.Errorf("unknown operation: %s %s %s", left.Type(), op, right.Type())
}

func toBoolean(value bool) object.Boolean {
	if value {
		return TRUE
	}
	return FALSE
}

func isTruthy(o object.Object) object.Boolean {
	switch {
	case o == FALSE:
		return FALSE
	case o.Type() == object.NUMBER && o.(object.Number) == 0.0:
		return FALSE
	case o.Type() == object.STRING && o.(object.String) == "":
		return FALSE
	case o.Type() == object.ARRAY && len(o.(*object.Array).Elements) == 0:
		return FALSE
	case o.Type() == object.HASH && len(o.(*object.Hash).Pairs) == 0:
		return FALSE
	case o == NULL:
		return FALSE
	default:
		return TRUE
	}
}

func (vm *VM) LastPopInstruction() object.Object {
	return vm.popped
}

func (vm *VM) push(o object.Object) error {
	if vm.stack.Size() >= STACK_SIZE {
		return fmt.Errorf("stack overflow")
	}

	vm.stack.Push(o)
	return nil
}

func (vm *VM) pop() object.Object {
	return vm.stack.Pop()
}
