package cwenanngeor

import (
	"fmt"
	"strconv"
)

type Parser struct {
	tz    Tokenizer
	tkbuf []*Token
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
		tz:    tz,
		tkbuf: make([]*Token, 0),
	}
}

func (p *Parser) readbuf() *Token {
	if len(p.tkbuf) == 0 {
		return nil
	}

	it := p.tkbuf[0]
	p.tkbuf = p.tkbuf[1:]

	return it
}

func (p *Parser) read() (*Token, error) {
	it := p.readbuf()

	if it != nil {
		return it, nil
	}

	tk, err := p.tz.Next()

	if err != nil {
		return nil, err
	}

	return tk, nil
}

func (p *Parser) unread(tk *Token) {
	p.tkbuf = append(p.tkbuf, tk)
}

func (p *Parser) parseData() (Node, error) {
	// Next token must be LITINT or LITFLOAT or IDENT.
	tk, err := p.read()

	if err != nil {
		return nil, err
	}

	switch tk.Type {
	case TT_LITINT:
		iv, err := strconv.ParseInt(tk.SVal, 10, 64)

		if err != nil {
			return nil, &ParserError{
				Token: tk,
				Msg:   fmt.Sprintf("`%s` is not a valid integer literal.", tk.SVal),
			}
		}

		return &LitIntNode{
			Value: iv,
			Token: tk,
		}, nil
	case TT_LITFLOAT:
		fv, err := strconv.ParseFloat(tk.SVal, 64)

		if err != nil {
			return nil, &ParserError{
				Token: tk,
				Msg:   fmt.Sprintf("`%s` is not a valid float literal.", tk.SVal),
			}
		}

		return &LitFloatNode{
			Value: fv,
			Token: tk,
		}, nil
	case TT_IDENT:
		return &ReadVarNode{
			Name:  tk.SVal,
			Token: tk,
		}, nil
	case TT_QUOT:
		// Next token must be IDENT
		tk, err = p.read()

		if tk.Type != TT_IDENT {
			return nil, &ParserError{
				Token: tk,
				Msg:   fmt.Sprintf("Expected identifier but got `%s`.", tk.SVal),
			}
		}

		return &QuotNode{
			Ident: tk.SVal,
			Token: tk,
		}, nil
	default:
		return nil, &ParserError{
			Token: tk,
			Msg:   fmt.Sprintf("Expected literal but got `%s`.", tk.SVal),
		}
	}
}

func (p *Parser) parseArg() (Arg, error) {
	tk, err := p.read()

	if err != nil {
		return VoidArg, nil
	}

	if tk.Type != TT_LPAREN {
		return VoidArg, &ParserError{
			Token: tk,
			Msg:   fmt.Sprintf("Expected `(` but got `%s`.", tk.SVal),
		}
	}

	tk, err = p.read()

	if err != nil {
		return VoidArg, nil
	}

	if tk.Type != TT_IDENT {
		return VoidArg, &ParserError{
			Token: tk,
			Msg:   fmt.Sprintf("Expected identifier but got `%s`.", tk.SVal),
		}
	}

	aname := tk.SVal

	tp, err := p.parseType()

	if err != nil {
		return VoidArg, err
	}

	tk, err = p.read()

	if err != nil {
		return VoidArg, err
	}

	if tk.Type != TT_RPAREN {
		return VoidArg, &ParserError{
			Token: tk,
			Msg:   fmt.Sprintf("Expected `)` but got `%s`.", tk.SVal),
		}
	}

	return Arg{
		Name: aname,
		Type: tp,
	}, nil
}

func (p *Parser) parseType() (Type, error) {
	tk, err := p.read()

	if err != nil {
		return InvalidType, err
	}

	if tk.Type == TT_IDENT {
		return &PrimType{
			Type: tk.SVal,
		}, nil
	}

	return InvalidType, &ParserError{
		Token: tk,
		Msg:   fmt.Sprintf("`%s` is not a type.", tk.SVal),
	}
}

func (p *Parser) Funcs() ([]*FuncNode, error) {
	return p.parseFuncs()
}

func (p *Parser) parseFuncs() ([]*FuncNode, error) {
	funcs := make([]*FuncNode, 8) // TODO: resize later
	fj := 0

	for {
		tk, err := p.read()

		if err != nil {
			return nil, err
		}

		if tk.Type == TT_EOF {
			break
		}

		if tk.Type != TT_LPAREN {
			return nil, &ParserError{
				Token: tk,
				Msg:   fmt.Sprintf("Expected `(` but got `%s`.", tk.SVal),
			}
		}

		p.unread(tk)

		fn, err := p.parseFunc()

		if err != nil {
			return nil, err
		}

		fn_, ok := fn.(*FuncNode)

		if !ok {
			panic("BUG: didn't get *FuncNode")
		}

		funcs[fj] = fn_
		fj++
	}

	return funcs[:fj], nil
}

func (p *Parser) parseFunc() (Node, error) {
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

	// then the next token must be FUNC

	tk, err = p.read()

	if err != nil {
		return nil, err
	}

	if tk.Type != TT_FUNC {
		return nil, &ParserError{
			Token: tk,
			Msg:   fmt.Sprintf("Expected `func` but got `%s`.", tk.SVal),
		}
	}

	// then the next token must be IDENT

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

	args := make([]Arg, 8) // TODO: resize later
	aj := 0

	// then the arguments follow. which is at least one LPAREN then until RPAREN
	tk, err = p.read()

	if err != nil {
		return nil, err
	}

	if tk.Type != TT_LPAREN {
		return nil, &ParserError{
			Token: tk,
			Msg:   fmt.Sprintf("Expected `(` but got `%s`.", tk.SVal),
		}
	}

	for {
		done := false

		tk, err = p.read()

		switch tk.Type {
		case TT_RPAREN:
			done = true
		case TT_LPAREN:
			p.unread(tk)

			arg, err := p.parseArg()

			if err != nil {
				return nil, err
			}

			args[aj] = arg
			aj++
		default:
			return nil, &ParserError{
				Token: tk,
				Msg:   fmt.Sprintf("Expected `(` or `)` but got `%s`", tk.SVal),
			}
		}

		if done {
			break
		}
	}

	ret, err := p.parseType()

	if err != nil {
		return nil, err
	}

	bodies := make([]Node, 8) // TODO: resize later
	bj := 0

	for {
		done := false

		tk, err = p.read()

		if err != nil {
			return nil, err
		}

		switch tk.Type {
		case TT_RPAREN:
			done = true
		case TT_LPAREN:
			p.unread(tk)

			sexp, err := p.parseSExp()

			if err != nil {
				return nil, err
			}

			bodies[bj] = sexp
			bj++
		}

		if done {
			break
		}
	}

	return &FuncNode{
		Args:    args[:aj],
		RetType: ret,
		Body:    bodies[:bj],
		Token:   firsttk,
		Name:    funcname,
	}, nil
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
		case TT_LITINT, TT_LITFLOAT, TT_IDENT, TT_QUOT:
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
