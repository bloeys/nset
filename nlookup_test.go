package nlookup_test

import (
	"testing"

	"github.com/bloeys/nlookup"
)

func TestNLookup(t *testing.T) {

	n := nlookup.NewNLookup[uint]()
	IsEq(t, 1, cap(n.Data))

	n.Add(0)
	n.Add(1)
	n.Add(63)

	AllTrue(t, n.Contains(0), n.Contains(1), n.Contains(63), !n.Contains(10), !n.Contains(599))
	AllTrue(t, n.ContainsAll(0, 1, 63), !n.ContainsAll(9, 0, 1), !n.ContainsAll(0, 1, 63, 99))
	AllTrue(t, n.ContainsAny(0, 1, 63), n.ContainsAny(9, 99, 999, 1), !n.ContainsAny(9, 99, 999))

	n.Remove(1)
	AllTrue(t, n.Contains(0), n.Contains(63), !n.Contains(1))

	n = nlookup.NewNLookupWithMax[uint](100)
	IsEq(t, 2, cap(n.Data))
}

func AllTrue(t *testing.T, values ...bool) bool {

	for i := 0; i < len(values); i++ {
		if !values[i] {
			t.Errorf("Expected 'true' but got 'false'\n")
		}
	}

	return true
}

func IsEq[T comparable](t *testing.T, expected, val T) bool {

	if val == expected {
		return true
	}

	t.Errorf("Expected '%v' but got '%v'\n", expected, val)
	return false
}
