package repl

import (
	"bufio"
	"fmt"
	"go-interpreter/evaluator"
	"go-interpreter/lexer"
	"go-interpreter/object"
	"go-interpreter/parser"
	"io"
)

const PROMPT = ">>> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprint(out, PROMPT)
		if !scanner.Scan() {
			break
		}

		t := scanner.Text()

		p := parser.NewParser(lexer.NewLexer(t))
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		// fmt.Fprintln(out, program.String())
		res, err := evaluator.Eval(program, env)
		if err != nil {
			fmt.Fprintln(out, "Error:", err)
			continue
		}

		if res != nil && res.Type() != object.EXPLIST_OBJ {
			fmt.Fprintln(out, res.Inspect())
		}
	}
}

func printParserErrors(out io.Writer, errors []error) {
	fmt.Fprintln(out, "parser errors:")

	for _, err := range errors {
		fmt.Fprintln(out, err)
	}
}
