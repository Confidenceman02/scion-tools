package bitwise

import (
	. "github.com/Confidenceman02/scion-tools/pkg/basics"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicOperations(t *testing.T) {
	asserts := assert.New(t)

	t.Run("And", func(t *testing.T) {
		SUT1 := And(1, 1)
		SUT2 := And(1, 2)

		asserts.Equal(Int(1), SUT1)
		asserts.Equal(Int(0), SUT2)
	})
}

func TestBitShifts(t *testing.T) {
	asserts := assert.New(t)

	t.Run("ShiftRightBy", func(t *testing.T) {
		SUT := ShiftRightBy(1, 2)

		asserts.Equal(Int(1), SUT)
	})
}
