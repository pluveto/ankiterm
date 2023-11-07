package typing

type Option[T any] struct {
	Value T
	IsSome bool
}

func Some[T any](v T) Option[T] {
	return Option[T]{
		Value: v,
		IsSome: true,
	}
}

func None[T any]() Option[T] {
	return Option[T]{
		IsSome: false,
	}
}

func (o Option[T]) Unwrap() T {
	return o.Value
}

func (o Option[T]) UnwrapOr(def T) T {
	if o.IsSome {
		return o.Value
	} else {
		return def
	}
}

func (o Option[T]) UnwrapOrElse(f func() T) T {
	if o.IsSome {
		return o.Value
	} else {
		return f()
	}
}

func (o Option[T]) IsNone() bool {
	return !o.IsSome
}
