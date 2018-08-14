package cwenanngeor

type Node interface {
	IsNode() bool
}

type Module struct {
	Name  string
	Funcs []*FuncNode
}

type Arg struct {
	Name string
	Type Type
}

type Type struct {
}

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

func ASTEqual(n1 Node, n2 Node) bool {
	switch n1.(type) {
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
