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
	"log"
	"sync"
)

func Run(batfile *string, wg *sync.WaitGroup) {

	defer wg.Done()
	lexer := NewLexer(*batfile)
	lexer.Lex()
	if len(lexer.tokens) <= 0 {
		return
	}
	parser := NewParser(lexer.tokens)
	parser.Parse(0, parser.TokensLength, false)

}

func RunAll() {

	wg := &sync.WaitGroup{}

	if len(BatFiles) < 1 {
		log.Fatal("No BatFiles Passed!")
	}

	// run all
	for _, bt := range BatFiles {
		wg.Add(1)
		go Run(&bt, wg)
	}

	// wait for all goroutins
	wg.Wait()
}
