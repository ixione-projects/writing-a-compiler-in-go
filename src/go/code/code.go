package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type OpCode byte

const (
	OP_CONSTANT OpCode = iota
	OP_ADD
)

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[OpCode]*Definition{
	OP_CONSTANT: {"OP_CONSTANT", []int{2}},
	OP_ADD:      {"OP_ADD", []int{}},
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[OpCode(op)]
	if !ok {
		return nil, fmt.Errorf("unexpected opcode: %d", op)
	}
	return def, nil
}

type Bytecode []byte
type Instruction []byte

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

func Disassemble(ins Bytecode) string {
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

		fmt.Fprintf(&out, "%04d %s\n", i, Instruction(ins[i:i+read]).String())

		i += read
	}

	return out.String()
}
