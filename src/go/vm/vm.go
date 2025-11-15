package vm

import (
	"fmt"

	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/code"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/compiler"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/object"
)

const STACK_SIZE = 2045

type VM struct {
	chunk *compiler.Chunk
	trace bool

	sp    int
	stack []object.Object
	ip    int
}

func New(chunk *compiler.Chunk, trace bool) *VM {
	return &VM{
		chunk: chunk,
		trace: trace,

		stack: make([]object.Object, STACK_SIZE),
	}
}

func (vm *VM) Run() error {
	for vm.ip = 0; vm.ip < len(vm.chunk.Bytecode); vm.ip++ {
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
	}

	return nil
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= STACK_SIZE {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp += 1

	return nil
}

func (vm *VM) pop() object.Object {
	vm.sp -= 1
	return vm.stack[vm.sp]
}
