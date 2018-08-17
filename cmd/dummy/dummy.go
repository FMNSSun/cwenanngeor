package main

import (
	"fmt"

	cwe "github.com/FMNSSun/cwenanngeor"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	m, err := cwe.LoadModule("test")
	
	if err != nil {
		panic(err.Error())
	}
	
	spew.Dump(m)
	fmt.Println("")
}