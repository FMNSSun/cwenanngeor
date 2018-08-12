package cwenanngeor

import (
	"testing"
)

func TestEmpty(t *testing.T) {
}

func TestSimple(t *testing.T) {
	checkTypes("func", []TokenType{TT_FUNC}, t)
	checkTypes("  func", []TokenType{TT_FUNC}, t)
	checkTypes(" func ", []TokenType{TT_FUNC}, t)
	checkTypes(" func\n func\n ", []TokenType{TT_FUNC, TT_FUNC}, t)
}

func TestSpecials(t *testing.T) {
	checkTypes("{}", []TokenType{TT_BEGIN, TT_END}, t)
	checkTypes("func()", []TokenType{TT_FUNC, TT_LPAREN, TT_RPAREN}, t)
	checkTypes(" ; ", []TokenType{TT_SEMICOLON}, t)
}

func TestLits(t *testing.T) {
	checkTypes("5", []TokenType{TT_LITINT}, t)
	checkTypes("5.0", []TokenType{TT_LITFLOAT}, t)
	checkTypes("1.", []TokenType{TT_LITFLOAT}, t)
	checkTypes("5func", []TokenType{TT_LITINT, TT_FUNC}, t)
	checkTypes("5.func", []TokenType{TT_LITFLOAT, TT_FUNC}, t)
	checkTypes("5.1func", []TokenType{TT_LITFLOAT, TT_FUNC}, t)
}

func checkTypes(str string, tts []TokenType, t *testing.T) {
	tz := NewTokenizerString(str)

	for _, v := range tts {
		tk, err := tz.Next()

		if err != nil {
			t.Fatalf("Unexpected error: %q", err.Error())
			return
		}

		if tk.Type != v {
			t.Fatalf("Expected TT %d but got %d: %q", v, tk.Type, str)
			return
		}
	}

	tk, err := tz.Next()

	if err != nil {
		t.Fatalf("Unexpected error: %q", err.Error())
		return
	}

	if tk != nil {
		t.Fatalf("Expected EOF but still got a token")
		return
	}
}
