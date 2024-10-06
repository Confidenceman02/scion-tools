package tuple

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTuple2(t *testing.T) {
	asserts := assert.New(t)
	t.Run("Pair", func(t *testing.T) {
		asserts.Equal(&tuple2[int, int]{&_tuple2[int, int]{1, 2}}, Pair(1, 2))
	})
}
