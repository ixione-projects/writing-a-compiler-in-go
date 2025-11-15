package compiler

import (
	"testing"

	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/code"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/object"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/parser"
)

type CompilerTest struct {
	input        string
	constants    []object.Object
	instructions []code.Instruction
}

func TestCompile(t *testing.T) {
	tests := []CompilerTest{
		{
			input:     `1 + 2`,
			constants: []object.Object{object.Number(1), object.Number(2)},
			instructions: []code.Instruction{
				code.Make(code.OP_CONSTANT, 0),
				code.Make(code.OP_CONSTANT, 1),
				code.Make(code.OP_ADD),
			},
		},
	}

	for i, test := range tests {
		p := parser.NewParser(test.input, false)
		program := p.ParseProgram()

		if 0 != len(p.Errors()) {
			t.Errorf("test[%d] - len(p.Errors()) ==> expected: <%d> but was: <%d>", i, 0, len(p.Errors()))
			for j, msg := range p.Errors() {
				t.Errorf("--------- p.Errors()[%d]: %s", j, msg)
			}
			t.Fatalf("test[%d] - %s", i, test.input)
		}

		c := NewCompiler(program)
		chunk, err := c.Compile()
		if err != nil {
			t.Fatalf("test[%d] - Compile(program) ==> expected: not <%#v>", i, err)
		}

		fail := false
		for j, expected := range test.constants {
			switch expected := expected.(type) {
			case object.Number:
				if expected != chunk.Constants[j] {
					t.Errorf("test[%d] - Constants[%d] ==> expected: <%f> but was: <%f>", i, j, expected, chunk.Constants[j])
					fail = true
				}
			default:
				t.Errorf("test[%d] - Constants[%d] ==> unexpected constant type: %T", i, j, expected)
				fail = true
			}
		}

		bytecode := code.Concat(test.instructions)

		if len(bytecode) != len(chunk.Bytecode) {
			t.Errorf("test[%d] - len(bytecode.Instructions) ==> expected: <%d> but was: <%d>", i, len(bytecode), len(chunk.Bytecode))
			fail = true
		}

		for j, code := range bytecode {
			if code != chunk.Bytecode[j] {
				t.Errorf("test[%d] - Instructions[%d] ==> expected: <%d> but was: <%d>", i, j, code, chunk.Bytecode[j])
				fail = true
			}
		}

		if fail {
			t.FailNow()
		}
	}
}
