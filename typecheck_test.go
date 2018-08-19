package cwenanngeor

import (
	"testing"
)

func TestInferType(t *testing.T) {
	checkInferedTypeSExp("(square.i 5)", &PrimType{Type: "int"}, t)
	mustErrorInferedTypeSExp("(square.i 5.0)",
		&PrimType{Type: "int"},
		&PrimType{Type: "float"}, t)
}

func mustErrorInferedTypeSExp(code string, wanted, got Type, t *testing.T) {
	p := NewParser(NewTokenizerString(code))
	n, err := p.parseSExp()

	if err != nil {
		t.Fatalf("Unexpected error for %s: %s", code, err.Error())
		return
	}

	typ, err := InferType(n, NewTypeWorlds(builtins))

	if err == nil {
		t.Fatalf("Expected error but got none for: %s. {%s}", code, typ)
		return
	}

	switch err.(type) {
	case *TypeError:
		te := err.(*TypeError)
		if !TypeEqual(te.Wanted, wanted) {
			t.Fatalf("Wanted mismatch %s != %s for %s.", te.Wanted, wanted, code)
			return
		}

		if !TypeEqual(te.Got, got) {
			t.Fatalf("Got mismatch %s != %s for %s.", te.Got, got, code)
			return
		}
	default:
		t.Fatalf("Unexpected error for %s: %s", code, err.Error())
		return
	}
}

func checkInferedTypeSExp(code string, exp Type, t *testing.T) {
	p := NewParser(NewTokenizerString(code))
	n, err := p.parseSExp()

	if err != nil {
		t.Fatalf("Unexpected error for %s: %s", code, err.Error())
		return
	}

	typ, err := InferType(n, NewTypeWorlds(builtins))

	if err != nil {
		t.Fatalf("Unexpected error for %s: %s", code, err.Error())
		return
	}

	if !TypeEqual(typ, exp) {
		t.Fatalf("Expected type %s but got %s for %s.", exp, typ, code)
		return
	}
}
