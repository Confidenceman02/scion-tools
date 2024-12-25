package internal

import (
	"cmp"
	"fmt"
	"github.com/Confidenceman02/scion-tools/pkg/basics"
	"reflect"
)

func CmpHelp(x any, y any) int {
	switch x1 := x.(type) {
	case basics.Int:
		switch y1 := y.(type) {
		case basics.Int:
			return cmp.Compare(x1, y1)
		default:
			panic("Not an Int")
		}
	case basics.Float:
		switch y1 := y.(type) {
		case basics.Float:
			return cmp.Compare(x1, y1)
		default:
			panic("Not a Float")
		}
	default:
		panic(fmt.Sprintf("Cmp Not implemented for: %v", reflect.TypeOf(x1)))
	}
}

type Cons_[A, B any] struct {
	A A
	B B
}
