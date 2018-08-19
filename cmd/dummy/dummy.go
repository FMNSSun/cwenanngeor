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
	fmt.Printf("\n\n%#v\n\n", m)
	
	err = cwe.TypeCheck(map[string]*cwe.Module{"test":m})
	
	if err != nil {
		panic(err.Error())
	}
}