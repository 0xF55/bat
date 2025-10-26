package main

import (
	"math/rand"
	"strconv"
	"strings"
	"unsafe"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const special = "#!_*@="

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
	l := len(special)
	for i := 0; i < num; i++ {
		n := rand.Intn(l)
		str += replacer.Replace(string(special[n]))
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
	l := len(charset)
	for i := 0; i < num; i++ {
		n := rand.Intn(l)
		str += string(charset[n])
	}

	return str
}

func Reverse(s string) string {
	n := len(s)
	out := make([]byte, n)
	for i := 0; i < n; i++ {
		out[n-1-i] = s[i]
	}
	return unsafe.String(&out[0], n)
}
