// A Result is either Ok meaning the computation succeeded, or it is an Err meaning that there was some failure.
package result

import (
	"fmt"
	"github.com/Confidenceman02/scion-tools/pkg/maybe"
	"reflect"
)

type Result[E, V any] interface {
	result() _result[E, V]
}

type _result[E, V any] struct{}

func (r _result[E, V]) result() _result[E, V] {
	return r
}

// VARIANTS

type Err[E, V any] struct {
	_result[E, V]
	Err E
}

type Ok[E, V any] struct {
	_result[E, V]
	Val V
}

// Mapping

// Apply a function to a result. If the result is Ok, it will be converted. If the result is an Err, the same error value will propagate through.
func Map[X, A, V any](f func(A) V, ra Result[X, A]) Result[X, V] {
	return ResultWith(
		ra,
		func(e Err[X, A]) Result[X, V] {
			return Err[X, V]{Err: e.Err}
		},
		func(o Ok[X, A]) Result[X, V] {
			return Ok[X, V]{Val: f(o.Val)}
		},
	)
}

// Apply a function if both results are Ok. If not, the first Err will propagate through.
func Map2[X, A, B, value any](f func(A, B) value, ra Result[X, A], rb Result[X, B]) Result[X, value] {
	return ResultWith(
		ra,
		func(e Err[X, A]) Result[X, value] {
			return Err[X, value]{Err: e.Err}
		},
		func(o Ok[X, A]) Result[X, value] {
			return ResultWith(
				rb,
				func(e1 Err[X, B]) Result[X, value] {
					return Err[X, value]{Err: e1.Err}
				},
				func(o1 Ok[X, B]) Result[X, value] {
					return Ok[X, value]{Val: f(o.Val, o1.Val)}
				},
			)
		},
	)
}

func Map3[X, A, B, C, value any](
	f func(A, B, C) value,
	ra Result[X, A],
	rb Result[X, B],
	rc Result[X, C],
) Result[X, value] {
	return ResultWith(
		ra,
		func(e Err[X, A]) Result[X, value] { return Err[X, value]{Err: e.Err} },
		func(ok1 Ok[X, A]) Result[X, value] {
			return ResultWith(
				rb,
				func(e1 Err[X, B]) Result[X, value] { return Err[X, value]{Err: e1.Err} },
				func(ok2 Ok[X, B]) Result[X, value] {
					return ResultWith(
						rc,
						func(e2 Err[X, C]) Result[X, value] { return Err[X, value]{Err: e2.Err} },
						func(ok3 Ok[X, C]) Result[X, value] {
							return Ok[X, value]{Val: f(ok1.Val, ok2.Val, ok3.Val)}
						},
					)
				},
			)
		},
	)
}

