package cwenanngeor

type Node interface {
	IsNode() bool
}

type Arg struct {
	Name string
	Type Type
}

type Type struct {
}

type FuncNode struct {
	Name string
	Args []Arg
	Body []Node
}

type IfNode struct {
	Condition Node
	Block     []Node
}

type SExpNode struct {
	FuncName string
	Exps     []Node
}

func (*SExpNode) IsNode() bool {
	return true
}
