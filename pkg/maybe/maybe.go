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
type Just[A any] struct {
	_maybe
	Value A
}

type Nothing struct {
	_maybe
}

// Common helpers

func WithDefault[A any](a A, m Maybe[A]) A {
	return MaybeWith(
		m,
		func(j Just[A]) A { return j.Value },
		func(_ Nothing) A { return a },
	)
}

// Transform a Maybe value with a given function:
//
// (a -> b) -> Maybe a -> Maybe b
func Map[A any, B any](f func(A) B, m Maybe[A]) Maybe[B] {
	return MaybeWith(
		m,
		func(j Just[A]) Maybe[B] {
			return Just[B]{Value: f(j.Value)}
		},
		func(_ Nothing) Maybe[B] { return Nothing{} },
	)
}

// Apply a function if all the arguments are Just a value.
//
// (a -> b -> value) -> Maybe a -> Maybe b -> Maybe value
func Map2[A any, B any, value any](f func(a A, b B) value, m1 Maybe[A], m2 Maybe[B]) Maybe[value] {
	return MaybeWith(
		m1,
		func(j Just[A]) Maybe[value] {
			return MaybeWith(
				m2,
				func(j1 Just[B]) Maybe[value] {
					return Just[value]{Value: f(j.Value, j1.Value)}
				},
				func(n Nothing) Maybe[value] { return Nothing{} },
			)
		},
		func(n Nothing) Maybe[value] { return Nothing{} },
	)
}

// (a -> b -> c -> value) -> Maybe a -> Maybe b -> Maybe c -> Maybe value
func Map3[A any, B any, C any, value any](f func(a A, b B, c C) value, m1 Maybe[A], m2 Maybe[B], m3 Maybe[C]) Maybe[value] {
	return MaybeWith(
		m1,
		func(j Just[A]) Maybe[value] {
			return MaybeWith(
				m2,
				func(j1 Just[B]) Maybe[value] {
					return MaybeWith(
						m3,
						func(j2 Just[C]) Maybe[value] { return Just[value]{Value: f(j.Value, j1.Value, j2.Value)} },
						func(n Nothing) Maybe[value] { return Nothing{} },
					)
				},
				func(n Nothing) Maybe[value] { return Nothing{} },
			)
		},
		func(n Nothing) Maybe[value] { return Nothing{} },
	)
}

// (a -> b -> c -> d -> value) -> Maybe a -> Maybe b -> Maybe c -> Maybe d -> Maybe value
func Map4[A any, B any, C any, D any, value any](f func(a A, b B, c C, d D) value, m1 Maybe[A], m2 Maybe[B], m3 Maybe[C], m4 Maybe[D]) Maybe[value] {
	return MaybeWith(
		m1,
		func(j Just[A]) Maybe[value] {
			return MaybeWith(
				m2,
				func(j1 Just[B]) Maybe[value] {
					return MaybeWith(
						m3,
						func(j2 Just[C]) Maybe[value] {
							return MaybeWith(
								m4,
								func(j3 Just[D]) Maybe[value] { return Just[value]{Value: f(j.Value, j1.Value, j2.Value, j3.Value)} },
								func(n Nothing) Maybe[value] { return Nothing{} },
							)
						},
						func(n Nothing) Maybe[value] { return Nothing{} },
					)
				},
				func(n Nothing) Maybe[value] { return Nothing{} },
			)
		},
		func(n Nothing) Maybe[value] { return Nothing{} },
	)
}

// (a -> b -> c -> d -> e -> value) -> Maybe a -> Maybe b -> Maybe c -> Maybe d -> Maybe e -> Maybe value
func Map5[A any, B any, C any, D any, E any, value any](f func(a A, b B, c C, d D, e E) value, m1 Maybe[A], m2 Maybe[B], m3 Maybe[C], m4 Maybe[D], m5 Maybe[E]) Maybe[value] {
	return MaybeWith(
		m1,
		func(j Just[A]) Maybe[value] {
			return MaybeWith(
				m2,
				func(j1 Just[B]) Maybe[value] {
					return MaybeWith(
						m3,
						func(j2 Just[C]) Maybe[value] {
							return MaybeWith(
								m4,
								func(j3 Just[D]) Maybe[value] {
									return MaybeWith(
										m5,
										func(j4 Just[E]) Maybe[value] {
											return Just[value]{Value: f(j.Value, j1.Value, j2.Value, j3.Value, j4.Value)}
										},
										func(n Nothing) Maybe[value] { return Nothing{} },
									)
								},
								func(n Nothing) Maybe[value] { return Nothing{} },
							)
						},
						func(n Nothing) Maybe[value] { return Nothing{} },
					)
				},
				func(n Nothing) Maybe[value] { return Nothing{} },
			)
		},
		func(n Nothing) Maybe[value] { return Nothing{} },
	)
}

// Chain together many computations that may fail
//
//	(a -> Maybe b) -> Maybe a -> Maybe b
func AndThen[A any, B any](f func(A) Maybe[B], m Maybe[A]) Maybe[B] {
	return MaybeWith(
		m,
		func(j Just[A]) Maybe[B] { return f(j.Value) },
		func(n Nothing) Maybe[B] { return Nothing{} },
	)
}

// Pattern matching for the poor man
func MaybeWith[V any, R any](
	m Maybe[V],
	j func(Just[V]) R,
	n func(Nothing) R,
) R {
	switch m := m.(type) {
	case Just[V]:
		return j(m)
	case Nothing:
		return n(m)
	}
	panic("unreachable")
}
