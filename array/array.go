package array

import (
	"fmt"
	"math"
	"sort"

	"github.com/lazybabe/gods/internal/rwmutex"
)

type Array[T comparable] struct {
	mu    rwmutex.RWMutex
	array []T
}

// New creates and returns an empty array.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func New[T comparable](safe ...bool) *Array[T] {
	return NewSize[T](0, 0, safe...)
}

// NewSize create and returns an array with given size and cap.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewSize[T comparable](size int, cap int, safe ...bool) *Array[T] {
	return &Array[T]{
		mu:    rwmutex.Create(safe...),
		array: make([]T, size, cap),
	}
}

// NewFrom creates and returns an array with given slice `array`.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewFrom[T comparable](array []T, safe ...bool) *Array[T] {
	return &Array[T]{
		mu:    rwmutex.Create(safe...),
		array: array,
	}
}

// Index returns the value by the specified index.
// If the given `index` is out of range of the array, it returns `nil`.
func (a *Array[T]) Index(index int) T {
	value, _ := a.Get(index)
	return value
}

// Get returns the value by the specified index.
// If the given `index` is out of range of the array, the `found` is false.
func (a *Array[T]) Get(index int) (value T, found bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if index < 0 || index >= len(a.array) {
		return
	}
	return a.array[index], true
}

// Set sets value to specified index.
func (a *Array[T]) Set(index int, value T) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if index < 0 || index >= len(a.array) {
		return fmt.Errorf("index %d out of array range %d", index, len(a.array))
	}
	a.array[index] = value
	return nil
}

// SortFunc sorts the array by custom function `less`.
func (a *Array[T]) Sort(less func(v1, v2 T) bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	sort.Slice(a.array, func(i, j int) bool {
		return less(a.array[i], a.array[j])
	})
}

// InsertBefore inserts the `value` to the front of `index`.
func (a *Array[T]) InsertBefore(index int, value T) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if index < 0 || index >= len(a.array) {
		return fmt.Errorf("index %d out of array range %d", index, len(a.array))
	}
	rear := append([]T{}, a.array[index:]...)
	a.array = append(a.array[0:index], value)
	a.array = append(a.array, rear...)
	return nil
}

// InsertAfter inserts the `value` to the back of `index`.
func (a *Array[T]) InsertAfter(index int, value T) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if index < 0 || index >= len(a.array) {
		return fmt.Errorf("index %d out of array range %d", index, len(a.array))
	}
	rear := append([]T{}, a.array[index+1:]...)
	a.array = append(a.array[0:index+1], value)
	a.array = append(a.array, rear...)
	return nil
}

// Remove removes an item by index.
// If the given `index` is out of range of the array, the `found` is false.
func (a *Array[T]) Remove(index int) (value T, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.doRemoveWithoutLock(index)
}

// doRemoveWithoutLock removes an item by index without lock.
func (a *Array[T]) doRemoveWithoutLock(index int) (value T, found bool) {
	if index < 0 || index >= len(a.array) {
		return value, false
	}
	// Determine array boundaries when deleting to improve deletion efficiency.
	if index == 0 {
		value := a.array[0]
		a.array = a.array[1:]
		return value, true
	} else if index == len(a.array)-1 {
		value := a.array[index]
		a.array = a.array[:index]
		return value, true
	}
	// If it is a non-boundary delete,
	// it will involve the creation of an array,
	// then the deletion is less efficient.
	value = a.array[index]
	a.array = append(a.array[:index], a.array[index+1:]...)
	return value, true
}

// RemoveValue removes an item by value.
// It returns true if value is found in the array, or else false if not found.
func (a *Array[T]) RemoveValue(value T) bool {
	if i := a.Search(value); i != -1 {
		a.Remove(i)
		return true
	}
	return false
}

// PushLeft pushes one or multiple items to the beginning of array.
func (a *Array[T]) PushLeft(value ...T) *Array[T] {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.array = append(value, a.array...)
	return a
}

// PushRight pushes one or multiple items to the end of array.
// It equals to Append.
func (a *Array[T]) PushRight(value ...T) *Array[T] {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.array = append(a.array, value...)
	return a
}

// PopLeft pops and returns an item from the beginning of array.
// Note that if the array is empty, the `found` is false.
func (a *Array[T]) PopLeft() (value T, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.array) == 0 {
		return value, false
	}
	value = a.array[0]
	a.array = a.array[1:]
	return value, true
}

// PopRight pops and returns an item from the end of array.
// Note that if the array is empty, the `found` is false.
func (a *Array[T]) PopRight() (value T, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	index := len(a.array) - 1
	if index < 0 {
		return value, false
	}
	value = a.array[index]
	a.array = a.array[:index]
	return value, true
}

