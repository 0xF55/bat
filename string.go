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
	"math/rand"
	"strconv"
	"strings"
	"unsafe"
)

var Lower = "abcdefghijklmnopqrstuvwxyz"
var Upper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
var Charset string
var Special string

func charUpper(char byte) string {
	ret := strings.ToUpper(string(char))
	return ret

}

func charLower(char byte) string {
	ret := strings.ToLower(string(char))
	return ret

}

/*Capitaltize String e.g. hello -> Hello | Symbol: ^ */
func Capitaltize(text string) string {

	ret := ""
	ret += charUpper(text[0])
	ret += text[1:]

	return ret

}

/*ZigZag string e.g. helloo -> hElLo | Symbol: ~ */

func ZigZag(text string) string {
	ret := ""
	for i, c := range text {
		if i%2 == 0 {
			ret += charLower(byte(c))
			continue
		}
		ret += charUpper(byte(c))
	}
	return ret
}

// RND:5
func RandNum(expr string) string {
	s := strings.Split(expr, ":")[1]
	num, err := strconv.Atoi(s)
	if err != nil || num == 0 {
		return ""
	}
	min := 1
	for i := 0; i < num-1; i++ {
		min *= 10
	}
	max := (min * 10) + 1

	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func RandSpecial(expr string) string {
	str := ""
	s := strings.Split(expr, ":")[1]
	num, err := strconv.Atoi(s)
	if err != nil || num == 0 {
		return ""
	}
	replacer := strings.NewReplacer("$", "`DS`", "=", "`EQ`", "@", "`AT`")
	l := len(Special)
	for i := 0; i < num; i++ {
		n := rand.Intn(l)
		str += replacer.Replace(string(Special[n]))
	}

	return str
}

func RandChar(expr string) string {
	str := ""
	s := strings.Split(expr, ":")[1]
	num, err := strconv.Atoi(s)
	if err != nil || num == 0 {
		return ""
	}
	l := len(Charset)
	for i := 0; i < num; i++ {
		n := rand.Intn(l)
		str += string(Charset[n])
	}

	return str
}

func Reverse(s string) string {
	cached := ReversedCache[s]
	if cached != "" {
		return cached
	}
	n := len(s)
	out := make([]byte, n)
	for i := 0; i < n; i++ {
		out[n-1-i] = s[i]
	}
	str := unsafe.String(&out[0], n)
	ReversedCache[s] = str
	return str
}
