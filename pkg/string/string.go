package string

import (
	"cmp"

	. "github.com/Confidenceman02/scion-tools/pkg/basics"
	. "github.com/Confidenceman02/scion-tools/pkg/maybe"
)

type String string

func (i String) Cmp(y Comparable[String]) int {
	return cmp.Compare(i, y.T())
}
func (i String) T() String {
	return i
}

// Int Conversions

func ToInt(x String) Maybe[Int] {
	var total int
	code0 := x[0]
	var start int = 0

	if code0 == 0x2B /* + */ || code0 == 0x2D /* - */ {
		start = 1
	}

	var i = start

	for i < len(x) {
		var code = x[i]
		if code < 0x30 || 0x39 < code /* 0 - 9 */ {
			return Nothing{}
		}

		total = 10*total + int(code-0x30)
		i++
	}

	if i == start {
		return Nothing{}
	} else {
		if code0 == 0x2D /* - */ {
			return Just[Int]{Value: Int(-total)}
		} else {
			return Just[Int]{Value: Int(total)}
		}
	}
}
