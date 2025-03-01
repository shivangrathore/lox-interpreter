package scanner

import (
	"fmt"
	"os"
)

// Token Type
type TokenType string

const (
	// Single-character tokens
	LEFT_PAREN  TokenType = "LEFT_PAREN"
	RIGHT_PAREN TokenType = "RIGHT_PAREN"
	LEFT_BRACE  TokenType = "LEFT_BRACE"
	RIGHT_BRACE TokenType = "RIGHT_BRACE"
	COMMA       TokenType = "COMMA"
	DOT         TokenType = "DOT"
	MINUS       TokenType = "MINUS"
	PLUS        TokenType = "PLUS"
	SEMICOLON   TokenType = "SEMICOLON"
	STAR        TokenType = "STAR"
	// One or two character tokens
	BANG          TokenType = "BANG"
	BANG_EQUAL    TokenType = "BANG_EQUAL"
	EQUAL         TokenType = "EQUAL"
	EQUAL_EQUAL   TokenType = "EQUAL_EQUAL"
	GREATER       TokenType = "GREATER"
	GREATER_EQUAL TokenType = "GREATER_EQUAL"
	LESS          TokenType = "LESS"
	LESS_EQUAL    TokenType = "LESS_EQUAL"
	SLASH         TokenType = "SLASH"
)

type Scanner struct {
	fileContents []byte
	currentIdx   int
	exitCode     int
	lines        int
}

type Token struct {
	tokenType TokenType
	lexeme    interface{}
}

func (s *Scanner) ExitCode() int {
	return s.exitCode
}

func NewScanner(fileContents []byte) *Scanner {
	return &Scanner{fileContents: fileContents, currentIdx: 0, exitCode: 0, lines: 1}
}

func NewToken(tokenType TokenType, lexeme interface{}) *Token {
	return &Token{tokenType: tokenType, lexeme: lexeme}
}

func NextToken(s *Scanner) (*Token, error) {
	current := rune(s.fileContents[s.currentIdx])
	s.currentIdx++
	switch current {
	case '=':
		if len(s.fileContents) > s.currentIdx && rune(s.fileContents[s.currentIdx]) == '=' {
			s.currentIdx++
			return NewToken(EQUAL_EQUAL, "=="), nil
		}
		return NewToken(EQUAL, "="), nil
	case '\n':
		s.lines++
		return nil, nil
	case '!':
		if len(s.fileContents) > s.currentIdx && rune(s.fileContents[s.currentIdx]) == '=' {
			s.currentIdx++
			return NewToken(BANG_EQUAL, "!="), nil
		}
		return NewToken(BANG, "!"), nil
	case '/':
		if len(s.fileContents) > s.currentIdx && rune(s.fileContents[s.currentIdx]) == '/' {
			s.currentIdx++
			for s.currentIdx < len(s.fileContents) && rune(s.fileContents[s.currentIdx]) != '\n' {
				s.currentIdx++
			}
			return nil, nil
		}

		return NewToken(SLASH, "/"), nil

	case '.':
		return NewToken(DOT, "."), nil
	case ',':
		return NewToken(COMMA, ","), nil
	case '+':
		return NewToken(PLUS, "+"), nil
	case '-':
		return NewToken(MINUS, "-"), nil
	case ';':
		return NewToken(SEMICOLON, ";"), nil
	case '(':
		return NewToken(LEFT_PAREN, "("), nil
	case ')':
		return NewToken(RIGHT_PAREN, ")"), nil
	case '{':
		return NewToken(LEFT_BRACE, "{"), nil
	case '}':
		return NewToken(RIGHT_BRACE, "}"), nil
	case '*':
		return NewToken(STAR, "*"), nil
	case '<':
		if len(s.fileContents) > s.currentIdx && rune(s.fileContents[s.currentIdx]) == '=' {
			s.currentIdx++
			return NewToken(LESS_EQUAL, "<="), nil
		}
		return NewToken(LESS, "<"), nil
	case '>':
		if len(s.fileContents) > s.currentIdx && rune(s.fileContents[s.currentIdx]) == '=' {
			s.currentIdx++
			return NewToken(GREATER_EQUAL, ">="), nil
		}
		return NewToken(GREATER, ">"), nil

	default:
		return nil, fmt.Errorf("[line %d] Error: Unexpected character: %c\n", s.lines, current)
	}
}

func (s *Scanner) Scan() {
	for s.currentIdx < len(s.fileContents) {
		token, err := NextToken(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err.Error())
			s.exitCode = 65
		} else if token != nil {
			fmt.Printf("%s %v null\n", token.tokenType, token.lexeme)
		}
	}
}
