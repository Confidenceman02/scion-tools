package bitwise

import (
	. "github.com/Confidenceman02/scion-tools/pkg/basics"
)

func And(a, b Int) Int {
	return a & b
}

// Bit Shifts

// Shift bits to the right by a given offset, filling new bits with whatever is the topmost bit. This can be used to
// divide numbers by powers of two.
func ShiftRightBy(offset Int, a Int) Int {
	return a >> offset
}
