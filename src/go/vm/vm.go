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
}

func New(chunk *compiler.Chunk, trace bool) *VM {
	return &VM{
		chunk: chunk,
		trace: trace,

		stack: util.NewStack[object.Object](STACK_SIZE),
	}
}

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
		case code.OP_ADD:
			right := vm.pop()
			left := vm.pop()
			if left.Type() != object.NUMBER || right.Type() != object.NUMBER {
				return fmt.Errorf("unknown operation: %s + %s", left.Type(), right.Type())
			}
			vm.push(left.(object.Number) + right.(object.Number))
		default:
			return fmt.Errorf("unexpected instruction: %d\n", op)
		}

		vm.ip += 1
	}

	return nil
}

func (vm *VM) StackTop() object.Object {
	return vm.stack.Peek()
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
