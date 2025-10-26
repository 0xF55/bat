package main

type TokenType int8

const (
	IDENTIFIER  TokenType = iota + 1
	VAR_NEEDED            // $var
	KEYWORD_OUT           // @
	KEYWORD_FOR           // for
	KEYWORD_END           // end
	ROUND_START           // (
	ROUND_END             // )
	LITERAL               // value
	COMMA                 // ,
)

type Token struct {
	Type  TokenType
	VALUE any
}
