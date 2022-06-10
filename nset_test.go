package nset_test

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/bloeys/nset"
)

const (
	maxBenchSize = 10_000_000
	RandSeed     = 9_812_938_704
)

var (
	dump          int
	fullRangeNSet *nset.NSet[uint32]
)

func TestNSet(t *testing.T) {

	n := nset.NewNSet[uint32]()
	n.Add(0)
	n.Add(1)
	n.Add(63)
	n.Add(math.MaxUint32)

	AllTrue(t, n.Contains(0), n.Contains(1), n.Contains(63), n.Contains(math.MaxUint32), !n.Contains(10), !n.Contains(599))
	AllTrue(t, n.ContainsAll(0, 1, 63), !n.ContainsAll(9, 0, 1), !n.ContainsAll(0, 1, 63, 99))
	AllTrue(t, n.ContainsAny(0, 1, 63), n.ContainsAny(9, 99, 999, 1), !n.ContainsAny(9, 99, 999))

	IsEq(t, nset.BucketCount-1, n.GetBucketIndex(math.MaxUint32))
	IsEq(t, math.MaxUint32/64/nset.BucketCount, n.GetStorageUnitIndex(math.MaxUint32))

	nCopy := n.Copy()
	n.Remove(1)

	AllTrue(t, n.Contains(0), n.Contains(63), !n.Contains(1), nCopy.ContainsAll(0, 1, 63, math.MaxUint32))

}

func TestNSetFullRange(t *testing.T) {

	if fullRangeNSet == nil {

		fullRangeNSet = nset.NewNSet[uint32]()
		println("Adding all uint32 to NSet...")
		for i := uint32(0); i < math.MaxUint32; i++ {
			fullRangeNSet.Add(i)
			if i%1_000_000_000 == 0 {
				fmt.Printf("i=%d billion\n", i)
			}
		}
		fullRangeNSet.Add(math.MaxUint32)
	}

	n := fullRangeNSet
	IsEq(t, 67_108_864, n.StorageUnitCount)
	for i := 0; i < len(n.Buckets); i++ {

		b := &n.Buckets[i]
		IsEq(t, 524288, b.StorageUnitCount)

		for j := 0; j < len(b.Data); j++ {
			if b.Data[j] != math.MaxUint64 {
				t.Errorf("Error: storage unit is NOT equal to MaxUint64 (i=%d,j=%d)! Expected math.MaxUint64 but got '%08b'\n",
					i,
					j,
					b.Data[j])
			}
		}
	}

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

func BenchmarkNSetAddRandNoSizeLimit(b *testing.B) {

	n := nset.NewNSet[uint32]()

	rand.Seed(RandSeed)
	for i := 0; i < b.N; i++ {
		n.Add(rand.Uint32())
	}
}

func BenchmarkMapAddRand(b *testing.B) {

	hMap := map[uint32]struct{}{}

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

func BenchmarkNSetContainsRandFullRange(b *testing.B) {

	//Init
	if fullRangeNSet == nil {

		b.StopTimer()

		fullRangeNSet = nset.NewNSet[uint32]()
		println("Preparing full range NSet...")
		for i := uint32(0); i < math.MaxUint32; i++ {
			fullRangeNSet.Add(i)
		}
		fullRangeNSet.Add(math.MaxUint32)

		b.StartTimer()
	}

	n := fullRangeNSet

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
