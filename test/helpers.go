package test

import "testing"

func ShouldPanic(t *testing.T) {
	if r := recover(); r == nil {
		t.Errorf("Expected panic, but didn't get one")
		t.FailNow()
	}
}

func CompareValues[T comparable](a T, b T, t *testing.T) {
	if a != b {
		t.Errorf("Expected %v, got %v", a, b)
		t.FailNow()
	}
}

func CompareSlices[T comparable](a []T, b []T, t *testing.T) {
	if len(a) != len(b) {
		t.Errorf("Slices are not the same length")
		t.FailNow()
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			t.Errorf("Expected %v at [%d], got %v", a[i], i, b[i])
			t.FailNow()
		}
	}
}
