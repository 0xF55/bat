/*
   Copyright [2025] [0xf55]

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.


*/

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
)

type Parser struct {
	Tokens       []Token
	TokensLength int
	CurrentToken Token
	Position     int
	PeekPosition int
	Writer       *bWriter
	Lists        map[string]*[]string
	Variables    map[string]string
}

/* Main Writer */
var Writer *oWriter

/* Cache for reversed strings in loops */
var ReversedCache map[string]string

func NewParser(tokens []Token) *Parser {

	parser := &Parser{}
	parser.Tokens = tokens
	parser.TokensLength = len(parser.Tokens)
	parser.Position = 0
	parser.PeekPosition = parser.Position + 1
	parser.CurrentToken = parser.Tokens[parser.Position]
	parser.Writer = &bWriter{}
	parser.Lists = make(map[string]*[]string)
	parser.Variables = make(map[string]string)

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
					new_list = append(new_list, p.Eval(value))
				}
			} // end loop

			// add list
			p.Lists[ident] = &new_list
			return

		} // end if
	}

	value := strings.Builder{}
	for p.TokenIs(LITERAL, VAR_NEEDED, COMMA) {
		value.WriteString(p.CurrentToken.VALUE.(string))
		p.NextToken()
	}

	p.Variables[ident] = p.Eval(value.String())

}

func (p *Parser) ParseOut() {

	value := ""
	p.NextToken()
	for p.TokenIs(LITERAL, VAR_NEEDED, COMMA) {

		value2 := p.CurrentToken.VALUE.(string)
		if value2 == "*" {
			for _, v := range p.Variables {
				p.Writer.Write(v)
			}
			return
		}
		value += value2
		p.NextToken()
	}

	p.Writer.Write(p.Eval(value))

}

func (p *Parser) ParseLoop(in_loop bool) {

	p.NextToken()
	ctx := LoopCTX{}
	// like for i=$list or 1..5 or file
	if p.CurrentToken.Type == IDENTIFIER {
		// declare the iterator
		val := p.CurrentToken.VALUE.(string)
		ctx.Iterator = val
		p.Variables[val] = ""
		p.NextToken()
		switch p.CurrentToken.Type {
		// should be range -> 1..10 or file
		case LITERAL:
			name := p.CurrentToken.VALUE.(string)
			Type, res := p.GetList(name)
			if res == nil {
				// it's a file ?
				file, err := os.Open(strings.Trim(name, "%"))
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
			Type, res := p.GetList(name)
			ctx.Type = Type
			ctx.Expr = res

		}
	} else {
		log.Fatal("Invalid Loop Expression")
	}

	loop_start := p.Position
	ret_pos := 0
	loops := 0

	for i := p.Position; i < p.TokensLength; i++ {
		/* Fixed Error in loops logic in 1.2.0*/
		if p.TokenAt(i).Type == KEYWORD_FOR {
			loops++
		}
		if p.TokenAt(i).Type == KEYWORD_END {
			if loops > 0 {
				loops--
				continue
			}
			ret_pos = i
			for ctx.Next(p) {
				p.Parse(loop_start, i, true)
			}
			break
		}

	}

	if ret_pos == 0 {
		log.Fatalf("Missing end keyword for the loop")
		return
	}

	if !in_loop {
		p.Goto(ret_pos)
		p.NextToken()
	}

}

func (p *Parser) ParseDirective() {

	p.NextToken()

	switch p.CurrentToken.Type {
	case LITERAL, VAR_NEEDED:
		ExecuteDirective(p.Eval(p.CurrentToken.VALUE.(string)))
	default:
		return
	}

	p.NextToken()

}

// start from position to position (main = 0,end)
func (p *Parser) Parse(start, end int, in_loop bool) {
	if p.TokensLength <= 1 {
		log.Fatal("Too short script :(")
		return
	}

	p.Goto(start)

	for p.CurrentToken.Type != 0 && p.Position < end {
		switch p.CurrentToken.Type {
		case IDENTIFIER:
			p.ParseAssignment()
		case KEYWORD_OUT:
			p.ParseOut()
		case KEYWORD_FOR:
			p.ParseLoop(in_loop)
		case DIRECTIVE:
			p.ParseDirective()
		case KEYWORD_END, COMMA:
			p.NextToken()
		default:
			p.NextToken()
		}
	}

	if !Quiet {
		fmt.Printf(color.GreenString("\r> Lines:        %d", Writer.lines))
	}
	Writer.Write(p.Writer)

	p.Writer.buffer.Reset()
	p.Writer.lines = 0

}
