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

				c := compiler.New()
				chunk, err := c.Compile(program)
				if err != nil {
					t.Fatalf("test[%d] - Compile(program) ==> expected: not <%#v>", i, err)
				}

				vm := New(chunk, true)
				err = vm.Run()
				if err != nil {
					t.Fatalf("test[%d] - vm.Run() ==> expected: not <%#v>", i, err)
				}

				switch expected := test.top.(type) {
				case object.Number:
					if expected != vm.StackTop() {
						t.Fatalf("test[%d] - vm.StackTop() ==> expected: <%f> but was: <%f>", i, expected, vm.StackTop())
					}
				default:
					t.Fatalf("test[%d] ==> unexpected constant type: %T", i, test.top)
				}
			}
		})
	}
}
