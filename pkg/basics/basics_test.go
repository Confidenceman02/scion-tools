package basics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasics(t *testing.T) {
	asserts := assert.New(t)

	t.Run("ToFloat", func(t *testing.T) {
		var SUT int
		SUT = 23

		asserts.Equal(23.0, ToFloat(SUT))
	})

	t.Run("Round 1.0", func(t *testing.T) {
		var SUT float64
		SUT = 1.0

		asserts.Equal(1, Round(SUT))
	})

	t.Run("Round 1.2", func(t *testing.T) {
		var SUT float64
		SUT = 1.2

		asserts.Equal(1, Round(SUT))
	})

	t.Run("Round 1.5", func(t *testing.T) {
		var SUT float64
		SUT = 1.5

		asserts.Equal(2, Round(SUT))
	})

	t.Run("Round 1.8", func(t *testing.T) {
		var SUT float64
		SUT = 1.8

		asserts.Equal(2, Round(SUT))
	})

	t.Run("Floor 1.0", func(t *testing.T) {
		var SUT float64
		SUT = 1.0

		asserts.Equal(1, Floor(SUT))
	})

	t.Run("Floor 1.2", func(t *testing.T) {
		var SUT float64
		SUT = 1.2

		asserts.Equal(1, Floor(SUT))
	})

	t.Run("Floor 1.5", func(t *testing.T) {
		var SUT float64
		SUT = 1.5

		asserts.Equal(1, Floor(SUT))
	})

	t.Run("Floor 1.8", func(t *testing.T) {
		var SUT float64
		SUT = 1.8

		asserts.Equal(1, Floor(SUT))
	})

	t.Run("Ceiling 1.0", func(t *testing.T) {
		var SUT float64
		SUT = 1.0

		asserts.Equal(1, Ceiling(SUT))
	})

	t.Run("Ceiling 1.2", func(t *testing.T) {
		var SUT float64
		SUT = 1.2

		asserts.Equal(2, Ceiling(SUT))
	})

	t.Run("Ceiling 1.5", func(t *testing.T) {
		var SUT float64
		SUT = 1.5

		asserts.Equal(2, Ceiling(SUT))
	})

	t.Run("Ceiling 1.8", func(t *testing.T) {
		var SUT float64
		SUT = 1.8

		asserts.Equal(2, Ceiling(SUT))
	})

	t.Run("Truncate 1.0", func(t *testing.T) {
		var SUT float64
		SUT = 1.0

		asserts.Equal(1, Truncate(SUT))
	})

	t.Run("Truncate 1.2", func(t *testing.T) {
		var SUT float64
		SUT = 1.2

		asserts.Equal(1, Truncate(SUT))
	})

	t.Run("Truncate 1.5", func(t *testing.T) {
		var SUT float64
		SUT = 1.5

		asserts.Equal(1, Truncate(SUT))
	})

	t.Run("Truncate 1.8", func(t *testing.T) {
		var SUT float64
		SUT = 1.8

		asserts.Equal(1, Truncate(SUT))
	})
}
