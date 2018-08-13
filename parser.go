package cwenanngeor

import (
	"fmt"
)

type Parser struct {
	tz Tokenizer
	tk *Token
}

type ParserError struct {
	Token *Token
	Msg   string
}

func (pe *ParserError) Error() string {
	return fmt.Sprintf("Parser error (file: %s, line: %d): %s",
		pe.Token.Pos.FilePath,
		pe.Token.Pos.LineNumber,
		pe.Msg)
}

func NewParser(tz Tokenizer) *Parser {
	return &Parser{
		tz: tz,
		tk: nil,
	}
}

func (p *Parser) read() (*Token, error) {
	if p.tk != nil {
		it := p.tk
		p.tk = nil
		return it, nil
	}

	tk, err := p.tz.Next()

	if err != nil {
		return nil, err
	}

	return tk, nil
}

func (p *Parser) unread(tk *Token) {
	if p.tk != nil {
		panic("BUG p.tk is not nil")
	}

	p.tk = tk
}

func (p *Parser) parseSExp() (*SExpNode, error) {
	// Next token must be an LPAREN
	tk, err := p.read()

	if err != nil {
		return nil, err
	}

	if tk.Type != TT_LPAREN {
		return nil, &ParserError{
			Token: tk,
			Msg:   fmt.Sprintf("Expected `(` but got `%s`.", tk.SVal),
		}
	}

	// Then followed by at least one IDENT.
	tk, err = p.read()

	if err != nil {
		return nil, err
	}

	if tk.Type != TT_IDENT {
		return nil, &ParserError{
			Token: tk,
			Msg:   fmt.Sprintf("Expected identifier but got `%s`.", tk.SVal),
		}
	}

	nodes := make([]Node, 256)
	nj := 0

	for {
		tk, err = p.read()

		if err != nil {
			return nil, err
		}

		if tk.Type == TT_LPAREN {
			p.unread(tk)
			node, err := p.parseSExp()

			if err != nil {
				return nil, err
			}

			nodes[nj] = node
			nj++
		}
	}

	return nil, nil
}
