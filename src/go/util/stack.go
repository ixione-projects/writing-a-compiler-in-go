package util

type Stack[T any] struct {
	es []T
	p  int
}

const INITIAL_STACK_CAPACITY = 10

func NewStack[T any](c int) *Stack[T] {
	return &Stack[T]{
		es: make([]T, 0, max(INITIAL_STACK_CAPACITY, c)),
	}
}

func (s *Stack[T]) Push(e T) {
	if s.p >= len(s.es) {
		s.grow()
	}

	s.es[s.p] = e
	s.p += 1
}

func (s *Stack[T]) Peek() T {
	if s.p == 0 {
		return *new(T)
	}
	return s.es[s.p-1]
}

func (s *Stack[T]) Pop() T {
	if s.p == 0 {
		return *new(T)
	}
	s.p -= 1
	return s.es[s.p]
}

func (s *Stack[T]) Size() int {
	return s.p
}

func (s *Stack[T]) grow() {
	l := len(s.es)
	if l+1 >= cap(s.es) {
		s.es = append(make([]T, 0, (cap(s.es)<<1)+cap(s.es)), s.es...)
	}
	s.es = s.es[:l+1]
}
