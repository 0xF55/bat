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
	"os"
	"regexp"
	"strconv"
	"strings"
)

func (p *Parser) GetVar(name string) string {

	name = strings.Trim(name, "$")
	switch name {
	default:
		lower := strings.ToLower(name)
		ret := ""
		// generators -> rndn,rndc,rnds:
		if strings.HasPrefix(lower, "rnd") && len(lower) > 4 {
			switch lower[3] {
			case 'n':
				ret = RandNum(lower)
			case 'c':
				ret = RandChar(lower)
			case 's':
				ret = RandSpecial(lower)
			default:
				return "" // soon
			}
			return ret
		}
		f := name[0]
		switch f {
		case '+', '-', '^', '~', '!':
			name = strings.Trim(name, string(f))
		}
		ret = string(p.Variables[name])
		switch f {
		case '+':
			ret = strings.ToUpper(ret)
		case '-':
			ret = strings.ToLower(ret)
		case '^':
			ret = Capitaltize(ret)
		case '~':
			ret = ZigZag(ret)
		case '!':
			ret = Reverse(ret)
		}
		if strings.TrimSpace(ret) == "" {
			return os.Getenv(name)
		}
		return ret
	}

}

func UnEscape(text string) string {

	replacer := strings.NewReplacer("`SP`", " ", "`COM`", ",", "`RSA`", "(", "`RCL`", ")", "`DS`", "$", "`EQ`", "=", "`AT`", "@")

	return replacer.Replace(text)

}

func (p *Parser) EvalExpression(text string) string {
	varegex := regexp.MustCompile(`\$\s*[\+\-\~\^\!\*]?\w+(?::\d+)?`) // $?var || $gen:n
	comma_sp := strings.Split(text, ",")
	var ret strings.Builder

	for _, expr := range comma_sp {
		expr = varegex.ReplaceAllStringFunc(expr, func(match string) string {
			return p.GetVar(match)
		})
		ret.WriteString(expr)
	}
	return ret.String()
}

func (p *Parser) Eval(text string) string {
	ret := UnEscape(p.EvalExpression(text))
	return ret
}

func (p *Parser) GetList(name string) (int8, any) {

	// is a file loop
	if name[0] == '%' {
		return File, nil
	}

	count := strings.Count(name, ".")
	sp := ".."

	if count > 0 {
		if count == 1 {
			sp = "."
		}
		Splitted := strings.Split(name, sp)
		if len(Splitted) != 2 || Splitted[1] == "" {
			log.Fatalf("Invalid Range Expression: %s", name)
		}

		n1, err1 := strconv.Atoi(Splitted[0])
		if err1 != nil {
			log.Fatalf("Numbers Only Allowed! Got: %s", Splitted[0])
		}
		n2, err2 := strconv.Atoi(Splitted[1])
		if err2 != nil {
			log.Fatalf("Numbers Only Allowed! Got: %s", Splitted[1])
		}

		Rng := RangeLoop{}
		Rng.Position = -1 + n1
		Rng.Length = n2
		if count == 1 {
			Rng.ZeroPad = len(Splitted[1])
		}

		return Range, &Rng

	} else {
		name = strings.Trim(name, "$")

		listPtr := p.Lists[name]
		if listPtr == nil {
			return 0, nil
		}

		new_list := ListLoop{
			Position: -1,
			List:     listPtr,
			Length:   len(*listPtr),
		}

		return List, &new_list
	}
}
