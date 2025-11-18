package iterator

type Iterator struct {
	data      []any
	index     int
	selector  []func(any) any
	predicate []func(any) bool
}

func NewIterator[T any](items []T) *Iterator {
	boxed := make([]any, len(items))
	for i, item := range items {
		boxed[i] = item
	}

	return &Iterator{
		data:      boxed,
		index:     0,
		selector:  []func(v any) any{},
		predicate: []func(v any) bool{},
	}
}

func (it *Iterator) matchPredicates(v any) bool {
	for _, pred := range it.predicate {
		if !pred(v) {
			return false
		}
	}
	return true
}

func (it *Iterator) HasNext() bool {
	return it.index < len(it.data)
}

func (it *Iterator) Next() any {
	if !it.HasNext() {
		return nil
	}
	v := it.data[it.index]
	it.index++
	return it.selector(v)
}

// Terminal Op
func (it *Iterator) ToList() []any {
	var result []any
	for it.HasNext() {
		v := it.data[it.index]
		if it.predicate(v) {
			result = append(result, it.Next())
		} else {
			it.index++
		}
	}
	return result
}

// Intermediate Op
func (it *Iterator) Select(fn func(v any) any) *Iterator {
	it.selector = fn
	return it
}

func (it *Iterator) Where(fn func(v any) bool) *Iterator {
	it.predicate = fn
	return it
}
