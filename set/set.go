package set

import (
	"fmt"
	"sort"

	"github.com/qinyuguang/gods/internal/rwmutex"
)

// Set is a unordered collection of unique members.
type Set[T comparable] struct {
	mu   rwmutex.RWMutex
	data map[T]struct{}
}

// New returns an empty set.
// The parameter `safe` is used to specify whether using set in concurrent-safety,
// which is false in default.
func New[T comparable](safe ...bool) *Set[T] {
	return &Set[T]{
		mu:   rwmutex.Create(safe...),
		data: make(map[T]struct{}),
	}
}

// NewFrom returns a set from `items`.
func NewFrom[T comparable](items []T, safe ...bool) *Set[T] {
	data := make(map[T]struct{})
	for i := range items {
		data[items[i]] = struct{}{}
	}
	return &Set[T]{
		mu:   rwmutex.Create(safe...),
		data: data,
	}
}

// Each calls 'fn' on every item in the set in no particular order,
// if `fn` returns true then continue iterating; or false to stop.
func (s *Set[T]) Each(fn func(item T) bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for v := range s.data {
		if !fn(v) {
			break
		}
	}
}

// Add adds one or multiple items to the set.
func (s *Set[T]) Add(items ...T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.data == nil {
		s.data = make(map[T]struct{})
	}
	for i := range items {
		s.data[items[i]] = struct{}{}
	}
}

// Remove deletes one or multiple items from set.
func (s *Set[T]) Remove(items ...T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.data != nil {
		for i := range items {
			delete(s.data, items[i])
		}
	}
}

// Contains checks whether the set contains `item`.
func (s *Set[T]) Contains(item T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.data[item]
	return ok
}

// Size returns the number of items in the set.
func (s *Set[T]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

// Clear deletes all items of the set.
func (s *Set[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[T]struct{})
}

// Slice returns all items of the set as slice.
func (s *Set[T]) Slice() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	slice := make([]T, 0, len(s.data))
	for k := range s.data {
		slice = append(slice, k)
	}
	return slice
}

// String returns items as a string.
func (s *Set[T]) String() string {
	out := make([]string, 0, s.Size())
	s.Each(func(v T) bool { out = append(out, fmt.Sprintf(`%v`, v)); return true })
	sort.Strings(out)
	return fmt.Sprintf("%v", out)
}

// Clone returns a new set by deep copy.
func (s *Set[T]) Clone() *Set[T] {
	return NewFrom(s.Slice(), s.mu.IsSafe())
}

// Equal checks whether the two sets equal.
func (s *Set[T]) Equal(other *Set[T]) bool {
	if other == nil {
		return false
	}
	if s == other {
		return true
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()
	if len(s.data) != len(other.data) {
		return false
	}
	for key := range s.data {
		if _, ok := other.data[key]; !ok {
			return false
		}
	}
	return true
}

// IsSubsetOf checks whether the current set is a sub-set of `other`.
func (s *Set[T]) IsSubsetOf(other *Set[T]) bool {
	if other == nil {
		return false
	}
	if s == other {
		return true
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()
	for key := range s.data {
		if _, ok := other.data[key]; !ok {
			return false
		}
	}
	return true
}

// Union returns a new set which is the union of `set` and `other`.
// Which means, all the items in `newSet` are in `set` or in `other`.
func (s *Set[T]) Union(others ...*Set[T]) *Set[T] {
	newSet := s.Clone()
	for _, other := range others {
		if other == nil {
			continue
		}
		other.mu.RLock()
		for k, v := range other.data {
			newSet.data[k] = v
		}
		other.mu.RUnlock()
	}
	return newSet
}

// Diff returns a new set which is the difference set from `set` to `other`.
// Which means, all the items in `newSet` are in `set` but not in `other`.
func (s *Set[T]) Diff(others ...*Set[T]) *Set[T] {
	newSet := s.Clone()
	for _, other := range others {
		if other == nil {
			continue
		}
		other.mu.RLock()
		for k := range other.data {
			delete(newSet.data, k)
		}
		other.mu.RUnlock()
	}
	return newSet
}

// Intersect returns a new set which is the intersection from `set` to `other`.
// Which means, all the items in `newSet` are in `set` and also in `other`.
func (s *Set[T]) Intersect(others ...*Set[T]) *Set[T] {
	newSet := s.Clone()
	for _, other := range others {
		if other == nil {
			return New[T](s.mu.IsSafe())
		}
		other.mu.RLock()
		for k := range newSet.data {
			if _, ok := other.data[k]; !ok {
				delete(newSet.data, k)
			}
		}
		other.mu.RUnlock()
	}
	return newSet
}
