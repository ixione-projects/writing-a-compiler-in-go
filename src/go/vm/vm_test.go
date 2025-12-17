package vm

import (
	"testing"

	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/compiler"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/object"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/parser"
)

type VMTest struct {
	input string
	top   object.Object
}

func TestStackTop(t *testing.T) {
	suites := []struct {
		name  string
		tests []VMTest
	}{
		{
			"TestNumberExpression",
			[]VMTest{
				{
					input: "1",
					top:   object.Number(1),
				},
				{
					input: "2",
					top:   object.Number(2),
				},
				{
					input: "1 + 2",
					top:   object.Number(3),
				},
				{
					input: "1 - 2",
					top:   object.Number(-1),
				},
				{
					input: "1 * 2",
					top:   object.Number(2),
				},
				{
					input: "4 / 2",
					top:   object.Number(2),
				},
				{
					input: "50 / 2 * 2 + 10 - 5",
					top:   object.Number(55),
				},
				{
					input: "5 + 5 + 5 + 5 - 10",
					top:   object.Number(10),
				},
				{
					input: "2 * 2 * 2 * 2 * 2",
					top:   object.Number(32),
				},
				{
					input: "5 * 2 + 10",
					top:   object.Number(20),
				},
				{
					input: "5 + 2 * 10",
					top:   object.Number(25),
				},
				{
					input: "5 * (2 + 10)",
					top:   object.Number(60),
				},
				{
					input: "-5",
					top:   object.Number(-5),
				},
				{
					input: "-10",
					top:   object.Number(-10),
				},
				{
					input: "-50 + 100 + -50",
					top:   object.Number(0),
				},
				{
					input: "(5 + 10 * 2 + 15 / 3) * 2 + -10",
					top:   object.Number(50),
				},
			},
		},
		{
			"TestBooleanExpression",
			[]VMTest{
				{
					input: "true",
					top:   object.Boolean(true),
				},
				{
					input: "false",
					top:   object.Boolean(false),
				},
				{
					input: "1 < 2",
					top:   object.Boolean(true),
				},
				{
					input: "1 > 2",
					top:   object.Boolean(false),
				},
				{
					input: "1 < 1",
					top:   object.Boolean(false),
				},
				{
					input: "1 > 1",
					top:   object.Boolean(false),
				},
				{
					input: "1 == 1",
					top:   object.Boolean(true),
				},
				{
					input: "1 != 1",
					top:   object.Boolean(false),
				},
				{
					input: "1 == 2",
					top:   object.Boolean(false),
				},
				{
					input: "1 != 2",
					top:   object.Boolean(true),
				},
				{
					input: "true == true",
					top:   object.Boolean(true),
				},
				{
					input: "false == false",
					top:   object.Boolean(true),
				},
				{
					input: "true == false",
					top:   object.Boolean(false),
				},
				{
					input: "true != false",
					top:   object.Boolean(true),
				},
				{
					input: "false != true",
					top:   object.Boolean(true),
				},
				{
					input: "(1 < 2) == true",
					top:   object.Boolean(true),
				},
				{
					input: "(1 < 2) == false",
					top:   object.Boolean(false),
				},
				{
					input: "(1 > 2) == true",
					top:   object.Boolean(false),
				},
				{
					input: "(1 > 2) == false",
					top:   object.Boolean(true),
				},
				{
					input: "!true",
					top:   object.Boolean(false),
				},
				{
					input: "!false",
					top:   object.Boolean(true),
				},
				{
					input: "!5",
					top:   object.Boolean(false),
				},
				{
					input: "!!true",
					top:   object.Boolean(true),
				},
				{
					input: "!!false",
					top:   object.Boolean(false),
				},
				{
					input: "!!5",
					top:   object.Boolean(true),
				},
				{
					input: "!(if (false) { 5; })",
					top:   object.Boolean(true),
				},
			},
		},
		{
			"TestConditionalExpression",
			[]VMTest{
				{
					input: "if (true) { 10 }",
					top:   object.Number(10),
				},
				{
					input: "if (true) { 10 } else { 20 }",
					top:   object.Number(10),
				},
				{
					input: "if (false) { 10 } else { 20 }",
					top:   object.Number(20),
				},
				{
					input: "if (1) { 10 }",
					top:   object.Number(10),
				},
				{
					input: "if (1 < 2) { 10 }",
					top:   object.Number(10),
				},
				{
					input: "if (1 < 2) { 10 } else { 20 }",
					top:   object.Number(10),
				},
				{
					input: "if (1 > 2) { 10 } else { 20 }",
					top:   object.Number(20),
				},
				{
					input: "if (1 > 2) { 10 }",
					top:   &object.Null{},
				},
				{
					input: "if (false) { 10 }",
					top:   &object.Null{},
				},
				{
					input: "if ((if (false) { 10 })) { 10 } else { 20 }",
					top:   object.Number(20),
				},
			},
		},
		{
			"TestLetStatement",
			[]VMTest{
				{
					input: "let one = 1; one",
					top:   object.Number(1),
				},
				{
					input: "let one = 1; let two = 2; one + two",
					top:   object.Number(3),
				},
				{
					input: "let one = 1; let two = one + one; one + two",
					top:   object.Number(3),
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

				c := compiler.NewCompiler(program)
				chunk, err := c.CompileProgram()
				if err != nil {
					t.Fatalf("test[%d] - Compile(program) ==> expected: not <%#v>", i, err)
				}

				vm := NewVM(chunk, true)
				err = vm.Run()
				if err != nil {
					t.Fatalf("test[%d] - vm.Run() ==> expected: not <%#v>", i, err)
				}

				switch expected := test.top.(type) {
				case object.Number:
					if expected != vm.LastPopInstruction() {
						t.Fatalf("test[%d] - vm.StackTop() ==> expected: <%f> but was: <%f>", i, expected, vm.LastPopInstruction())
					}
				case object.Boolean:
					if expected != vm.LastPopInstruction() {
						t.Fatalf("test[%d] - vm.StackTop() ==> expected: <%t> but was: <%t>", i, expected, vm.LastPopInstruction())
					}
				case *object.Null:
					if _, ok := vm.LastPopInstruction().(*object.Null); !ok {
						t.Fatalf("test[%d] - vm.StackTop() ==> unexpected type, expected: <%T> but was: <%T>", i, &object.Null{}, vm.LastPopInstruction())
					}
				default:
					t.Fatalf("test[%d] ==> unexpected constant type: %T", i, test.top)
				}
			}
		})
	}
}
