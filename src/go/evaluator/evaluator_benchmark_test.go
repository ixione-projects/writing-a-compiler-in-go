package evaluator

import (
	"testing"

	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/object"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/parser"
)

func BenchmarkEvaluate(b *testing.B) {
	for _, suite := range suites {
		b.Run(suite.name, func(b *testing.B) {
			for b.Loop() {
				for i, test := range suite.tests {
					benchmarkEvaluator(b, i, test)
				}
			}
		})
	}
}

func benchmarkEvaluator(tb testing.TB, i int, test EvaluatorTest) {
	p := parser.NewParser(test.input, false)
	program := p.ParseProgram()

	if 0 != len(p.Errors()) {
		tb.Errorf("test[%d] - len(p.Errors()) ==> expected: <%d> but was: <%d>", i, 0, len(p.Errors()))
		for j, msg := range p.Errors() {
			tb.Errorf("--------- p.Errors()[%d]: %s", j, msg)
		}
		tb.FailNow()
	}

	_, interrupt := Evaluate(program, object.NewEnvironment(nil))
	if interrupt != nil {
		tb.Errorf("test[%d] - Evaluate() (*object.Error) ==> expected: <%#v> but was: <%s>", i, nil, interrupt.(*object.Error).Message)
		tb.Fatalf("test[%d] - <%s>", i, program.String())
	}
}
