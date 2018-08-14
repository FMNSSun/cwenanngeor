package cwenanngeor

import (
	"testing"
)

func TestParseSExp(t *testing.T) {

	checkAST(
		"(foo 5 6)",
		&SExpNode{
			FuncName: "foo",
			Exps: []Node{
				&LitIntNode{SVal: "5"},
				&LitIntNode{SVal: "6"},
			},
		}, t)

	checkAST(
		"(foo 5 (bar 1))",
		&SExpNode{
			FuncName: "foo",
			Exps: []Node{
				&LitIntNode{SVal: "5"},
				&SExpNode{
					FuncName: "bar",
					Exps: []Node{
						&LitIntNode{SVal: "1"},
					},
				},
			},
		}, t)
}

func checkAST(code string, exp Node, t *testing.T) {
	p := NewParser(NewTokenizerString(code))

	n, err := p.parseSExp()

	if err != nil {
		t.Fatalf("Unexpected error: %s.", err.Error())
	}

	if !ASTEqual(n, exp) {
		t.Fatalf("ASTs do not match! %v %v", n, exp)
	}
}
