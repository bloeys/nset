package nlookup_test

import (
	"testing"

	"github.com/bloeys/nlookup"
)

func TestNLookup(t *testing.T) {

	n := nlookup.NewNLookup[uint]()

	IsEq(t, 1, cap(n.Data))

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
