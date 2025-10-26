package main

import "fmt"

const (
	Major int8 = 1
	Minor int8 = 0
	Patch int8 = 0
)

func ShowVersion() {
	fmt.Printf("Bat Version:\t%d.%d.%d\n", Major, Minor, Patch)
}
