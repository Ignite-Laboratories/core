package testing

import (
	"github.com/ignite-laboratories/core/atomic"
	"testing"
)

func Test_Atomic_Slice_NewSlice(t *testing.T) {
	s := atomic.NewSlice[int]()
	s.Add(5)
	s.Add(7)
	r := s.All()
	CompareIntSlices(r, []int{5, 7}, t)
}
