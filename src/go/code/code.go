package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type OpCode byte

const (
	OP_CONSTANT OpCode = iota
	OP_POP
	OP_ADD
	OP_SUB
	OP_MUL
	OP_DIV
	OP_EQUAL
	OP_NOT_EQUAL
	OP_GREATER
	OP_MINUS
	OP_BANG
	OP_JUMP_NOT_TRUTHY
	OP_JUMP
	OP_TRUE
	OP_FALSE
	OP_NULL
)

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[OpCode]*Definition{
	OP_CONSTANT:        {"OP_CONSTANT", []int{2}},
	OP_POP:             {"OP_POP", []int{}},
	OP_ADD:             {"OP_ADD", []int{}},
	OP_SUB:             {"OP_SUB", []int{}},
	OP_MUL:             {"OP_MUL", []int{}},
	OP_DIV:             {"OP_DIV", []int{}},
	OP_EQUAL:           {"OP_EQUAL", []int{}},
	OP_NOT_EQUAL:       {"OP_NOT_EQUAL", []int{}},
	OP_GREATER:         {"OP_GREATER", []int{}},
	OP_MINUS:           {"OP_MINUS", []int{}},
	OP_BANG:            {"OP_BANG", []int{}},
	OP_JUMP_NOT_TRUTHY: {"OP_JUMP_NOT_TRUTHY", []int{2}},
	OP_JUMP:            {"OP_JUMP", []int{2}},
	OP_TRUE:            {"OP_TRUE", []int{}},
	OP_FALSE:           {"OP_FALSE", []int{}},
	OP_NULL:            {"OP_NULL", []int{}},
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[OpCode(op)]
	if !ok {
		return nil, fmt.Errorf("unexpected opcode: %s", OpCode(op))
	}
	return def, nil
}

type Bytecode []byte
type Instruction []byte

func (ins Bytecode) Disassemble() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		read := 1
		for _, width := range def.OperandWidths {
			read += width
		}

		fmt.Fprintf(&out, "%04d %s\n", i, Instruction(ins[i:i+read]))

		i += read
	}

	return out.String()
}

func (ins Instruction) String() string {
	def, err := Lookup(byte(ins[0]))
	if err != nil {
		return fmt.Sprintf("ERROR: %s\n", err)
	}

	operands, _ := ReadOperands(def, ins[1:])
	count := len(def.OperandWidths)
	if len(operands) != count {
		return fmt.Sprintf("ERROR: unexpected operand len, expecteed %d but was %d\n", count, len(operands))
	}

	switch count {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unexpected instruction type %s\n", def.Name)
}

func Make(op OpCode, operands ...int) Instruction {
	def, ok := definitions[op]
	if !ok {
		return Instruction{}
	}

	length := 1
	for _, w := range def.OperandWidths {
		length += w
	}

	instruction := make(Instruction, length)
	instruction[0] = byte(op)

	off := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[off:], uint16(o))
		}

		off += width
	}

	return instruction
}

func ReadOperands(def *Definition, ins Instruction) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))

	off := 0
	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[off:]))
		}

		off += width
	}

	return operands, off
}

func ReadUint16(ins []byte) uint16 {
	return binary.BigEndian.Uint16(ins)
}

func Concat(ins []Instruction) Bytecode {
	bytecode := Bytecode{}
	for _, in := range ins {
		bytecode = append(bytecode, in...)
	}
	return bytecode
}

var codes = [...]string{
	OP_CONSTANT:        "OP_CONSTANT",
	OP_POP:             "OP_POP",
	OP_ADD:             "OP_ADD",
	OP_SUB:             "OP_SUB",
	OP_MUL:             "OP_MUL",
	OP_DIV:             "OP_DIV",
	OP_EQUAL:           "OP_EQUAL",
	OP_NOT_EQUAL:       "OP_NOT_EQUAL",
	OP_GREATER:         "OP_GREATER",
	OP_MINUS:           "OP_MINUS",
	OP_BANG:            "OP_BANG",
	OP_JUMP_NOT_TRUTHY: "OP_JUMP_NOT_TRUTHY",
	OP_JUMP:            "OP_JUMP",
	OP_TRUE:            "OP_TRUE",
	OP_FALSE:           "OP_FALSE",
	OP_NULL:            "OP_NULL",
}

func (oc OpCode) String() string {
	return codes[oc]
}
