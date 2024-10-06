package list

import (
	"github.com/Confidenceman02/scion-tools/pkg/basics"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmpty(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Empty", func(t *testing.T) {
		SUT := Empty[basics.Int]()

		asserts.Nil(SUT._list()._cons)
	})
}

func TestSingleton(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Singleton", func(t *testing.T) {
		ls := Singleton[basics.Int](10)

		SUT := ls._list()
		asserts.Equal(&list[basics.Int]{&cons[basics.Int]{10, nil}}, SUT)
	})
}

func TestRepeat(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Repeat", func(t *testing.T) {
		SUT := Repeat[basics.Int](2, 10)

		asserts.Nil(SUT._list()._cons.tail._list()._cons.tail._cons)
		asserts.Equal(basics.Int(10), SUT._list()._cons.head)
		asserts.Equal(basics.Int(10), SUT._list()._cons.tail._cons.head)
	})
}

func TestCons(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Cons", func(t *testing.T) {
		ls := Singleton(10)
		SUT := Cons(20, ls)

		// Head
		asserts.Equal(20, SUT._list()._cons.head)
		// Tail
		asserts.Equal(ls._list(), SUT._list()._cons.tail)
	})
}
