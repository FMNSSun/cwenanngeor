package cwenanngeor

import (
	"strings"
)

func MangleName(module string, fname string, args []Arg, ret Type) string {
	s := make([]string, 2)
	s[0] = module
	s[1] = fname
	for _, arg := range args {
		s = append(s, MangleType(arg.Type))
	}

	return strings.Join(s, "_")
}

func MangleType(t Type) string {
	switch t.(type) {
	case *PrimType:
		return t.(*PrimType).Type
	}
	panic("BUG: MangleType")
}
