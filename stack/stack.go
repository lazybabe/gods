package stack

import (
	"github.com/lazybabe/gods/array"
)

type Stack[T comparable] struct {
	data *array.Array[T]
}

// New creates and returns an empty stack.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func New[T comparable](safe ...bool) *Stack[T] {
	return &Stack[T]{
		data: array.New[T](safe...),
	}
}

// NewFrom creates and returns a stack, and push the elements of `data` at the top of the stack one by one.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewFrom[T comparable](data []T, safe ...bool) *Stack[T] {
	return &Stack[T]{
		data: array.NewFrom(data, safe...),
	}
}

// Push places 'value' at the top of the stack.
func (s *Stack[T]) Push(value T) {
	s.data.PushRight(value)
}

// Pop removes the stack's top element and returns it.
// If the stack is empty it returns the zero value.
func (s *Stack[T]) Pop() T {
	value, _ := s.data.PopRight()
	return value
}

// Peek returns the stack's top element but does not remove it.
// If the stack is empty the zero value is returned.
func (s *Stack[T]) Peek() (t T) {
	return s.data.Index(s.data.Size() - 1)
}

// Size returns the number of elements in the stack.
func (s *Stack[T]) Size() int {
	return s.data.Size()
}

// Clone returns a new stack, which is a copy of current stack.
func (s *Stack[T]) Clone() *Stack[T] {
	return &Stack[T]{
		data: s.data.Clone(),
	}
}

// IsEmpty returns true if the stack is empty, otherwise returns false.
func (s *Stack[T]) IsEmpty() bool {
	return s.data.Size() == 0
}
