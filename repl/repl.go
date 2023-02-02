package repl

import (
	"bufio"
	"fmt"
	"go-interpreter/lexer"
	"go-interpreter/token"
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

		l := lexer.NewLexer(scanner.Text())

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
