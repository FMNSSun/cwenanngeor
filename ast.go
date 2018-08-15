package cwenanngeor

type Node interface {
	IsNode() bool
}

type Arg struct {
	Name string
	Type Type
}

var VoidArg Arg = Arg{}

type Type struct {
	Kind TypeKind
	Type string
}

type TypeKind uint8

const TK_PRIM = TypeKind(0)
const TK_LIST = TypeKind(1)

var VoidType Type = Type{}

type FuncNode struct {
	Name    string
	Args    []Arg
	Body    []Node
	RetType Type
	Token   *Token
}

func (*FuncNode) IsNode() bool {
	return true
}

type LitIntNode struct {
	Value int64
	Token *Token
}

func (*LitIntNode) IsNode() bool {
	return true
}

type QuotNode struct {
	Ident string
	Token *Token
}

func (*QuotNode) IsNode() bool {
	return true
}

type LitFloatNode struct {
	Value float64
	Token *Token
}

func (*LitFloatNode) IsNode() bool {
	return true
}

type ReadVarNode struct {
	Name  string
	Token *Token
}

func (*ReadVarNode) IsNode() bool {
	return true
}

type IfNode struct {
	Condition Node
	Block     []Node
}

type SExpNode struct {
	FuncName string
	Exps     []Node
	Token    *Token
}

func (*SExpNode) IsNode() bool {
	return true
}

func TypeEqual(t1 Type, t2 Type) bool {
	return t1.Kind == t2.Kind && t1.Type == t2.Type
}

func ArgEqual(a1 Arg, a2 Arg) bool {
	return a1.Name == a2.Name && TypeEqual(a1.Type, a2.Type)
}

func ASTEqual(n1 Node, n2 Node) bool {
	switch n1.(type) {
	case *FuncNode:
		switch n2.(type) {
		case *FuncNode:
			fn1 := n1.(*FuncNode)
			fn2 := n2.(*FuncNode)

			if fn1.Name != fn2.Name {
				return false
			}

			if len(fn1.Args) != len(fn2.Args) {
				return false
			}

			if len(fn1.Body) != len(fn2.Body) {
				return false
			}

			if !TypeEqual(fn1.RetType, fn2.RetType) {
				return false
			}

			for i := 0; i < len(fn1.Args); i++ {
				if !ArgEqual(fn1.Args[i], fn2.Args[i]) {
					return false
				}
			}

			for i := 0; i < len(fn1.Body); i++ {
				if !ASTEqual(fn1.Body[i], fn2.Body[i]) {
					return false
				}
			}

			return true
		default:
			return false
		}
	case *LitFloatNode:
		switch n2.(type) {
		case *LitFloatNode:
			return n1.(*LitFloatNode).Value == n2.(*LitFloatNode).Value
		default:
			return false
		}
	case *LitIntNode:
		switch n2.(type) {
		case *LitIntNode:
			return n1.(*LitIntNode).Value == n2.(*LitIntNode).Value
		default:
			return false
		}
	case *ReadVarNode:
		switch n2.(type) {
		case *ReadVarNode:
			return n1.(*ReadVarNode).Name == n2.(*ReadVarNode).Name
		default:
			return false
		}
	case *QuotNode:
		switch n2.(type) {
		case *QuotNode:
			return n1.(*QuotNode).Ident == n2.(*QuotNode).Ident
		default:
			return false
		}
	case *SExpNode:
		switch n2.(type) {
		case *SExpNode:
			n1_ := n1.(*SExpNode)
			n2_ := n2.(*SExpNode)

			if len(n1_.Exps) != len(n2_.Exps) {
				return false
			}

			sz := len(n1_.Exps)

			for i := 0; i < sz; i++ {
				if !ASTEqual(n1_.Exps[i], n2_.Exps[i]) {
					return false
				}
			}

			return n1_.FuncName == n2_.FuncName
		default:
			return false
		}
	}

	panic("BUG: ASTEqual")
}
