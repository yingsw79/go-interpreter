package main

import (
	"fmt"
	"go-interpreter/repl"
	"os"
)

func main() {
	fmt.Println("Welcome to YL")
	repl.Start(os.Stdin, os.Stdout)
}