// SubSlice returns a slice of elements from the array as specified
// by the `offset` and `size` parameters.
// If in concurrent safe usage, it returns a copy of the slice; else a pointer.
//
// If offset is non-negative, the sequence will start at that offset in the array.
// If offset is negative, the sequence will start that far from the end of the array.
//
// If length is given and is positive, then the sequence will have up to that many elements in it.
// If the array is shorter than the length, then only the available array elements will be present.
// If length is given and is negative then the sequence will stop that many elements from the end of the array.
// If it is omitted, then the sequence will have everything from offset up until the end of the array.
//
// Any possibility crossing the left border of array, it will fail.
func (a *Array[T]) SubSlice(offset int, length ...int) []T {
	a.mu.RLock()
	defer a.mu.RUnlock()
	size := len(a.array)
	if len(length) > 0 {
		size = length[0]
	}
	if offset > len(a.array) {
		return nil
	}
	if offset < 0 {
		offset = len(a.array) + offset
		if offset < 0 {
			return nil
		}
	}
	if size < 0 {
		offset += size
		size = -size
		if offset < 0 {
			return nil
		}
	}
	if offset+size > len(a.array) {
		size = len(a.array) - offset
	}
	s := make([]T, size)
	copy(s, a.array[offset:])
	return s
}

// Append is alias of PushRight, please See PushRight.
func (a *Array[T]) Append(value ...T) *Array[T] {
	a.PushRight(value...)
	return a
}

// Size returns the length of array.
func (a *Array[T]) Size() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.array)
}

// Slice returns the underlying data of array.
// Note that, if it's in concurrent-safe usage, it returns a copy of underlying data,
// or else a pointer to the underlying data.
func (a *Array[T]) Slice() []T {
	a.mu.RLock()
	defer a.mu.RUnlock()
	array := make([]T, len(a.array))
	copy(array, a.array)
	return array
}

// Clone returns a new array, which is a copy of current array.
func (a *Array[T]) Clone() (newArray *Array[T]) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	array := make([]T, len(a.array))
	copy(array, a.array)
	return NewFrom(array, a.mu.IsSafe())
}

// Clear deletes all items of current array.
func (a *Array[T]) Clear() *Array[T] {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.array) > 0 {
		a.array = make([]T, 0)
	}
	return a
}

// Contains checks whether a value exists in the array.
func (a *Array[T]) Contains(value T) bool {
	return a.Search(value) != -1
}

// Search searches array by `value`, returns the index of `value`,
// or returns -1 if not exists.
func (a *Array[T]) Search(value T) int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	result := -1
	for index, v := range a.array {
		if v == value {
			result = index
			break
		}
	}
	return result
}

// Unique uniques the array, clear repeated items.
// Example: [2, 3, 1, 2, 1, 4] -> [2, 3, 1, 4]
func (a *Array[T]) Unique() *Array[T] {
	a.mu.Lock()
	defer a.mu.Unlock()
	result := make([]T, 0, len(a.array))
	seen := make(map[T]struct{}, len(a.array))
	for i := 0; i < len(a.array); i++ {
		item := a.array[i]
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}
	a.array = result
	return a
}

// Fill fills an array with num entries of the value `value`,
// keys starting at the `startIndex` parameter.
func (a *Array[T]) Fill(startIndex int, num int, value T) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if startIndex < 0 || startIndex > len(a.array) {
		return fmt.Errorf("index %d out of array range %d", startIndex, len(a.array))
	}
	for i := startIndex; i < startIndex+num; i++ {
		if i > len(a.array)-1 {
			a.array = append(a.array, value)
		} else {
			a.array[i] = value
		}
	}
	return nil
}

// Chunk returns an array of elements split into groups the length of size.
// If array can't be split evenly, the final chunk will be the remaining elements.
func (a *Array[T]) Chunk(size int) [][]T {
	if size < 1 {
		return nil
	}
	a.mu.RLock()
	defer a.mu.RUnlock()
	length := len(a.array)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	var result [][]T
	for i, end := 0, 0; chunks > 0; chunks-- {
		end = (i + 1) * size
		if end > length {
			end = length
		}
		result = append(result, a.array[i*size:end])
		i++
	}
	return result
}

// Reverse makes array with elements in reverse order.
func (a *Array[T]) Reverse() *Array[T] {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i, j := 0, len(a.array)-1; i < j; i, j = i+1, j-1 {
		a.array[i], a.array[j] = a.array[j], a.array[i]
	}
	return a
}

// Each calls 'fn' on every item in the array in ascending order.
// If `f` returns true, then it continues iterating; or false to stop.
func (a *Array[T]) Each(f func(k int, v T) bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for k, v := range a.array {
		if !f(k, v) {
			break
		}
	}
}

// String returns current array as a string, which implements like json.Marshal does.
func (a *Array[T]) String() string {
	out := make([]string, 0, a.Size())
	a.Each(func(_ int, v T) bool { out = append(out, fmt.Sprintf(`%v`, v)); return true })
	return fmt.Sprintf("%v", out)
}
