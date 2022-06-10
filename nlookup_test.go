package nset_test

import (
	"math/rand"
	"testing"

	"github.com/bloeys/nset"
)

const (
	maxBenchSize = 10_000_000
	RandSeed     = 9_812_938_704
)

var (
	dump int
)

func TestNSet(t *testing.T) {

	n := nset.NewNSet[uint32]()
	IsEq(t, 1, cap(n.Data))

	n.Add(0)
	n.Add(1)
	n.Add(63)

	AllTrue(t, n.Contains(0), n.Contains(1), n.Contains(63), !n.Contains(10), !n.Contains(599))
	AllTrue(t, n.ContainsAll(0, 1, 63), !n.ContainsAll(9, 0, 1), !n.ContainsAll(0, 1, 63, 99))
	AllTrue(t, n.ContainsAny(0, 1, 63), n.ContainsAny(9, 99, 999, 1), !n.ContainsAny(9, 99, 999))

	n.Remove(1)
	AllTrue(t, n.Contains(0), n.Contains(63), !n.Contains(1))

	n = nset.NewNSetWithMax[uint32](100)
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

func BenchmarkNSetAdd(b *testing.B) {

	n := nset.NewNSet[uint32]()

	for i := uint32(0); i < uint32(b.N); i++ {
		n.Add(i % maxBenchSize)
	}
}

func BenchmarkMapAdd(b *testing.B) {

	hMap := map[uint32]struct{}{}

	for i := uint32(0); i < uint32(b.N); i++ {
		hMap[i%maxBenchSize] = struct{}{}
	}
}

func BenchmarkNSetAddRand(b *testing.B) {

	n := nset.NewNSet[uint32]()

	rand.Seed(RandSeed)
	for i := 0; i < b.N; i++ {
		n.Add(rand.Uint32() % maxBenchSize)
	}
}

func BenchmarkMapAddRand(b *testing.B) {

	hMap := map[uint32]struct{}{}

	rand.Seed(RandSeed)
	for i := 0; i < b.N; i++ {
		hMap[rand.Uint32()%maxBenchSize] = struct{}{}
	}
}

func BenchmarkNSetAddPresized(b *testing.B) {

	n := nset.NewNSetWithMax[uint32](maxBenchSize - 1)

	for i := uint32(0); i < uint32(b.N); i++ {
		n.Add(i % maxBenchSize)
	}
}

func BenchmarkMapAddPresized(b *testing.B) {

	hMap := make(map[uint32]struct{}, maxBenchSize-1)

	for i := uint32(0); i < uint32(b.N); i++ {
		hMap[i%maxBenchSize] = struct{}{}
	}
}

func BenchmarkNSetAddPresizedRand(b *testing.B) {

	n := nset.NewNSetWithMax[uint32](maxBenchSize - 1)

	rand.Seed(RandSeed)
	for i := 0; i < b.N; i++ {
		n.Add(rand.Uint32() % maxBenchSize)
	}
}

func BenchmarkMapAddPresizedRand(b *testing.B) {

	hMap := make(map[uint32]struct{}, maxBenchSize-1)

	rand.Seed(RandSeed)
	for i := 0; i < b.N; i++ {
		hMap[rand.Uint32()%maxBenchSize] = struct{}{}
	}
}

func BenchmarkNSetContains(b *testing.B) {

	//Init
	b.StopTimer()
	n := nset.NewNSet[uint32]()

	for i := uint32(0); i < maxBenchSize; i++ {
		n.Add(i)
	}
	b.StartTimer()

	//Work
	found := 0
	for i := uint32(0); i < uint32(b.N); i++ {
		if n.Contains(i) {
			found++
		}
	}

	dump = found
}

func BenchmarkMapContains(b *testing.B) {

	//Init
	b.StopTimer()
	hMap := map[uint32]struct{}{}

	for i := uint32(0); i < maxBenchSize; i++ {
		hMap[i] = struct{}{}
	}
	b.StartTimer()

	//Work
	found := 0
	for i := uint32(0); i < uint32(b.N); i++ {
		if _, ok := hMap[i]; ok {
			found++
		}
	}

	dump = found
}

func BenchmarkNSetContainsRand(b *testing.B) {

	//Init
	b.StopTimer()
	n := nset.NewNSet[uint32]()

	for i := uint32(0); i < maxBenchSize; i++ {
		n.Add(i)
	}
	b.StartTimer()

	//Work
	found := 0
	rand.Seed(RandSeed)
	for i := 0; i < b.N; i++ {

		randVal := rand.Uint32()
		if n.Contains(randVal) {
			found++
		}
	}

	dump = found
}

func BenchmarkMapContainsRand(b *testing.B) {

	//Init
	b.StopTimer()
	hMap := map[uint32]struct{}{}

	for i := uint32(0); i < maxBenchSize; i++ {
		hMap[i] = struct{}{}
	}
	b.StartTimer()

	//Work
	found := 0
	rand.Seed(RandSeed)
	for i := 0; i < b.N; i++ {

		randVal := rand.Uint32()
		if _, ok := hMap[randVal]; ok {
			found++
		}
	}

	dump = found
}

func BenchmarkNSetDelete(b *testing.B) {

	//Init
	b.StopTimer()
	n := nset.NewNSet[uint32]()

	for i := uint32(0); i < maxBenchSize; i++ {
		n.Add(i)
	}
	b.StartTimer()

	//Work
	for i := uint32(0); i < uint32(b.N); i++ {
		n.Remove(i)
	}
}

func BenchmarkMapDelete(b *testing.B) {

	//Init
	b.StopTimer()
	hMap := map[uint32]struct{}{}

	for i := uint32(0); i < maxBenchSize; i++ {
		hMap[i] = struct{}{}
	}
	b.StartTimer()

	//Work
	for i := uint32(0); i < uint32(b.N); i++ {
		delete(hMap, i)
	}
}

func BenchmarkNSetDeleteRand(b *testing.B) {

	//Init
	b.StopTimer()
	n := nset.NewNSet[uint32]()

	for i := uint32(0); i < maxBenchSize; i++ {
		n.Add(i)
	}
	b.StartTimer()

	//Work
	rand.Seed(RandSeed)
	for i := 0; i < b.N; i++ {

		randVal := rand.Uint32()
		n.Remove(randVal)
	}
}

func BenchmarkMapDeleteRand(b *testing.B) {

	//Init
	b.StopTimer()
	hMap := map[uint32]struct{}{}

	for i := uint32(0); i < maxBenchSize; i++ {
		hMap[i] = struct{}{}
	}
	b.StartTimer()

	//Work
	rand.Seed(RandSeed)
	for i := 0; i < b.N; i++ {

		randVal := rand.Uint32()
		delete(hMap, randVal)
	}
}
