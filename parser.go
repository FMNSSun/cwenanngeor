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

func (p *Parser) parseData() (Node, error) {
	// Next token must be LITINT or LITFLOAT or IDENT.
	tk, err := p.read()

	if err != nil {
		return nil, err
	}

	switch tk.Type {
	case TT_LITINT:
		return &LitIntNode{
			SVal:  tk.SVal,
			Token: tk,
		}, nil
	case TT_LITFLOAT:
		return &LitFloatNode{
			SVal:  tk.SVal,
			Token: tk,
		}, nil
	case TT_IDENT:
		return &ReadVarNode{
			Name:  tk.SVal,
			Token: tk,
		}, nil
	default:
		return nil, &ParserError{
			Token: tk,
			Msg:   fmt.Sprintf("Expected literal but got `%s`.", tk.SVal),
		}
	}
}

func (p *Parser) parseSExp() (Node, error) {
	// Next token must be an LPAREN
	tk, err := p.read()

	firsttk := tk // remember for later

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

	funcname := tk.SVal
	nodes := make([]Node, 8) //TODO: resize this later
	nj := 0

	for {
		tk, err = p.read()

		if err != nil {
			return nil, err
		}

		switch tk.Type {
		case TT_LPAREN:
			p.unread(tk)
			node, err := p.parseSExp()

			if err != nil {
				return nil, err
			}

			nodes[nj] = node
			nj++
		case TT_LITINT, TT_LITFLOAT, TT_IDENT:
			p.unread(tk)
			node, err := p.parseData()

			if err != nil {
				return nil, err
			}

			nodes[nj] = node
			nj++
		case TT_RPAREN:
			return &SExpNode{
				FuncName: funcname,
				Exps:     nodes[:nj],
				Token:    firsttk,
			}, nil
		default:
			return nil, &ParserError{
				Token: tk,
				Msg:   fmt.Sprintf("Expected `(`, identifier or literal but got `%s`.", tk.SVal),
			}
		}
	}

	return nil, nil
}
