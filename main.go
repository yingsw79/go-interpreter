package main

import (
	"fmt"
	"go-interpreter/repl"
	"os"
)

func main() {
	fmt.Println("Feel free to type in commands")
	repl.Start(os.Stdin, os.Stdout)
}
