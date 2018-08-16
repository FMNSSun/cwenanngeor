package cwenanngeor

import (
	"testing"
)

func testInferType(t *testing.T) {
	checkInferedTypeSExp("(dummy 5)", &PrimType{Type: "int"}, t)
	mustErrorInferedTypeSExp("(dummy 5.0)",
		&PrimType{Type: "int"},
		&PrimType{Type: "float"}, t)
}

func mustErrorInferedTypeSExp(code string, wanted, got Type, t *testing.T) {
	p := NewParser(NewTokenizerString(code))
	n, err := p.parseSExp()

	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
		return
	}

	_, err = InferType(n, NewTypeWorlds(builtins))

	if err != nil {
		t.Fatalf("Expected error but got none.")
		return
	}

	switch err.(type) {
	case *TypeError:
		te := err.(*TypeError)
		if !TypeEqual(te.Wanted, wanted) {
			t.Fatalf("Wanted mismatch %s != %s.", te.Wanted, wanted)
			return
		}

		if !TypeEqual(te.Got, got) {
			t.Fatalf("Got mismatch %s != %s.", te.Got, got)
			return
		}
	default:
		t.Fatalf("Unexpected error: %s", err.Error())
		return
	}
}

func checkInferedTypeSExp(code string, exp Type, t *testing.T) {
	p := NewParser(NewTokenizerString(code))
	n, err := p.parseSExp()

	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
		return
	}

	typ, err := InferType(n, NewTypeWorlds(builtins))

	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
		return
	}

	if !TypeEqual(typ, exp) {
		t.Fatalf("Expected type %s but got %s", exp, typ)
		return
	}
}
