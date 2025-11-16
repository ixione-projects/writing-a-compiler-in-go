package code

import "testing"

func TestMake(t *testing.T) {
	tests := []struct {
		op          OpCode
		operands    []int
		instruction Instruction
	}{
		{OP_CONSTANT, []int{65534}, Instruction{byte(OP_CONSTANT), 255, 254}},
		{OP_ADD, []int{}, Instruction{byte(OP_ADD)}},
	}

	for i, test := range tests {
		instruction := Make(test.op, test.operands...)

		if len(test.instruction) != len(instruction) {
			t.Errorf("test[%d] - len(instruction) ==> expected: <%d> but was: <%d>", i, len(test.instruction), len(instruction))
			continue
		}

		for j, expected := range test.instruction {
			if expected != instruction[j] {
				t.Errorf("test[%d] - instruction[%d] ==> expected: <%d> but was: <%d>", i, j, expected, instruction[j])
			}
		}
	}
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op       OpCode
		operands []int
		n        int
	}{
		{OP_CONSTANT, []int{65535}, 2},
	}

	for i, test := range tests {
		instruction := Make(test.op, test.operands...)

		def, err := Lookup(byte(test.op))
		if err != nil {
			t.Fatalf("test[%d] - Lookup(op) ==> expected: not <%#v>", i, err)
		}

		read, n := ReadOperands(def, instruction[1:])
		if test.n != n {
			t.Errorf("test[%d] - ReadOperands() ==> expected: <%d> but was: <%d>", i, test.n, n)
		}

		for j, expected := range test.operands {
			if expected != read[j] {
				t.Errorf("test[%d] - read[%d] ==> expected: <%d> but was: <%d>", i, j, expected, read[j])
			}
		}
	}
}

func TestDisassemble(t *testing.T) {
	tests := []struct {
		instructions []Instruction
		disassembled string
	}{
		{
			instructions: []Instruction{
				Make(OP_CONSTANT, 1),
				Make(OP_CONSTANT, 2),
				Make(OP_CONSTANT, 65535),
			},
			disassembled: `0000 OP_CONSTANT 1
0003 OP_CONSTANT 2
0006 OP_CONSTANT 65535
`,
		},
		{
			instructions: []Instruction{
				Make(OP_ADD),
				Make(OP_CONSTANT, 2),
				Make(OP_CONSTANT, 65535),
			},
			disassembled: `0000 OP_ADD
0001 OP_CONSTANT 2
0004 OP_CONSTANT 65535
`,
		},
	}

	for i, test := range tests {
		actual := Concat(test.instructions)
		if test.disassembled != actual.Disassemble() {
			t.Errorf("test[%d] - Disassemble() ==> expected: <%s> but was: <%s>", i, test.disassembled, actual.Disassemble())
		}
	}
}
