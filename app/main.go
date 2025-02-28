package main

import (
	"fmt"
	"os"
)

const (
	LEFT_PAREN  rune = '('
	RIGHT_PAREN rune = ')'
	LEFT_BRACE  rune = '{'
	RIGHT_BRACE rune = '}'
	STAR        rune = '*'
	DOT         rune = '.'
	COMMA       rune = ','
	PLUS        rune = '+'
	MINUS       rune = '-'
	SEMICOLON   rune = ';'
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

	fileContents := string(rawfileContents)
	lines := 1
	for _, current := range fileContents {
		switch current {
		case '\n':
			lines++
		case LEFT_PAREN:
			fmt.Println("LEFT_PAREN ( null")
		case RIGHT_PAREN:
			fmt.Println("RIGHT_PAREN ) null")
		case LEFT_BRACE:
			fmt.Println("LEFT_BRACE { null")
		case RIGHT_BRACE:
			fmt.Println("RIGHT_BRACE } null")
		case STAR:
			fmt.Println("STAR * null")
		case DOT:
			fmt.Println("DOT . null")
		case COMMA:
			fmt.Println("COMMA , null")
		case PLUS:
			fmt.Println("PLUS + null")
		case MINUS:
			fmt.Println("MINUS - null")
		case SEMICOLON:
			fmt.Println("SEMICOLON ; null")
		default:
			fmt.Printf("[line %d] Error: Unexpected character: %c\n", lines, current)
		}
	}
	fmt.Println("EOF  null")
}
