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

func checkTypes(str string, tts []TokenType, t *testing.T) {
    tz := NewTokenizerString(str)
    
    for _, v := range tts {
        tk, err := tz.Next()
        
        if err != nil {
            t.Fatalf("Unexpected error: %q", err.Error())
            return
        }
        
        if tk.Type != v {
            t.Fatalf("Expected TT %d but got %d", v, tk.Type)
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