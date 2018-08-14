package cwenanngeor

import (
	"testing"
)

func TestParseSExp(t *testing.T) {

	checkASTSExp(
		"(foo 'fooo)",
		&SExpNode{
			FuncName: "foo",
			Exps: []Node{
				&QuotNode{Ident: "fooo"},
			},
		}, t)

	checkASTSExp(
		"(foo 5 6)",
		&SExpNode{
			FuncName: "foo",
			Exps: []Node{
				&LitIntNode{Value: 5},
				&LitIntNode{Value: 6},
			},
		}, t)

	checkASTSExp(
		"(foo 5 (bar 1))",
		&SExpNode{
			FuncName: "foo",
			Exps: []Node{
				&LitIntNode{Value: 5},
				&SExpNode{
					FuncName: "bar",
					Exps: []Node{
						&LitIntNode{Value: 1},
					},
				},
			},
		}, t)
}

func checkASTSExp(code string, exp Node, t *testing.T) {
	p := NewParser(NewTokenizerString(code))

	n, err := p.parseSExp()

	if err != nil {
		t.Fatalf("Unexpected error: %s.", err.Error())
	}

	if !ASTEqual(n, exp) {
		t.Fatalf("ASTs do not match! %v %v", n, exp)
	}
}
