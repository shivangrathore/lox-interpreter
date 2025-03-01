package scanner

import (
	"fmt"
	"os"
)

type Scanner struct {
	fileContents []byte
	currentIdx   int
	ExitCode     int
	lines        int
}

type Token struct {
	tokenType string
	lexeme    interface{}
}

func NewScanner(fileContents []byte) *Scanner {
	return &Scanner{fileContents: fileContents, currentIdx: 0, ExitCode: 0, lines: 1}
}

func NewToken(tokenType string, lexeme interface{}) *Token {
	return &Token{tokenType: tokenType, lexeme: lexeme}
}

func NextToken(s *Scanner) *Token {
	current := rune(s.fileContents[s.currentIdx])
	s.currentIdx++
	switch current {
	case '=':
		if len(s.fileContents) > s.currentIdx && rune(s.fileContents[s.currentIdx]) == '=' {
			s.currentIdx++
			return NewToken("EQUAL_EQUAL", "==")
		}
		return NewToken("EQUAL", "=")
	case '\n':
		s.lines++
	case '!':
		if len(s.fileContents) > s.currentIdx && rune(s.fileContents[s.currentIdx]) == '=' {
			s.currentIdx++
			return NewToken("BANG_EQUAL", "!=")
		}
		return NewToken("BANG", "!")
	case '.':
		return NewToken("DOT", ".")
	case ',':
		return NewToken("COMMA", ",")
	case '+':
		return NewToken("PLUS", "+")
	case '-':
		return NewToken("MINUS", "-")
	case ';':
		return NewToken("SEMICOLON", ";")
	case '(':
		return NewToken("LEFT_PAREN", "(")
	case ')':
		return NewToken("RIGHT_PAREN", ")")
	case '{':
		return NewToken("LEFT_BRACE", "{")
	case '}':
		return NewToken("RIGHT_BRACE", "}")
	case '*':
		return NewToken("STAR", "*")
	case '<':
		if len(s.fileContents) > s.currentIdx && rune(s.fileContents[s.currentIdx]) == '=' {
			s.currentIdx++
			return NewToken("LESS_EQUAL", "<=")
		}
		return NewToken("LESS", "<")
	case '>':
		if len(s.fileContents) > s.currentIdx && rune(s.fileContents[s.currentIdx]) == '=' {
			s.currentIdx++
			return NewToken("GREATER_EQUAL", ">=")
		}
		return NewToken("GREATER", ">")

	default:
		fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", s.lines, current)
		s.ExitCode = 65
	}
	return nil

}

func (s Scanner) Scan() {
	for s.currentIdx < len(s.fileContents) {
		token := NextToken(&s)
		if token != nil {
			fmt.Printf("%s %v null\n", token.tokenType, token.lexeme)
		}
	}
}
