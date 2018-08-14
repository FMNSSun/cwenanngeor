package cwenanngeor

import (
	"testing"
)

func TestParseSExp(t *testing.T) {
	code := "(foo 5 6)"

	p := NewParser(NewTokenizerString(code))

	n, err := p.parseSExp()

	if err != nil {
		t.Fatalf("Unexpected error: %s.", err.Error())
	}

	exp := &SExpNode{
		FuncName: "foo",
		Exps: []Node{
			&LitIntNode{SVal: "5"},
			&LitIntNode{SVal: "6"},
		},
	}

	if !ASTEqual(n, exp) {
		t.Fatalf("ASTs do not match! %v %v", n, exp)
	}
}
