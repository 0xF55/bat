package main

import (
	"log"
	"strings"
)

func ExecuteDirective(expr string) {

	sp := strings.Split(expr, ":")

	if len(sp) == 2 && sp[1] != "" {

		switch sp[0] {
		case "charset":
			mode := sp[1]
			switch mode {
			case "lower":
				Charset = Lower
			case "upper":
				Charset = Upper
			case "all":
				Charset = Lower + Upper
			default:
				// excpected custom charset here
				Charset = mode
			}
		case "special":
			Special = sp[1]
		}

	} else {
		log.Fatalf("Invalid Directive: %s", expr)
	}
}
