package iterator

type Iterator[T any] struct {
	data  []T
	index int
}

func FromSlice[T any](s []T) *Iterator[T] {
	return &Iterator[T]{data: s}
}

func (it *Iterator[T]) Next() (T, bool) {
	if it.index >= len(it.data) {
		var zero T
		return zero, false
	}
	v := it.data[it.index]
	it.index++
	return v, true
}

type Query[T any] struct {
	src        func() (T, bool)
	predicates []func(T) bool
	selectors  []func(any) any
}

func NewQuery[T any](src *Iterator[T]) *Query[T] {
	return &Query[T]{
		src: func() (T, bool) { return src.Next() },
	}
}

func Where[T any](q *Query[T], pred func(T) bool) *Query[T] {
	q.predicates = append(q.predicates, pred)
	return q
}

func Select[T any, R any](q *Query[T], sel func(T) R) *Query[R] {
	return &Query[R]{
		src: func() (R, bool) {
			for {
				v, ok := q.Next()
				if !ok {
					var zero R
					return zero, false
				}
				return sel(v), true
			}
		},
	}
}

func (q *Query[T]) Next() (T, bool) {
	for {
		v, ok := q.src()
		if !ok {
			var zero T
			return zero, false
		}

		okAll := true
		for _, pred := range q.predicates {
			if !pred(v) {
				okAll = false
				break
			}
		}
		if !okAll {
			continue
		}

		if len(q.selectors) == 0 {
			return v, true
		}

		var r any = v
		for _, sel := range q.selectors {
			r = sel(r)
		}

		casted, _ := r.(T)
		return casted, true
	}
}
