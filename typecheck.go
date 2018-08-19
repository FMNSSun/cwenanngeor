package cwenanngeor

import (
	"fmt"
)

type TypeError struct {
	Wanted Type
	Got    Type
	Token  *Token
	Extra  string
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
	if te.Extra == "" {
		return fmt.Sprintf("Type error %s: Wanted type `%s` but got type `%s`.",
			te.Token.Pos, te.Wanted, te.Got)
	} else {
		return fmt.Sprintf("Type error %s %s: Wanted type `%s` but got type `%s`.",
			te.Extra, te.Token.Pos, te.Wanted, te.Got)
	}
}

var builtins map[string]Type = map[string]Type{
	"square.i": &FuncType{
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

func TypeCompatibleWith(a Type, b Type) bool {
	switch a.(type) {
	case *VoidType:
		switch b.(type) {
		case *VoidType:
			return true
		default:
			return false
		}
	case *PrimType:
		switch b.(type) {
		case *PrimType:
			return TypeEqual(a, b)
		case *UnionType:
			ut := b.(*UnionType)

			for _, typ := range ut.Types {
				if TypeEqual(a, typ) {
					return true
				}
			}

			return false
		default:
			return false
		}
	case *UnionType:
		switch b.(type) {
		case *UnionType:
			// all types of a must be types of b as well.
			ut_a := a.(*UnionType)
			ut_b := b.(*UnionType)

			for _, typ_a := range ut_a.Types {
				found := false
				for _, typ_b := range ut_b.Types {
					if TypeEqual(typ_a, typ_b) {
						found = true
						break
					}
				}

				if !found {
					return false
				}
			}

			return true
		default:
			return false
		}
	}

	panic("BUG: Can't tell if compatible or not?")
}

func InferType(node Node, typeWorlds TypeWorlds) (Type, error) {
	switch node.(type) {
	case *LitFloatNode:
		return &PrimType{Type: "float"}, nil
	case *LitIntNode:
		return &PrimType{Type: "int"}, nil
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
					Extra:  fmt.Sprintf("in a call to `%s`", sexp.FuncName),
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

		for _, fn := range v.Funcs {
			typeWorlds := NewTypeWorlds(builtins)

			var lastType Type = &VoidType{}

			for _, node := range fn.FuncNode.Body {
				typ, err := InferType(node, typeWorlds)

				if err != nil {
					return err
				}

				lastType = typ
			}

			if !TypeCompatibleWith(lastType, fn.FuncNode.RetType) {
				return &TypeError{
					Wanted: fn.FuncNode.RetType,
					Got:    lastType,
					Token:  fn.FuncNode.Token,
					Extra:  fmt.Sprintf("(type of last statement is not compatible with return type of the function)"),
				}
			}
		}
	}

	return nil
}
