package maybe

// Definition

type Maybe[A any] interface {
	maybe() _maybe
}

type _maybe struct{}

func (m _maybe) maybe() _maybe {
	return m
}

// Variants
type Just[V any] struct {
	_maybe
	Value V
}
type Nothing struct {
	_maybe
}

// Common helpers

/*
(a -> b) -> Maybe a -> Maybe b
*/
func Map[A any, B any](f func(A) B, m Maybe[A]) Maybe[B] {
	return MaybeWith(
		m,
		func(j *Just[A]) Maybe[B] {
			return Just[B]{Value: f(j.Value)}
		},
		func(n *Nothing) Maybe[B] { return Nothing{} },
	)
}

/*
(a -> b -> value) -> Maybe a -> Maybe b -> Maybe value
*/
func Map2[A any, B any, value any](f func(a A, b B) value, m1 Maybe[A], m2 Maybe[B]) Maybe[value] {
	return MaybeWith(
		m1,
		func(j *Just[A]) Maybe[value] {
			return MaybeWith(
				m2,
				func(j1 *Just[B]) Maybe[value] {
					return Just[value]{Value: f(j.Value, j1.Value)}
				},
				func(n *Nothing) Maybe[value] { return Nothing{} },
			)
		},
		func(n *Nothing) Maybe[value] { return Nothing{} },
	)
}

/*
Chain together many computations that may fail

(a -> Maybe b) -> Maybe a -> Maybe b
*/
func AndThen[A any, B any](f func(a A) Maybe[B], m Maybe[A]) Maybe[B] {
	return MaybeWith(
		m,
		func(j *Just[A]) Maybe[B] { return f(j.Value) },
		func(n *Nothing) Maybe[B] { return Nothing{} },
	)
}

/*
Pattern matching for the poor man
*/
func MaybeWith[V any, R any](
	m Maybe[V],
	j func(*Just[V]) R,
	n func(*Nothing) R,
) R {
	switch m := m.(type) {
	case Just[V]:
		return j(&m)
	case Nothing:
		return n(&m)
	}
	panic("unreachable")
}
