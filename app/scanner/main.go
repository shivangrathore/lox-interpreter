package scanner

import (
	"fmt"
	"github.com/codecrafters-io/interpreter-starter-go/app/utils"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
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
	EOF         TokenType = "EOF"

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
	CLASS         TokenType = "CLASS"
	AND           TokenType = "AND"
	OR            TokenType = "OR"
	IF            TokenType = "IF"
	ELSE          TokenType = "ELSE"
	FUN           TokenType = "FUN"
	FOR           TokenType = "FOR"
	NIL           TokenType = "NIL"
	TRUE          TokenType = "TRUE"
	FALSE         TokenType = "FALSE"
	PRINT         TokenType = "PRINT"
	RETURN        TokenType = "RETURN"
	SUPER         TokenType = "SUPER"
	THIS          TokenType = "THIS"
	VAR           TokenType = "VAR"
	WHILE         TokenType = "WHILE"
)

func singleCharacters(c rune) TokenType {
	var chars map[rune]TokenType = map[rune]TokenType{
		'(': LEFT_PAREN,
		')': RIGHT_PAREN,
		'{': LEFT_BRACE,
		'}': RIGHT_BRACE,
		',': COMMA,
		'.': DOT,
		'-': MINUS,
		'+': PLUS,
		';': SEMICOLON,
		'*': STAR,
	}
	return chars[c]
}

func matchOperators(s string) TokenType {
	var operators map[string]TokenType = map[string]TokenType{
		"!":  BANG,
		"!=": BANG_EQUAL,
		"=":  EQUAL,
		"==": EQUAL_EQUAL,
		">":  GREATER,
		">=": GREATER_EQUAL,
		"<":  LESS,
		"<=": LESS_EQUAL,
		"/":  SLASH,
	}
	return operators[s]
}

func matchKeywords(s string) TokenType {
	var keywords map[string]TokenType = map[string]TokenType{
		"class":  CLASS,
		"and":    AND,
		"or":     OR,
		"if":     IF,
		"else":   ELSE,
		"fun":    FUN,
		"for":    FOR,
		"nil":    NIL,
		"true":   TRUE,
		"false":  FALSE,
		"print":  PRINT,
		"return": RETURN,
		"super":  SUPER,
		"this":   THIS,
		"var":    VAR,
		"while":  WHILE,
	}

	tt, exists := keywords[s]
	if exists {
		return tt
	}
	return IDENTIFIER
}

type Scanner struct {
	source   []byte
	current  int
	exitCode int
	lines    int
	start    int
	tokens   []Token
}

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   interface{}
	line      int
}

func (t Token) String() string {
	literal := ""
	if t.literal == nil {
		literal = "null"
	} else if t.tokenType == NUMBER {
		parts := strings.Split(fmt.Sprintf("%f", t.literal), ".")
		integerPart := parts[0]
		decimalPart := parts[1]

		integerPart = strings.TrimLeft(integerPart, "0")
		decimalPart = strings.TrimRight(decimalPart, "0")
		if len(decimalPart) == 0 {
			decimalPart = "0"
		}
		if len(integerPart) == 0 {
			integerPart = "0"
		}
		literal = fmt.Sprintf("%s.%s", integerPart, decimalPart)
	} else {
		literal = fmt.Sprintf("%v", t.literal)
	}
	return fmt.Sprintf("%s %s %s", t.tokenType, t.lexeme, literal)
}

func (s *Scanner) ExitCode() int {
	return s.exitCode
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() (rune, int) {
	if s.isAtEnd() {
		return 0, 0
	}
	r, size := utf8.DecodeRune(s.source[s.current:])
	s.current += size
	return r, size
}

func (s *Scanner) peak() rune {
	if s.isAtEnd() {
		return 0
	}
	r, _ := utf8.DecodeRune(s.source[s.current:])
	return r
}

func (s *Scanner) addToken(tokenType TokenType, lexeme string, literal interface{}) {
	s.tokens = append(s.tokens, Token{
		tokenType: tokenType,
		lexeme:    lexeme,
		literal:   literal,
	})
}

func (s *Scanner) operators(r rune) {
	op := string(r)
	if s.match('=') {
		op += "="
	}
	s.addToken(matchOperators(op), op, nil)
}

func (s *Scanner) setExitCode(code int) {
	if s.exitCode == 0 {
		s.exitCode = code
	}
}

func (s *Scanner) scanString() {
	str := ""
	for !s.isAtEnd() {
		if s.peak() == '"' {
			s.advance()
			s.addToken(STRING, fmt.Sprintf("\"%s\"", str), str)
			return
		} else if rune(s.source[s.current]) == '\n' {
			s.setExitCode(65)
			fmt.Fprintf(os.Stderr, "[line %d] Error: Unterminated string.\n", s.lines)
			return
		} else {
			char, _ := s.advance()
			str += string(char)
		}
	}
	s.setExitCode(65)
	fmt.Fprintf(os.Stderr, "[line %d] Error: Unterminated string.\n", s.lines)

}

func (s *Scanner) scanNumber() {
	isDecimal := false
	for !s.isAtEnd() {
		if utils.IsDigit(s.peak()) {
			s.advance()
		} else if s.peak() == '.' && !isDecimal && utils.IsDigit(rune(s.source[s.current+1])) {
			isDecimal = true
			s.advance()
		} else {
			break
		}
	}

	lexeme := string(s.source[s.start:s.current])
	value, err := strconv.ParseFloat(lexeme, 64)
	if err != nil {
		s.setExitCode(65)
		fmt.Fprintf(os.Stderr, "[line %d] Error: Invalid number: %s\n", s.lines, lexeme)
		return
	}
	s.addToken(NUMBER, lexeme, value)
}
func (s *Scanner) scanIdentifier() {
	for !s.isAtEnd() && (utils.IsAlpha(s.peak()) || utils.IsDigit(s.peak())) {
		s.advance()
	}
	lexeme := string(s.source[s.start:s.current])
	s.addToken(matchKeywords(lexeme), lexeme, nil)
}

func NewScanner(contents []byte) *Scanner {
	return &Scanner{source: contents, current: 0, exitCode: 0, lines: 1, start: 0, tokens: []Token{}}
}

func (s *Scanner) scanToken() {
	s.start = s.current
	r, _ := s.advance()
	switch r {
	case '<', '>', '=', '!':
		s.operators(r)
	case '\n':
		s.lines++
	case ' ', '\r', '\t':
	case '"':
		s.scanString()
	case '(', ')', '{', '}', ',', '.', '-', '+', ';', '*':
		s.addToken(singleCharacters(r), string(r), nil)
	case '/':
		if s.match('/') {
			for s.peak() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH, "/", nil)
		}
	default:
		if utils.IsDigit(r) {
			s.scanNumber()
			return
		}
		if utils.IsAlpha(r) {
			s.scanIdentifier()
			return
		}
		s.setExitCode(65)
		fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", s.lines, r)
	}
}

func (s *Scanner) match(c rune) bool {
	if s.current >= len(s.source) {
		return false
	}
	if rune(s.source[s.current]) != c {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) Scan() {
	for !s.isAtEnd() {
		s.scanToken()
	}
	s.addToken(EOF, "", nil)
	for _, token := range s.tokens {
		fmt.Println(token)
	}
}
