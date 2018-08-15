package cwenanngeor

import (
	"fmt"
)

type TypeError struct {
	Wanted Type
	Got    Type
}

func (te *TypeError) Error() string {
	return fmt.Sprintf("Type error: Wanted type `%s` but got type `%s`.", te.Wanted, te.Got)
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

func InferType(node Node, typeWorld map[string]Type) (Type, error) {
	switch node.(type) {
	case *SExpNode:
		sexp := node.(*SExpNode)

		ts, ok := typeWorld[sexp.FuncName]

		if !ok {
			return InvalidType, fmt.Errorf("Func `%s` does not exist.", sexp.FuncName)
		}

		switch ts.(type) {
		case *FuncType:
		default:
			return InvalidType, fmt.Errorf("`%s` is not a function.", sexp.FuncName)
		}

		funcType := ts.(*FuncType)

		for i, exp := range sexp.Exps {
			typ, err := InferType(exp, typeWorld)

			if err != nil {
				return InvalidType, err
			}

			if !TypeEqual(typ, funcType.ArgTypes[i]) {
				return InvalidType, &TypeError{
					Wanted: funcType.ArgTypes[i],
					Got:    typ,
				}
			}
		}
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

			scope := make(map[string]Type)

			for _, node := range fv.Body {
				_, err := InferType(node, scope)

				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
