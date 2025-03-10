package atomic

// Slice is a thread-safe collection of items.  Operations are performed using
// channels to ensure they complete fully before the next operation is permitted.
type Slice[T any] struct {
	data   chan func([]T) []T
	signal chan struct{}
}

// NewSlice creates a new instance of a Slice.
func NewSlice[T any]() *Slice[T] {
	as := &Slice[T]{
		data:   make(chan func([]T) []T),
		signal: make(chan struct{}),
	}
	go func() {
		var slice []T
		for f := range as.data {
			slice = f(slice)
		}
		close(as.signal)
	}()
	return as
}

// Add places the provided element at the end of the slice.
func (as *Slice[T]) Add(element T) {
	as.data <- func(slice []T) []T {
		return append(slice, element)
	}
}

// RemoveIf removes an element from the slice if the provided predicate returns true.
// If you would like to remove a specific element, you can provide an anonymous test
// as the predicate - or you could "fuzzily" match elements with a more complex predicate.
func (as *Slice[T]) RemoveIf(predicate func(T) bool) {
	as.data <- func(slice []T) []T {
		var result []T
		for _, v := range slice {
			if !predicate(v) {
				result = append(result, v)
			}
		}
		return result
	}

}

// Length returns the number of elements currently in the slice.
func (as *Slice[T]) Length() int {
	result := make(chan int)
	as.data <- func(slice []T) []T {
		result <- len(slice)
		return slice
	}
	return <-result
}

// Get returns the element at the provided index.
func (as *Slice[T]) Get(index int) T {
	result := make(chan T)
	as.data <- func(slice []T) []T {
		result <- slice[index]
		return slice
	}
	return <-result
}

// All returns all elements in a copy of the current slice.
func (as *Slice[T]) All() []T {
	result := make(chan []T)
	as.data <- func(slice []T) []T {
		result <- append([]T{}, slice...)
		return slice
	}
	return <-result
}

// Close closes all the inner channels.
func (as *Slice[T]) Close() {
	close(as.data)
	<-as.signal
}
