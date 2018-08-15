package cwenanngeor

import (
	"testing"
)

func TestParseSExp(t *testing.T) {

	checkASTFunc(
		"(func main ((a int)) float)",
		&FuncNode{
			Name:    "main",
			RetType: Type{Kind: TK_PRIM, Type: "float"},
			Body:    []Node{},
			Args: []Arg{
				Arg{
					Type: Type{Kind: TK_PRIM, Type: "int"},
					Name: "a",
				},
			},
		}, t)

	checkASTFunc(
		"(func main ((a int)) float (add 5 6) (sub 4 7))",
		&FuncNode{
			Name:    "main",
			RetType: Type{Kind: TK_PRIM, Type: "float"},
			Body: []Node{
				&SExpNode{
					FuncName: "add",
					Exps: []Node{
						&LitIntNode{Value: 5},
						&LitIntNode{Value: 6},
					},
				},
				&SExpNode{
					FuncName: "sub",
					Exps: []Node{
						&LitIntNode{Value: 4},
						&LitIntNode{Value: 7},
					},
				},
			},
			Args: []Arg{
				Arg{
					Type: Type{Kind: TK_PRIM, Type: "int"},
					Name: "a",
				},
			},
		}, t)

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

func checkASTFunc(code string, exp Node, t *testing.T) {
	p := NewParser(NewTokenizerString(code))

	n, err := p.parseFunc()

	if err != nil {
		t.Fatalf("Unexpected error: %s.", err.Error())
	}

	if !ASTEqual(n, exp) {
		t.Fatalf("ASTs do not match! %+v %+v", n, exp)
	}
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
