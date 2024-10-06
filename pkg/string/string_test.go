package string

import (
	"testing"

	. "github.com/Confidenceman02/scion-tools/pkg/basics"
	. "github.com/Confidenceman02/scion-tools/pkg/maybe"
	"github.com/stretchr/testify/assert"
)

func TestIntConversions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("ToInt", func(t *testing.T) {
		i1 := String("123")
		i2 := String("-42")
		i3 := String("3.1")
		i4 := String("31a")

		SUT1 := ToInt(i1)
		SUT2 := ToInt(i2)
		SUT3 := ToInt(i3)
		SUT4 := ToInt(i4)

		asserts.Equal(Just[Int]{Value: Int(123)}, SUT1)
		asserts.Equal(Just[Int]{Value: Int(-42)}, SUT2)
		asserts.Equal(Nothing{}, SUT3)
		asserts.Equal(Nothing{}, SUT4)
	})
}

func TestBasicsComparisons(t *testing.T) {
	asserts := assert.New(t)

	asserts.Equal(String("xyz"), Max(String("abc"), String("xyz")))
	asserts.Equal(String("abc"), Min(String("abc"), String("xyz")))
}
