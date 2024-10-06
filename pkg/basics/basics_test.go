package basics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasics(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Add", func(t *testing.T) {
		var SUT Int
		SUT = 1

		asserts.Equal(Int(2), Add(SUT, SUT))
	})

	t.Run("ToFloat", func(t *testing.T) {
		var SUT Int
		SUT = 23

		asserts.Equal(Float(23), ToFloat(SUT))
	})

	t.Run("Round 1.0", func(t *testing.T) {
		var SUT Float
		SUT = 1.0

		asserts.Equal(Int(1), Round(SUT))
	})

	t.Run("Round 1.2", func(t *testing.T) {
		var SUT Float
		SUT = 1.2

		asserts.Equal(Int(1), Round(SUT))
	})

	t.Run("Round 1.5", func(t *testing.T) {
		var SUT Float
		SUT = 1.5

		asserts.Equal(Int(2), Round(SUT))
	})

	t.Run("Round 1.8", func(t *testing.T) {
		var SUT Float
		SUT = 1.8

		asserts.Equal(Int(2), Round(SUT))
	})

	t.Run("Floor 1.0", func(t *testing.T) {
		var SUT Float
		SUT = 1.0

		asserts.Equal(Int(1), Floor(SUT))
	})

	t.Run("Floor 1.2", func(t *testing.T) {
		var SUT Float
		SUT = 1.2

		asserts.Equal(Int(1), Floor(SUT))
	})

	t.Run("Floor 1.5", func(t *testing.T) {
		var SUT Float
		SUT = 1.5

		asserts.Equal(Int(1), Floor(SUT))
	})

	t.Run("Floor 1.8", func(t *testing.T) {
		var SUT Float
		SUT = 1.8

		asserts.Equal(Int(1), Floor(SUT))
	})

	t.Run("Ceiling 1.0", func(t *testing.T) {
		var SUT Float
		SUT = 1.0

		asserts.Equal(Int(1), Ceiling(SUT))
	})

	t.Run("Ceiling 1.2", func(t *testing.T) {
		var SUT Float
		SUT = 1.2

		asserts.Equal(Int(2), Ceiling(SUT))
	})

	t.Run("Ceiling 1.5", func(t *testing.T) {
		var SUT Float
		SUT = 1.5

		asserts.Equal(Int(2), Ceiling(SUT))
	})

	t.Run("Ceiling 1.8", func(t *testing.T) {
		var SUT Float
		SUT = 1.8

		asserts.Equal(Int(2), Ceiling(SUT))
	})

	t.Run("Truncate 1.0", func(t *testing.T) {
		var SUT Float
		SUT = 1.0

		asserts.Equal(Int(1), Truncate(SUT))
	})

	t.Run("Truncate 1.2", func(t *testing.T) {
		var SUT Float
		SUT = 1.2

		asserts.Equal(Int(1), Truncate(SUT))
	})

	t.Run("Truncate 1.5", func(t *testing.T) {
		var SUT Float
		SUT = 1.5

		asserts.Equal(Int(1), Truncate(SUT))
	})

	t.Run("Truncate 1.8", func(t *testing.T) {
		var SUT Float
		SUT = 1.8

		asserts.Equal(Int(1), Truncate(SUT))
	})
}
