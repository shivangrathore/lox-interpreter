package scanner

import (
	"fmt"
	"os"
	"strings"

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
	IDENTIFIER    TokenType = "IDENTIFIER"
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
	if current == '=' {
		if len(s.fileContents) > s.currentIdx && rune(s.fileContents[s.currentIdx]) == '=' {
			s.currentIdx++
			return NewToken(EQUAL_EQUAL, "==", nil), nil
		}
		return NewToken(EQUAL, "=", nil), nil
	} else if current == '\n' {
		s.lines++
		return nil, nil
	} else if current == '!' {
		if len(s.fileContents) > s.currentIdx && rune(s.fileContents[s.currentIdx]) == '=' {
			s.currentIdx++
			return NewToken(BANG_EQUAL, "!=", nil), nil
		}
		return NewToken(BANG, "!", nil), nil
	} else if current == '/' {
		if len(s.fileContents) > s.currentIdx && rune(s.fileContents[s.currentIdx]) == '/' {
			s.currentIdx++
			for s.currentIdx < len(s.fileContents) && rune(s.fileContents[s.currentIdx]) != '\n' {
				s.currentIdx++
			}
			return nil, nil
		}
		return NewToken(SLASH, "/", nil), nil
	} else if current == '\t' || current == ' ' {
		return nil, nil
	} else if current == '"' {
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
	} else if current == '.' {
		return NewToken(DOT, ".", nil), nil
	} else if current == ',' {
		return NewToken(COMMA, ",", nil), nil
	} else if current == '+' {
		return NewToken(PLUS, "+", nil), nil
	} else if current == '-' {
		return NewToken(MINUS, "-", nil), nil
	} else if current == ';' {
		return NewToken(SEMICOLON, ";", nil), nil
	} else if current == '(' {
		return NewToken(LEFT_PAREN, "(", nil), nil
	} else if current == ')' {
		return NewToken(RIGHT_PAREN, ")", nil), nil
	} else if current == '{' {
		return NewToken(LEFT_BRACE, "{", nil), nil
	} else if current == '}' {
		return NewToken(RIGHT_BRACE, "}", nil), nil
	} else if current == '*' {
		return NewToken(STAR, "*", nil), nil
	} else if current == '<' {
		if len(s.fileContents) > s.currentIdx && rune(s.fileContents[s.currentIdx]) == '=' {
			s.currentIdx++
			return NewToken(LESS_EQUAL, "<=", nil), nil
		}
		return NewToken(LESS, "<", nil), nil
	} else if current == '>' {
		if len(s.fileContents) > s.currentIdx && rune(s.fileContents[s.currentIdx]) == '=' {
			s.currentIdx++
			return NewToken(GREATER_EQUAL, ">=", nil), nil
		}
		return NewToken(GREATER, ">", nil), nil
	} else if utils.IsDigit(current) {
		lexeme := string(current)
		isDec := false
		for len(s.fileContents) > s.currentIdx {
			if rune(s.fileContents[s.currentIdx]) == '.' {
				if len(s.fileContents) > s.currentIdx+1 && utils.IsDigit(rune(s.fileContents[s.currentIdx+1])) && !isDec {
					isDec = true
				} else {
					break
				}
			} else if !utils.IsDigit(rune(s.fileContents[s.currentIdx])) {
				break
			}
			digChar := rune(s.fileContents[s.currentIdx])
			lexeme += string(digChar)
			s.currentIdx++
		}
		// TODO: Might need later

		// integer := 0.0
		// fraction := 0.0
		// fractionDigits := 1.0
		// isFraction := false
		// for _, c := range lexeme {
		// 	if c == '.' {
		// 		isFraction = true
		// 		continue
		// 	}
		// 	dig := float64(9 - (int('9') - int(c)))
		// 	if isFraction {
		// 		fractionDigits *= 10
		// 		fraction = fraction*10 + dig
		// 	} else {
		// 		integer = integer*10 + dig
		// 	}
		// }

		// literal := integer + (fraction / fractionDigits)
		// literalString := ""
		// {
		// 	str := fmt.Sprintf("%f", literal)
		// 	parts := strings.Split(str, ".")
		// 	if len(parts) == 1 {
		// 		literalString = parts[0] + ".0"
		// 		return NewToken(NUMBER, lexeme, &literalString), nil
		// 	}
		//
		// 	integerPart := parts[0]
		// 	fractionPart := parts[1]
		//
		// 	fractionPart = strings.TrimRight(fractionPart, "0")
		//
		// 	if len(fractionPart) == 0 {
		// 		literalString = integerPart + ".0"
		// 		return NewToken(NUMBER, lexeme, &literalString), nil
		// 	} else {
		// 		literalString = integerPart + "." + fractionPart
		// 		return NewToken(NUMBER, lexeme, &literalString), nil
		// 	}
		// }
		literalString := ""
		parts := strings.Split(lexeme, ".")
		if len(parts) == 1 {
			integerPart := strings.TrimLeft(parts[0], "0")
			if len(integerPart) == 0 {
				integerPart = "0"
			}
			literalString = integerPart + ".0"
			return NewToken(NUMBER, lexeme, &literalString), nil
		}
		intergerPart := strings.TrimLeft(parts[0], "0")
		fractionPart := strings.TrimRight(parts[1], "0")

		if len(intergerPart) == 0 {
			intergerPart = "0"
		}
		if len(fractionPart) == 0 {
			fractionPart = "0"
		}
		literalString = intergerPart + "." + fractionPart
		return NewToken(NUMBER, lexeme, &literalString), nil
	} else if utils.IsAlpha(current) || current == '_' {
		lexeme := string(current)
		for s.currentIdx < len(s.fileContents) {
			current = rune(s.fileContents[s.currentIdx])
			if utils.IsAlpha(current) || utils.IsDigit(current) || current == '_' {
				lexeme += string(current)
				s.currentIdx++
			} else {
				break
			}
		}
		return NewToken(IDENTIFIER, lexeme, nil), nil
	} else {
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
