package main

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func GetVar(name string) string {

	name = strings.Trim(name, "$")
	switch name {
	case "S":
		return " "
	case "D":
		return "`DS`" // dolar $
	case "E":
		return "`EQ`" // equal =
	case "A":
		return "`AT`" // at @
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
		ret = string(Variables[name])
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

	replacer := strings.NewReplacer("`DS`", "$", "`EQ`", "=", "`AT`", "@")

	return replacer.Replace(text)

}

func EvalExpression(text string) string {
	varegex := regexp.MustCompile(`\$\s*[\+\-\~\^\!\*]?\w+(?::\d+)?`) // $?var || $gen:n
	comma_sp := strings.Split(text, ",")
	var ret strings.Builder

	for _, expr := range comma_sp {
		expr = varegex.ReplaceAllStringFunc(expr, func(match string) string {
			return GetVar(match)
		})
		ret.WriteString(expr)
	}
	return ret.String()
}

func Eval(text string) string {
	ret := UnEscape(EvalExpression(text))
	return ret
}

func GetList(name string) (int8, any) {

	if strings.Contains(name, "..") {
		Splitted := strings.Split(name, "..")
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

		return Range, &Rng

	} else {
		name = strings.Trim(name, "$")

		listPtr := Lists[name]
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
