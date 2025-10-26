package main

import (
	"bufio"
	"strconv"
)

type ListLoop struct {
	List     *[]string
	Position int // start
	Length   int // end
}

type RangeLoop struct {
	Position int // start
	Length   int // end
}

type FileLoop struct {
	Scanner *bufio.Scanner
}

const (
	List int8 = iota + 1
	Range
	File
)

type LoopCTX struct {
	Iterator string
	Type     int8
	Expr     interface{}
}

func (ctx *LoopCTX) Next() bool {

	switch ctx.Type {
	case List:
		listLoop := ctx.Expr.(*ListLoop)
		listLoop.Position++
		if listLoop.Position >= listLoop.Length || listLoop.Length == 0 {
			return false
		}
		Variables[ctx.Iterator] = (*listLoop.List)[listLoop.Position]
		return true

	case Range:
		rangeLoop := ctx.Expr.(*RangeLoop)
		rangeLoop.Position++
		if rangeLoop.Position >= rangeLoop.Length || rangeLoop.Length == 0 {
			return false
		}
		num := strconv.Itoa(rangeLoop.Position)
		Variables[ctx.Iterator] = num
		return true

	case File:
		fileLoop := ctx.Expr.(*FileLoop)
		if fileLoop.Scanner.Scan() {
			Variables[ctx.Iterator] = fileLoop.Scanner.Text()
			return true
		}
		return false
	default:
		return false
	}
}
