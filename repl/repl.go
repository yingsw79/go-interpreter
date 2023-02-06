package repl

import (
	"bufio"
	"fmt"
	"go-interpreter/lexer"
	"go-interpreter/parser"
	"io"
)

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		if !scanner.Scan() {
			return
		}

		t := scanner.Text()
		if t == "exit" {
			return
		}

		p := parser.NewParser(lexer.NewLexer(t))
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		fmt.Fprintln(out, program.String())
	}
}

func printParserErrors(out io.Writer, errors []error) {
	fmt.Fprintln(out, "parser errors:")

	for _, err := range errors {
		fmt.Fprintln(out, err)
	}
}
