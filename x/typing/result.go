package typing

type Result[T any, E any] struct {
	err   *E
	value *T
}

func Ok[T any, E any](value T) Result[T, E] {
	return Result[T,E]{
		value: &value,
		err:   nil,
	}
}

func Err[T any, E any](err E) Result[T,E] {
	return Result[T, E]{
		value: nil,
		err:   &err,
	}
}

func (r Result[T,E]) Unwrap() T {
	if r.err != nil {
		panic(*r.err)
	}
	return *r.value
}


