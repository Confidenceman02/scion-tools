package tuple

type Tuple2[A, B any] interface {
	tuple2() *_tuple2[A, B]
}

type _tuple2[A, B any] struct {
	a A
	b B
}

func (t *_tuple2[A, B]) tuple2() *_tuple2[A, B] {
	return t
}

type tuple2[A, B any] struct {
	*_tuple2[A, B]
}

// Create

// Create a 2-tuple.
func Pair[A, B any](a A, b B) Tuple2[A, B] {
	return &tuple2[A, B]{&_tuple2[A, B]{a: a, b: b}}
}

// Access

// Extract the first value from a tuple.
func First[A, B any](t Tuple2[A, B]) A {
	return t.tuple2().a
}

// Extract the second value from a tuple.
func Second[A, B any](t Tuple2[A, B]) B {
	return t.tuple2().b
}
