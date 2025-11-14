package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/evaluator"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/object"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/parser"
)

const MONKEY_FACE = `
            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \ |'  | |
 | \   \  \ 0 | 0 / /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

const PROMPT = "> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment(nil)

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		p := parser.NewParser(line, false)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			io.WriteString(out, MONKEY_FACE)
			io.WriteString(out, "Woops! We ran into some monkey business here!\n")
			io.WriteString(out, "parser errors:\n")
			for _, msg := range p.Errors() {
				io.WriteString(out, "\t"+msg+"\n")
			}
			continue
		}

		value, error := evaluator.Evaluate(program, env)
		if error != nil {
			io.WriteString(out, error.Inspect()+"\n")
			continue
		}

		if value != nil && value.Type() != object.NULL {
			io.WriteString(out, value.Inspect()+"\n")
		}
	}
}