func Map4[X, A, B, C, D, value any](
	f func(A, B, C, D) value,
	ra Result[X, A],
	rb Result[X, B],
	rc Result[X, C],
	rd Result[X, D],
) Result[X, value] {
	return ResultWith(
		ra,
		func(e Err[X, A]) Result[X, value] { return Err[X, value]{Err: e.Err} },
		func(ok1 Ok[X, A]) Result[X, value] {
			return ResultWith(
				rb,
				func(e1 Err[X, B]) Result[X, value] { return Err[X, value]{Err: e1.Err} },
				func(ok2 Ok[X, B]) Result[X, value] {
					return ResultWith(
						rc,
						func(e2 Err[X, C]) Result[X, value] { return Err[X, value]{Err: e2.Err} },
						func(ok3 Ok[X, C]) Result[X, value] {
							return ResultWith(
								rd,
								func(e4 Err[X, D]) Result[X, value] { return Err[X, value]{Err: e4.Err} },
								func(ok4 Ok[X, D]) Result[X, value] {
									return Ok[X, value]{Val: f(ok1.Val, ok2.Val, ok3.Val, ok4.Val)}
								},
							)
						},
					)
				},
			)
		},
	)
}
func Map5[X, A, B, C, D, E, value any](
	f func(A, B, C, D, E) value,
	ra Result[X, A],
	rb Result[X, B],
	rc Result[X, C],
	rd Result[X, D],
	re Result[X, E],
) Result[X, value] {
	return ResultWith(
		ra,
		func(e Err[X, A]) Result[X, value] { return Err[X, value]{Err: e.Err} },
		func(ok1 Ok[X, A]) Result[X, value] {
			return ResultWith(
				rb,
				func(e1 Err[X, B]) Result[X, value] { return Err[X, value]{Err: e1.Err} },
				func(ok2 Ok[X, B]) Result[X, value] {
					return ResultWith(
						rc,
						func(e2 Err[X, C]) Result[X, value] { return Err[X, value]{Err: e2.Err} },
						func(ok3 Ok[X, C]) Result[X, value] {
							return ResultWith(
								rd,
								func(e4 Err[X, D]) Result[X, value] { return Err[X, value]{Err: e4.Err} },
								func(ok4 Ok[X, D]) Result[X, value] {
									return ResultWith(
										re,
										func(e5 Err[X, E]) Result[X, value] { return Err[X, value]{Err: e5.Err} },
										func(ok5 Ok[X, E]) Result[X, value] {
											return Ok[X, value]{Val: f(ok1.Val, ok2.Val, ok3.Val, ok4.Val, ok5.Val)}
										},
									)
								},
							)
						},
					)
				},
			)
		},
	)
}

// Chaining

// Chain together a sequence of computations that may fail.
func AndThen[X, A, B any](f func(A) Result[X, B], r Result[X, A]) Result[X, B] {
	return ResultWith(
		r,
		func(e Err[X, A]) Result[X, B] { return Err[X, B]{Err: e.Err} },
		func(o Ok[X, A]) Result[X, B] { return f(o.Val) },
	)
}

// If the result is Ok return the value, but if the result is an Err then return a given default value.
func WithDefault[E, V any](r Result[E, V], defaultValue V) V {
	return ResultWith(
		r,
		func(e Err[E, V]) V { return defaultValue },
		func(o Ok[E, V]) V { return o.Val },
	)
}

// Convert to a simpler Maybe if the actual error message is not needed or you need to interact with some code that primarily uses maybes.
func ToMaybe[E, V any](r Result[E, V]) maybe.Maybe[V] {
	return ResultWith(
		r,
		func(e Err[E, V]) maybe.Maybe[V] { return maybe.Nothing{} },
		func(o Ok[E, V]) maybe.Maybe[V] { return maybe.Just[V]{Value: o.Val} },
	)
}

// Convert from a simple Maybe to interact with some code that primarily uses Results.
func FromMaybe[X, V any](e X, m maybe.Maybe[V]) Result[X, V] {
	return maybe.MaybeWith(
		m,
		func(j maybe.Just[V]) Result[X, V] { return Ok[X, V]{Val: j.Value} },
		func(maybe.Nothing) Result[X, V] { return Err[X, V]{Err: e} },
	)
}

// Transform an Err value. For example, say the errors we get have too much information:
func MapError[X, Y, V any](f func(X) Y, r Result[X, V]) Result[Y, V] {
	return ResultWith(
		r,
		func(e Err[X, V]) Result[Y, V] { return Err[Y, V]{Err: f(e.Err)} },
		func(o Ok[X, V]) Result[Y, V] { return Ok[Y, V]{Val: o.Val} },
	)
}

// Utilities

func ResultWith[E, V, R any](
	r Result[E, V],
	err func(Err[E, V]) R,
	ok func(Ok[E, V]) R) R {
	switch r := r.(type) {
	case Err[E, V]:
		return err(r)
	case Ok[E, V]:
		return ok(r)
	default:
		var zero [0]V
		panic(
			fmt.Sprintf(
				"\nI was expecting a type of: \n    maybe.Maybe[%v]\n\nBut instead got a\n    %v\n",
				reflect.TypeOf(zero).Elem(),
				reflect.TypeOf(r),
			),
		)
	}
}
