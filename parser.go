package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type Parser struct {
	Tokens       []Token
	TokensLength int
	CurrentToken Token
	Position     int
	PeekPosition int
}

/* Maps For Storing Vars,Lists */
var Variables map[string]string
var Lists map[string]*[]string
var Writer *oWriter

func NewParser(tokens []Token) *Parser {

	parser := &Parser{}
	parser.Tokens = tokens
	parser.TokensLength = len(parser.Tokens)
	parser.Position = 0
	parser.PeekPosition = parser.Position + 1
	parser.CurrentToken = parser.Tokens[parser.Position]

	Variables = make(map[string]string)
	Lists = make(map[string]*[]string)
	// default output value
	Writer = NewWriter()

	return parser

}

func (p *Parser) NextToken() Token {

	p.Position++
	if p.Position >= p.TokensLength {
		t := Token{}
		p.CurrentToken = t
		return t
	}
	p.PeekPosition = p.PeekPosition + 1
	p.CurrentToken = p.Tokens[p.Position]
	return p.CurrentToken

}

func (p *Parser) TokenIs(args ...TokenType) bool {

	for _, arg := range args {
		if p.CurrentToken.Type == arg {
			return true
		}
	}
	return false

}

func (p *Parser) PeekToken() Token {
	return p.Tokens[p.PeekPosition]
}

func (p *Parser) PeekIs(Type TokenType) bool {

	return p.PeekToken().Type == Type

}

func (p *Parser) BackToken() Token {
	return p.Tokens[p.PeekPosition]
}

func (p *Parser) BackIs(Type TokenType) bool {

	return p.BackToken().Type == Type

}

func (p *Parser) TokenAt(index int) Token {
	if index >= p.TokensLength {
		return Token{}
	}
	return p.Tokens[index]
}

func (p *Parser) Goto(index int) {
	if index >= p.TokensLength {
		return
	}
	p.Position = index
	p.PeekPosition = p.Position + 1
	p.CurrentToken = p.Tokens[p.Position]

}

func (p *Parser) ParseAssignment() {

	ident := p.CurrentToken.VALUE.(string) // var | list name
	p.NextToken()
	pos := p.Position
	if p.CurrentToken.Type == ROUND_START {
		for p.TokenIs(LITERAL, VAR_NEEDED, COMMA, ROUND_START, ROUND_END) {
			p.NextToken()
			if p.CurrentToken.Type == 0 {
				break
			}

		}

		if (p.Position - pos) < 2 {
			log.Fatal("Invalid List Expression")
			return
		}
		tokens := p.Tokens[pos:p.Position]
		// is a list
		if tokens[0].Type == ROUND_START {
			// Round Bracket not closed
			if tokens[len(tokens)-1].Type != ROUND_END {
				log.Fatal("Invalid List Expression")
				return
			}
			new_list := make([]string, 0)
			// add items to the list
			for _, token := range tokens {
				t := token.Type
				if t == LITERAL || t == VAR_NEEDED {
					value := token.VALUE.(string)
					if strings.TrimSpace(value) == "" {
						continue
					}
					new_list = append(new_list, Eval(value))
				}
			} // end loop

			// add list
			Lists[ident] = &new_list
			return

		} // end if
	}

	value := ""
	for p.TokenIs(LITERAL, VAR_NEEDED, COMMA) {
		value += p.CurrentToken.VALUE.(string)
		p.NextToken()
	}

	Variables[ident] = Eval(value)

}

func (p *Parser) ParseOut() {

	value := ""
	p.NextToken()
	for p.TokenIs(LITERAL, VAR_NEEDED, COMMA) {

		value2 := p.CurrentToken.VALUE.(string)
		if value2 == "*" {
			for k, v := range Variables {
				if k != "out" {
					Writer.Write(v)
				}
			}
			return
		}
		value += value2
		p.NextToken()
	}

	Writer.Write(Eval(value))

}

func (p *Parser) ParseLoop() {

	p.NextToken()
	ctx := LoopCTX{}
	// like for i=$list or 1..5 or file
	if p.CurrentToken.Type == IDENTIFIER {
		// declare the iterator
		val := p.CurrentToken.VALUE.(string)
		ctx.Iterator = val
		Variables[val] = ""
		p.NextToken()
		switch p.CurrentToken.Type {
		// should be range -> 1..10 or file
		case LITERAL:
			name := p.CurrentToken.VALUE.(string)
			Type, res := GetList(name)
			if res == nil {
				// it's a file ?
				file, err := os.Open(name)
				if err != nil {
					log.Fatalf("Unknown: %s", name)
					return
				}
				fileLoop := &FileLoop{}
				fileLoop.Scanner = bufio.NewScanner(file)
				Type = File
				res = fileLoop
			}
			ctx.Type = Type
			ctx.Expr = res
		// for i=$list
		case VAR_NEEDED:
			name := p.CurrentToken.VALUE.(string)
			name = strings.Trim(name, "$")
			Type, res := GetList(name)
			ctx.Type = Type
			ctx.Expr = res

		}
	} else {
		log.Fatal("Invalid Loop Expression")
		return
	}

	loop_start := p.Position
	ret_pos := 0
	for i := p.Position; i < p.TokensLength; i++ {
		if p.TokenAt(i).Type == KEYWORD_END {
			ret_pos = i
			for ctx.Next() {
				p.Parse(loop_start, i)
			}
			break
		}

	}

	if ret_pos == 0 {
		log.Fatalf("Missing end block for the loop")
		return
	}

	p.Goto(ret_pos)
	p.NextToken()

}

// start from position to position (main = 0,end)
func (p *Parser) Parse(start, end int) {
	if p.TokensLength <= 1 {
		log.Fatal("Too short script :(")
		return
	}
	p.Position = start
	p.PeekPosition = p.Position + 1
	p.CurrentToken = p.Tokens[p.Position]

	for p.CurrentToken.Type != 0 {
		switch p.CurrentToken.Type {
		case IDENTIFIER:
			p.ParseAssignment()
		case KEYWORD_OUT:
			p.ParseOut()
		case KEYWORD_FOR:
			p.ParseLoop()
		case KEYWORD_END, COMMA:
			p.NextToken()
		default:
			p.NextToken()
		}
	}

}
