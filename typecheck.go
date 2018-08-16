package cwenanngeor

import (
	"fmt"
)

type TypeError struct {
	Wanted Type
	Got    Type
	Token  *Token
}

type TypeWorld map[string]Type
type TypeWorlds []TypeWorld

func (tws TypeWorlds) Lookup(val string) Type {
	for i := len(tws) - 1; i >= 0; i-- {
		it := tws[i][val]

		if it != nil {
			return it
		}
	}

	return nil
}

func NewTypeWorlds(typeWorlds ...TypeWorld) TypeWorlds {
	return typeWorlds
}

func (te *TypeError) Error() string {
	return fmt.Sprintf("Type error (file: %q, line: %d): Wanted type `%s` but got type `%s`.",
		te.Token.Pos.FilePath, te.Token.Pos.LineNumber, te.Wanted, te.Got)
}

var builtins map[string]Type = map[string]Type{
	"cast.i.i": &FuncType{
		ArgTypes: []Type{
			&PrimType{
				Type: "int",
			},
		},
		RetType: &PrimType{
			Type: "int",
		},
	},
}

func InferType(node Node, typeWorlds TypeWorlds) (Type, error) {
	switch node.(type) {
	case *SExpNode:
		sexp := node.(*SExpNode)

		ts := typeWorlds.Lookup(sexp.FuncName)

		if ts == nil {
			return InvalidType, fmt.Errorf("`%s` does not exist.", sexp.FuncName)
		}

		switch ts.(type) {
		case *FuncType:
		default:
			return InvalidType, fmt.Errorf("`%s` is not a function.", sexp.FuncName)
		}

		funcType := ts.(*FuncType)

		for i, exp := range sexp.Exps {
			typ, err := InferType(exp, typeWorlds)

			if err != nil {
				return InvalidType, err
			}

			if !TypeEqual(typ, funcType.ArgTypes[i]) {
				return InvalidType, &TypeError{
					Wanted: funcType.ArgTypes[i],
					Got:    typ,
					Token:  sexp.Token,
				}
			}
		}

		return funcType.RetType, nil
	}

	return InvalidType, fmt.Errorf("Can't infer type.")
}

func TypeCheck(modules map[string]*Module) error {
	for k, v := range modules {
		if k != v.Name {
			panic("BUG: TypeCheck 1.")
		}

		for fk, fv := range v.Funcs {
			if fk != fv.Name {
				panic("BUG: TypeCheck 2.")
			}

			typeWorlds := NewTypeWorlds(builtins)

			for _, node := range fv.Body {
				_, err := InferType(node, typeWorlds)

				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
