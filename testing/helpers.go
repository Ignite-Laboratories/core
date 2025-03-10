package testing

import "testing"

func CompareIntSlices(a []int, b []int, t *testing.T) {
	if len(a) != len(b) {
		t.Errorf("Slices are not the same length")
		t.FailNow()
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			t.Errorf("Expected %d at [%d], got %d", a[i], i, b[i])
			t.FailNow()
		}
	}
}
