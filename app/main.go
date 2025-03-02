package main

import (
	"fmt"
	"github.com/codecrafters-io/interpreter-starter-go/app/scanner"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	rawfileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	tokenScanner := scanner.NewScanner(rawfileContents)
	tokenScanner.Scan()
	fmt.Println("EOF  null")
	os.Exit(tokenScanner.ExitCode())
}
