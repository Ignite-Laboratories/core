package testing

import (
	"github.com/ignite-laboratories/core/atomic"
	"testing"
)

func Test_Atomic_Slice_NewSlice(t *testing.T) {
	s := atomic.NewSlice[int]()
	s.Add(5)
	s.Add(7)
	s.Add(9)
	s.RemoveIf(func(i int) bool {
		return i == 9
	})
	r := s.All()
	CompareIntSlices(r, []int{5, 7}, t)
}

func Test_Atomic_Slice_Add(t *testing.T) {
	s := atomic.NewSlice[int]()
	s.Add(5)
	s.Add(7)
	s.Add(9)
	r := s.All()
	CompareIntSlices(r, []int{5, 7, 9}, t)
}

func Test_Atomic_Slice_RemoveIf(t *testing.T) {
	s := atomic.NewSlice[int]()
	s.Add(5)
	s.Add(9)
	s.Add(7)
	s.Add(9)
	s.RemoveIf(func(i int) bool {
		return i == 9
	})
	r := s.All()
	CompareIntSlices(r, []int{5, 7}, t)
}

func Test_Atomic_Slice_RemoveIf_NoMatches(t *testing.T) {
	s := atomic.NewSlice[int]()
	s.Add(5)
	s.Add(9)
	s.Add(7)
	s.Add(9)
	s.RemoveIf(func(i int) bool {
		return i == 1
	})
	r := s.All()
	CompareIntSlices(r, []int{5, 9, 7, 9}, t)
}

func Test_Atomic_Slice_Length(t *testing.T) {
	s := atomic.NewSlice[int]()
	CompareInts(s.Length(), 0, t)
	s.Add(5)
	CompareInts(s.Length(), 1, t)
	s.Add(9)
	CompareInts(s.Length(), 2, t)
	s.Add(7)
	CompareInts(s.Length(), 3, t)
	s.Add(9)
	CompareInts(s.Length(), 4, t)
}

func Test_Atomic_Slice_Get(t *testing.T) {
	s := atomic.NewSlice[int]()
	s.Add(5)
	s.Add(9)
	s.Add(7)
	s.Add(9)

	CompareInts(5, s.Get(0), t)
	CompareInts(9, s.Get(1), t)
	CompareInts(7, s.Get(2), t)
	CompareInts(9, s.Get(3), t)
}

func Test_Atomic_Slice_All(t *testing.T) {
	s := atomic.NewSlice[int]()
	r := s.All()
	CompareIntSlices(r, []int{}, t)

	s.Add(5)
	s.Add(7)
	s.Add(9)
	r = s.All()
	CompareIntSlices(r, []int{5, 7, 9}, t)
}
