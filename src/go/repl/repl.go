package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/compiler"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/object"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/parser"
	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/vm"
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

	symbols := compiler.NewSymbolTable()
	constants := []object.Object{}
	globals := make([]object.Object, vm.GLOBALS_SIZE)
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

		c := compiler.NewCompilerWithState(program, symbols, constants)
		chunk, err := c.CompileProgram()
		if err != nil {
			fmt.Fprintf(out, "Woops! Compilation failed!\n")
			fmt.Fprintf(out, "compilation error: %s\n", err)
		}

		vm := vm.NewVMWithState(chunk, false, globals)
		err = vm.Run()
		if err != nil {
			fmt.Fprintf(out, "Woops! Executing bytecode failed!\n")
			fmt.Fprintf(out, "execution error: %s\n", err)
		}

		io.WriteString(out, vm.LastPopInstruction().Inspect()+"\n")
	}
}
