package cwenanngeor

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Token struct {
	SVal string
	Type TokenType
	Pos  *FilePos
}

type FilePos struct {
	LineNumber uint32
	FilePath   string
}

type TokenType uint8

const TT_FUNC = TokenType(1)
const TT_BEGIN = TokenType(2)
const TT_END = TokenType(3)
const TT_IDENT = TokenType(4)
const TT_LITINT = TokenType(5)
const TT_SEMICOLON = TokenType(6)
const TT_IF = TokenType(7)
const TT_ELSE = TokenType(8)
const TT_LPAREN = TokenType(9)
const TT_RPAREN = TokenType(10)
const TT_LITFLOAT = TokenType(11)
const TT_EOF = TokenType(12)
const TT_QUOT = TokenType(13)

type Tokenizer interface {
	Next() (*Token, error)
}

type tokenizer struct {
	r      *bufio.Reader
	fpath  string
	lineno uint32
	rn     rune
}

func NewTokenizerString(code string) Tokenizer {
	return &tokenizer{
		r:      bufio.NewReader(strings.NewReader(code)),
		fpath:  "<memory>",
		lineno: 1,
		rn:     eof,
	}
}

type TokenizerError struct {
	Pos *FilePos
	Err error
}

func (te *TokenizerError) Error() string {
	return fmt.Sprintf("Error (file: %q, line: %d): %s", te.Pos.FilePath,
		te.Pos.LineNumber, te.Err.Error(),
	)
}

var eof = rune(0)

func (t *tokenizer) filepos() *FilePos {
	return &FilePos{
		LineNumber: t.lineno,
		FilePath:   t.fpath,
	}
}

func (t *tokenizer) unread(rn rune) {
	if t.rn != eof {
		panic("BUG t.rn is not empty!")
	}

	t.rn = rn
}

func (t *tokenizer) read() (rune, error) {
	if t.rn != eof {
		it := t.rn
		t.rn = eof
		return it, nil
	}

	rn, _, err := t.r.ReadRune()

	if err == io.EOF {
		return eof, nil
	}

	if err != nil {
		return rn, err
	}

	if rn == '\n' {
		t.lineno++
	}

	return rn, nil
}

func (t *tokenizer) litintfloat(rn rune) (*Token, error) {
	var buf bytes.Buffer
	buf.WriteRune(rn)

	seenDot := false

	for {
		rn, err := t.read()

		if err != nil {
			return nil, err
		}

		if rn == eof {
			break
		}

		if isdigit(rn) {
			buf.WriteRune(rn)
		} else if rn == '.' {
			if seenDot {
				return nil, &TokenizerError{
					Pos: t.filepos(),
					Err: err,
				}
			}
			seenDot = true
			buf.WriteRune(rn)
		} else {
			t.unread(rn)
			break
		}
	}

	str := buf.String()

	if str == "-" {
		return nil, &TokenizerError{
			Pos: t.filepos(),
			Err: fmt.Errorf("Literal %q is missing at least one digit.", str),
		}
	}

	if seenDot {
		return &Token{
			SVal: str,
			Type: TT_LITFLOAT,
			Pos:  t.filepos(),
		}, nil

	} else {
		return &Token{
			SVal: str,
			Type: TT_LITINT,
			Pos:  t.filepos(),
		}, nil
	}
}

func (t *tokenizer) ident(rn rune) (*Token, error) {
	var buf bytes.Buffer
	buf.WriteRune(rn)

	for {
		rn, err := t.read()

		if err != nil {
			return nil, err
		}

		if rn == eof {
			break
		}

		if isident(rn) {
			buf.WriteRune(rn)
		} else {
			t.unread(rn)
			break
		}
	}

	str := buf.String()

	switch str {
	case "func":
		return &Token{
			SVal: str,
			Type: TT_FUNC,
			Pos:  t.filepos(),
		}, nil
	}

	return &Token{
		SVal: str,
		Type: TT_IDENT,
		Pos:  t.filepos(),
	}, nil
}

func (t *tokenizer) Next() (*Token, error) {
	var rn rune
	var err error

	for {
		rn, err = t.read()

		if err != nil {
			return nil, &TokenizerError{
				Pos: t.filepos(),
				Err: err,
			}
		}

		if !iswhitespace(rn) || rn == eof {
			break
		}
	}

	switch rn {
	case ';':
		return &Token{
			SVal: ";",
			Type: TT_SEMICOLON,
			Pos:  t.filepos(),
		}, nil
	case '{':
		return &Token{
			SVal: "{",
			Type: TT_BEGIN,
			Pos:  t.filepos(),
		}, nil
	case '}':
		return &Token{
			SVal: "}",
			Type: TT_END,
			Pos:  t.filepos(),
		}, nil
	case '(':
		return &Token{
			SVal: "(",
			Type: TT_LPAREN,
			Pos:  t.filepos(),
		}, nil
	case ')':
		return &Token{
			SVal: ")",
			Type: TT_RPAREN,
			Pos:  t.filepos(),
		}, nil
	case '\'':
		return &Token{
			SVal: "'",
			Type: TT_QUOT,
			Pos:  t.filepos(),
		}, nil
	}

	if isletter(rn) {
		return t.ident(rn)
	} else if isdigit(rn) || rn == '-' {
		return t.litintfloat(rn)
	}

	if rn == eof {
		return &Token{
			SVal: "<eof>",
			Type: TT_EOF,
			Pos:  t.filepos(),
		}, nil
	}

	return nil, &TokenizerError{
		Pos: t.filepos(),
		Err: fmt.Errorf("Unexpected rune: %v", rn),
	}
}
