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
	suites := []struct {
		name  string
		tests []CompilerTest
	}{
		{
			name: "TestNumberExpression",
			tests: []CompilerTest{
				{
					input:     `1 + 2`,
					constants: []object.Object{object.Number(1), object.Number(2)},
					instructions: []code.Instruction{
						code.Make(code.OP_CONSTANT, 0),
						code.Make(code.OP_CONSTANT, 1),
						code.Make(code.OP_ADD),
						code.Make(code.OP_POP),
					},
				},
				{
					input:     `1; 2`,
					constants: []object.Object{object.Number(1), object.Number(2)},
					instructions: []code.Instruction{
						code.Make(code.OP_CONSTANT, 0),
						code.Make(code.OP_POP),
						code.Make(code.OP_CONSTANT, 1),
						code.Make(code.OP_POP),
					},
				},
				{
					input:     `1 - 2`,
					constants: []object.Object{object.Number(1), object.Number(2)},
					instructions: []code.Instruction{
						code.Make(code.OP_CONSTANT, 0),
						code.Make(code.OP_CONSTANT, 1),
						code.Make(code.OP_SUB),
						code.Make(code.OP_POP),
					},
				},
				{
					input:     `1 * 2`,
					constants: []object.Object{object.Number(1), object.Number(2)},
					instructions: []code.Instruction{
						code.Make(code.OP_CONSTANT, 0),
						code.Make(code.OP_CONSTANT, 1),
						code.Make(code.OP_MUL),
						code.Make(code.OP_POP),
					},
				},
				{
					input:     `2 / 1`,
					constants: []object.Object{object.Number(2), object.Number(1)},
					instructions: []code.Instruction{
						code.Make(code.OP_CONSTANT, 0),
						code.Make(code.OP_CONSTANT, 1),
						code.Make(code.OP_DIV),
						code.Make(code.OP_POP),
					},
				},
				{
					input:     `-1`,
					constants: []object.Object{object.Number(1)},
					instructions: []code.Instruction{
						code.Make(code.OP_CONSTANT, 0),
						code.Make(code.OP_MINUS),
						code.Make(code.OP_POP),
					},
				},
			},
		},
		{
			name: "TestBooleanExpression",
			tests: []CompilerTest{
				{
					input:     `true`,
					constants: []object.Object{},
					instructions: []code.Instruction{
						code.Make(code.OP_TRUE),
						code.Make(code.OP_POP),
					},
				},
				{
					input:     `false`,
					constants: []object.Object{},
					instructions: []code.Instruction{
						code.Make(code.OP_FALSE),
						code.Make(code.OP_POP),
					},
				},
				{
					input:     `1 > 2`,
					constants: []object.Object{object.Number(1), object.Number(2)},
					instructions: []code.Instruction{
						code.Make(code.OP_CONSTANT, 0),
						code.Make(code.OP_CONSTANT, 1),
						code.Make(code.OP_GREATER),
						code.Make(code.OP_POP),
					},
				},
				{
					input:     `1 < 2`,
					constants: []object.Object{object.Number(2), object.Number(1)},
					instructions: []code.Instruction{
						code.Make(code.OP_CONSTANT, 0),
						code.Make(code.OP_CONSTANT, 1),
						code.Make(code.OP_GREATER),
						code.Make(code.OP_POP),
					},
				},
				{
					input:     `1 == 2`,
					constants: []object.Object{object.Number(1), object.Number(2)},
					instructions: []code.Instruction{
						code.Make(code.OP_CONSTANT, 0),
						code.Make(code.OP_CONSTANT, 1),
						code.Make(code.OP_EQUAL),
						code.Make(code.OP_POP),
					},
				},
				{
					input:     `1 != 2`,
					constants: []object.Object{object.Number(1), object.Number(2)},
					instructions: []code.Instruction{
						code.Make(code.OP_CONSTANT, 0),
						code.Make(code.OP_CONSTANT, 1),
						code.Make(code.OP_NOT_EQUAL),
						code.Make(code.OP_POP),
					},
				},
				{
					input:     `true == false`,
					constants: []object.Object{},
					instructions: []code.Instruction{
						code.Make(code.OP_TRUE),
						code.Make(code.OP_FALSE),
						code.Make(code.OP_EQUAL),
						code.Make(code.OP_POP),
					},
				},
				{
					input:     `true != false`,
					constants: []object.Object{},
					instructions: []code.Instruction{
						code.Make(code.OP_TRUE),
						code.Make(code.OP_FALSE),
						code.Make(code.OP_NOT_EQUAL),
						code.Make(code.OP_POP),
					},
				},
				{
					input:     `!true`,
					constants: []object.Object{},
					instructions: []code.Instruction{
						code.Make(code.OP_TRUE),
						code.Make(code.OP_BANG),
						code.Make(code.OP_POP),
					},
				},
			},
		},
		{
			name: "TestConditionalExpression",
			tests: []CompilerTest{
				{
					input: `
					if (true) { 10 }; 3333;
					`,
					constants: []object.Object{object.Number(10), object.Number(3333)},
					instructions: []code.Instruction{
						code.Make(code.OP_TRUE),
						code.Make(code.OP_JUMP_NOT_TRUTHY, 6),
						code.Make(code.OP_CONSTANT, 0),
						code.Make(code.OP_JUMP, 1),
						code.Make(code.OP_NULL),
						code.Make(code.OP_POP),
						code.Make(code.OP_CONSTANT, 1),
						code.Make(code.OP_POP),
					},
				},
				{
					input: `
					if (true) { 10 } else { 20 }; 3333;
					`,
					constants: []object.Object{object.Number(10), object.Number(20), object.Number(3333)},
					instructions: []code.Instruction{
						code.Make(code.OP_TRUE),
						code.Make(code.OP_JUMP_NOT_TRUTHY, 6),
						code.Make(code.OP_CONSTANT, 0),
						code.Make(code.OP_JUMP, 3),
						code.Make(code.OP_CONSTANT, 1),
						code.Make(code.OP_POP),
						code.Make(code.OP_CONSTANT, 2),
						code.Make(code.OP_POP),
					},
				},
			},
		},
		{
			name: "TestLetStatement",
			tests: []CompilerTest{
				{
					input: `
					let one = 1;
					let two = 2;
					`,
					constants: []object.Object{object.Number(1), object.Number(2)},
					instructions: []code.Instruction{
						code.Make(code.OP_CONSTANT, 0),
						code.Make(code.OP_SET_GLOBAL, 0),
						code.Make(code.OP_CONSTANT, 1),
						code.Make(code.OP_SET_GLOBAL, 1),
					},
				},
				{
					input: `
					let one = 1;
					one;
					`,
					constants: []object.Object{object.Number(1)},
					instructions: []code.Instruction{
						code.Make(code.OP_CONSTANT, 0),
						code.Make(code.OP_SET_GLOBAL, 0),
						code.Make(code.OP_GET_GLOBAL, 0),
						code.Make(code.OP_POP),
					},
				},
				{
					input: `
					let one = 1;
					let two = one;
					two;`,
					constants: []object.Object{object.Number(1)},
					instructions: []code.Instruction{
						code.Make(code.OP_CONSTANT, 0),
						code.Make(code.OP_SET_GLOBAL, 0),
						code.Make(code.OP_GET_GLOBAL, 0),
						code.Make(code.OP_SET_GLOBAL, 1),
						code.Make(code.OP_GET_GLOBAL, 1),
						code.Make(code.OP_POP),
					},
				},
			},
		},
	}

	for _, suite := range suites {
		t.Run(suite.name, func(t *testing.T) {
			for i, test := range suite.tests {
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
				chunk, err := c.CompileProgram()
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
					t.Fatalf("test[%d] - \n%s", i, chunk.Bytecode.Disassemble())
				}
			}
		})
	}
}
