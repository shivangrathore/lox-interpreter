package scanner

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/app/utils"
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
	STRING        TokenType = "STRING"
	NUMBER        TokenType = "NUMBER"
)

type Scanner struct {
	fileContents []byte
	currentIdx   int
	exitCode     int
	lines        int
}

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   *string
}

func (s *Scanner) ExitCode() int {
	return s.exitCode
}

func NewScanner(fileContents []byte) *Scanner {
	return &Scanner{fileContents: fileContents, currentIdx: 0, exitCode: 0, lines: 1}
}

func NewToken(tokenType TokenType, lexeme string, literal *string) *Token {
	return &Token{tokenType: tokenType, lexeme: lexeme, literal: literal}
}

func NextToken(s *Scanner) (*Token, error) {
	current := rune(s.fileContents[s.currentIdx])
	s.currentIdx++
	switch current {
	case '=':
		if len(s.fileContents) > s.currentIdx && rune(s.fileContents[s.currentIdx]) == '=' {
			s.currentIdx++
			return NewToken(EQUAL_EQUAL, "==", nil), nil
		}
		return NewToken(EQUAL, "=", nil), nil

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		lexeme := string(current)
		isDec := false
		for len(s.fileContents) > s.currentIdx {
			if rune(s.fileContents[s.currentIdx]) == '.' {
				if len(s.fileContents) > s.currentIdx+1 && utils.IsDigit(rune(s.fileContents[s.currentIdx+1])) {
					isDec = true
				} else {
					break
				}
			} else if !utils.IsDigit(rune(s.fileContents[s.currentIdx])) {
				break
			}
			lexeme += string(s.fileContents[s.currentIdx])
			s.currentIdx++
		}
		literal := lexeme
		if !isDec {
			literal += ".0"
		}
		return NewToken(NUMBER, lexeme, &literal), nil

	case '\n':
		s.lines++
		return nil, nil

	case '!':
		if len(s.fileContents) > s.currentIdx && rune(s.fileContents[s.currentIdx]) == '=' {
			s.currentIdx++
			return NewToken(BANG_EQUAL, "!=", nil), nil
		}
		return NewToken(BANG, "!", nil), nil

	case '/':
		if len(s.fileContents) > s.currentIdx && rune(s.fileContents[s.currentIdx]) == '/' {
			s.currentIdx++
			for s.currentIdx < len(s.fileContents) && rune(s.fileContents[s.currentIdx]) != '\n' {
				s.currentIdx++
			}
			return nil, nil
		}
		return NewToken(SLASH, "/", nil), nil

	case '\t', ' ':
		return nil, nil

	case '"':
		currStr := ""
		for s.currentIdx < len(s.fileContents) {
			char := rune(s.fileContents[s.currentIdx])
			s.currentIdx++
			if char == '"' {
				return NewToken(STRING, fmt.Sprintf("\"%s\"", currStr), &currStr), nil
			} else if char == '\n' {
				return nil, fmt.Errorf("[line %d] Error: Unterminated string.\n", s.lines)
			}
			currStr += string(char)
		}
		return nil, fmt.Errorf("[line %d] Error: Unterminated string.\n", s.lines)

	case '.':
		return NewToken(DOT, ".", nil), nil
	case ',':
		return NewToken(COMMA, ",", nil), nil
	case '+':
		return NewToken(PLUS, "+", nil), nil
	case '-':
		return NewToken(MINUS, "-", nil), nil
	case ';':
		return NewToken(SEMICOLON, ";", nil), nil
	case '(':
		return NewToken(LEFT_PAREN, "(", nil), nil
	case ')':
		return NewToken(RIGHT_PAREN, ")", nil), nil
	case '{':
		return NewToken(LEFT_BRACE, "{", nil), nil
	case '}':
		return NewToken(RIGHT_BRACE, "}", nil), nil
	case '*':
		return NewToken(STAR, "*", nil), nil
	case '<':
		if len(s.fileContents) > s.currentIdx && rune(s.fileContents[s.currentIdx]) == '=' {
			s.currentIdx++
			return NewToken(LESS_EQUAL, "<=", nil), nil
		}
		return NewToken(LESS, "<", nil), nil
	case '>':
		if len(s.fileContents) > s.currentIdx && rune(s.fileContents[s.currentIdx]) == '=' {
			s.currentIdx++
			return NewToken(GREATER_EQUAL, ">=", nil), nil
		}
		return NewToken(GREATER, ">", nil), nil

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
			literal := "null"
			if token.literal != nil {
				literal = *token.literal
			}
			fmt.Printf("%s %s %s\n", token.tokenType, token.lexeme, literal)
		}
	}
}
