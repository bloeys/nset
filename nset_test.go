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

	n1 := nset.NewNSet[uint32]()
	n1.Add(0)
	n1.Add(1)
	n1.Add(63)
	n1.Add(math.MaxUint32)

	AllTrue(t, n1.Contains(0), n1.Contains(1), n1.Contains(63), n1.Contains(math.MaxUint32), !n1.Contains(10), !n1.Contains(599))
	AllTrue(t, n1.ContainsAll(0, 1, 63), !n1.ContainsAll(9, 0, 1), !n1.ContainsAll(0, 1, 63, 99))
	AllTrue(t, n1.ContainsAny(0, 1, 63), n1.ContainsAny(9, 99, 999, 1), !n1.ContainsAny(9, 99, 999))

	IsEq(t, nset.BucketCount-1, n1.GetBucketIndex(math.MaxUint32))
	IsEq(t, math.MaxUint32/64/nset.BucketCount, n1.GetStorageUnitIndex(math.MaxUint32))

	nCopy := n1.Copy()
	n1.Remove(1)

	AllTrue(t, n1.Contains(0), n1.Contains(63), !n1.Contains(1), nCopy.ContainsAll(0, 1, 63, math.MaxUint32))

	//Intersections
	n2 := nset.NewNSet[uint32]()
	n2.AddMany(1000, 63, 5, 10)

	n3 := nset.NewNSet[uint32]()
	n3.AddMany(math.MaxUint32)

	AllTrue(t, n1.HasIntersection(n2), n2.HasIntersection(n1), n3.HasIntersection(n1), !n3.HasIntersection(n2))

	n4 := nset.NewNSet[uint32]()
	n4.AddMany(0, 1, 64, math.MaxUint32)

	n5 := nset.NewNSet[uint32]()
	n5.AddMany(0, 1, 63, 64, math.MaxUint32)

	n4n5 := n4.GetIntersection(n5)
	AllTrue(t, n4n5.ContainsAll(0, 1, 64, math.MaxUint32), !n4n5.Contains(63))

	//Union
	n6 := nset.NewNSet[uint32]()
	n6.AddMany(4, 7, 100, 1000)

	n7 := nset.NewNSet[uint32]()
	n7.AddMany(math.MaxUint32)
	n7OldStorageUnitCount := n7.StorageUnitCount
	n7.Union(n6)

	AllTrue(t, n6.ContainsAll(4, 7, 100, 1000), !n6.Contains(math.MaxUint32), n7.ContainsAll(4, 7, 100, 1000, math.MaxUint32), n7.StorageUnitCount == n7OldStorageUnitCount+n6.StorageUnitCount)

}

func TestNSetFullRange(t *testing.T) {
	return
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
