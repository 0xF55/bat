package main

import (
	"fmt"
	"log"
	"os"
)

type Lexer struct {
	input        string
	leninput     int
	position     int
	peekposition int
	currentchar  byte
	tokens       []Token
}

func NewLexer(source_file string) *Lexer {

	// read file
	data, err := os.ReadFile(source_file)
	if err != nil {
		log.Fatal("Failed to read bat script file")
	}

	lexer := &Lexer{}
	lexer.input = string(data)
	lexer.leninput = len(lexer.input)
	lexer.init()
	return lexer
}

// initialize the lexer
func (l *Lexer) init() {
	if l.leninput == 0 {
		l.currentchar = 0
		return
	}
	l.position = 0
	l.peekposition = 1
	l.currentchar = l.input[l.position]
}

// get current char
func (l *Lexer) currch() byte {
	return l.currentchar
}

// is delimiter
func (l *Lexer) isdelem(char byte) bool {
	return char == ' ' || char == '=' || char == '\n' || char == '\t' || char == '\r' || char == ',' || char == '(' || char == ')'
}

// read one char
func (l *Lexer) readch() byte {

	l.position++
	if l.position >= l.leninput {
		l.currentchar = 0
		return 0
	}
	l.peekposition = l.position + 1
	l.currentchar = l.input[l.position]

	return l.currentchar
}

func (l *Lexer) skipWhitespace() {
	for l.currentchar == ' ' || l.currentchar == '\t' || l.currentchar == '\n' || l.currentchar == '\r' {

		if l.readch() == 0 {
			break
		}
	}
}

// read literal
func (l *Lexer) readliteral() string {
	char := l.currch()
	pos := l.position
	for char != 0 && !l.isdelem(char) {
		char = l.readch()
		if char == '$' {
			break
		}
	}
	ret := l.input[pos:l.position]
	if char != '=' {
		l.position--
		l.currentchar = l.input[l.position]
	}
	return ret
}

// lex token
func (l *Lexer) NextToken() bool {

	l.skipWhitespace()

	if l.currch() == 0 {
		return false
	}

	var tokenType TokenType
	token := Token{}
	char := l.currch()

	switch char {
	case '$':
		tokenType = VAR_NEEDED
		token.VALUE = l.readliteral()
	case '@':
		tokenType = KEYWORD_OUT
	case '(':
		tokenType = ROUND_START
	case ')':
		tokenType = ROUND_END
	case ',':
		tokenType = COMMA
		token.VALUE = ","
	default:
		word := l.readliteral()
		switch word {
		case "for":
			tokenType = KEYWORD_FOR
		case "end":
			tokenType = KEYWORD_END
		default:
			len_tokens := len(l.tokens)
			if l.currch() == '=' {
				tokenType = IDENTIFIER
			} else if len_tokens > 0 && l.tokens[len_tokens-1].Type == KEYWORD_FOR {
				tokenType = IDENTIFIER
			} else {
				tokenType = LITERAL
			}
			token.VALUE = word
		}
	}

	l.readch()
	token.Type = tokenType
	l.tokens = append(l.tokens, token)
	return true
}

// token by token
func (l *Lexer) Lex() {

	for l.NextToken() {
	}

}

// for debugging
func (l *Lexer) Print() {
	for _, token := range l.tokens {
		fmt.Printf("%d %s\n", token.Type, token.VALUE)
	}
}
