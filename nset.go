package nset

import (
	"fmt"
	"reflect"
	"strings"
)

var _ fmt.Stringer = &NSet[uint8]{}

type BucketType uint8
type StorageType uint64

const (
	BucketCount        = 128
	StorageTypeBits    = 64
	BucketIndexingBits = 7
)

//IntsIf is limited to uint32 because we can store ALL 4 Billion uint32 numbers
//in 256MB with NSet (instead of the normal 16GB for an array of all uint32s).
//But if we allow uint64 (or int, since int can be 64-bit) users can easily put a big 64-bit number and use more RAM than maybe Google and crash.
type IntsIf interface {
	uint8 | uint16 | uint32
}

type Bucket struct {
	Data             []StorageType
	StorageUnitCount uint32
}

type NSet[T IntsIf] struct {
	Buckets [BucketCount]Bucket
	//StorageUnitCount the number of uint64 integers that are used to indicate presence of numbers in the set
	StorageUnitCount uint32
	shiftAmount      T
}

func (n *NSet[T]) Add(x T) {

	bucket := n.GetBucketFromValue(x)

	unitIndex := n.GetStorageUnitIndex(x)
	if unitIndex >= bucket.StorageUnitCount {

		storageUnitsToAdd := unitIndex - bucket.StorageUnitCount + 1
		bucket.Data = append(bucket.Data, make([]StorageType, storageUnitsToAdd)...)

		n.StorageUnitCount += storageUnitsToAdd
		bucket.StorageUnitCount += storageUnitsToAdd
	}

	bucket.Data[unitIndex] |= n.GetBitMask(x)
}

func (n *NSet[T]) AddMany(values ...T) {

	for i := 0; i < len(values); i++ {

		x := values[i]
		bucket := n.GetBucketFromValue(x)

		unitIndex := n.GetStorageUnitIndex(x)
		if unitIndex >= bucket.StorageUnitCount {

			storageUnitsToAdd := unitIndex - bucket.StorageUnitCount + 1
			bucket.Data = append(bucket.Data, make([]StorageType, storageUnitsToAdd)...)

			n.StorageUnitCount += storageUnitsToAdd
			bucket.StorageUnitCount += storageUnitsToAdd
		}

		bucket.Data[unitIndex] |= n.GetBitMask(x)
	}

}

func (n *NSet[T]) Remove(x T) {

	b := n.GetBucketFromValue(x)
	unitIndex := n.GetStorageUnitIndex(x)
	if unitIndex >= b.StorageUnitCount {
		return
	}

	b.Data[unitIndex] ^= n.GetBitMask(x)
}

func (n *NSet[T]) Contains(x T) bool {
	return n.isSet(x)
}

func (n *NSet[T]) ContainsAny(values ...T) bool {

	for _, x := range values {
		if n.isSet(x) {
			return true
		}
	}

	return false
}

func (n *NSet[T]) ContainsAll(values ...T) bool {

	for _, x := range values {
		if !n.isSet(x) {
			return false
		}
	}

	return true
}

func (n *NSet[T]) isSet(x T) bool {
	b := n.GetBucketFromValue(x)
	unitIndex := n.GetStorageUnitIndex(x)
	return unitIndex < b.StorageUnitCount && b.Data[unitIndex]&n.GetBitMask(x) != 0
}

func (n *NSet[T]) GetBucketFromValue(x T) *Bucket {
	return &n.Buckets[n.GetBucketIndex(x)]
}

func (n *NSet[T]) GetBucketIndex(x T) BucketType {
	//Use the top 'n' bits as the index to the bucket
	return BucketType(x >> n.shiftAmount)
}

func (n *NSet[T]) GetStorageUnitIndex(x T) uint32 {
	//The top 'n' bits are used to select the bucket so we need to remove them before finding storage
	//unit and bit mask. This is done by shifting left by 4 which removes the top 'n' bits,
	//then shifting right by 4 which puts the bits back to their original place, but now
	//the top 'n' bits are zeros.
	return uint32(((x << BucketIndexingBits) >> BucketIndexingBits) / StorageTypeBits)
}

func (n *NSet[T]) GetBitMask(x T) StorageType {
	//Removes top 'n' bits
	return 1 << (((x << BucketIndexingBits) >> BucketIndexingBits) % StorageTypeBits)
}

func (n *NSet[T]) Union(otherSet *NSet[T]) {

	for i := 0; i < BucketCount; i++ {

		b1 := &n.Buckets[i]
		b2 := &otherSet.Buckets[i]

		if b1.StorageUnitCount < b2.StorageUnitCount {

			storageUnitsToAdd := b2.StorageUnitCount - b1.StorageUnitCount
			b1.Data = append(b1.Data, make([]StorageType, storageUnitsToAdd)...)

			b1.StorageUnitCount += storageUnitsToAdd
			n.StorageUnitCount += storageUnitsToAdd
		}

		for j := 0; j < len(b1.Data) && j < len(b2.Data); j++ {
			b1.Data[j] |= b2.Data[j]
		}
	}
}

func (n *NSet[T]) GetIntersection(otherSet *NSet[T]) *NSet[T] {

	outSet := NewNSet[T]()

	for i := 0; i < BucketCount; i++ {

		b1 := &n.Buckets[i]
		b2 := &otherSet.Buckets[i]

		//bucketIndexBits are the bits removed from the original value to use for bucket indexing.
		//We will use this to restore the original value 'x' once an intersection is detected
		bucketIndexBits := T(i << n.shiftAmount)
		for j := 0; j < len(b1.Data) && j < len(b2.Data); j++ {

			if b1.Data[j]&b2.Data[j] == 0 {
				continue
			}

			mask := StorageType(1 << 0)                                     //This will be used to check set bits. Numbers will be reconstructed only for set bits
			commonBits := b1.Data[j] & b2.Data[j]                           //Bits that are set on both storage units (aka the intersection)
			firstStorageUnitValue := T(j*StorageTypeBits) | bucketIndexBits //StorageUnitIndex = noBucketBitsX / StorageTypeBits. So: noBucketBitsX = StorageUnitIndex * StorageTypeBits; Then: x = noBucketBitsX | bucketIndexBits
			for k := T(0); k < StorageTypeBits; k++ {

				if commonBits&mask > 0 {
					outSet.Add(firstStorageUnitValue + k)
					// fmt.Printf("Bucket=%d, Storage unit=%d, bitPos=%d, value=%d\n", i, j, k, firstStorageUnitValue+k)
				}

				mask <<= 1
			}

		}
	}

	return outSet
}

func (n *NSet[T]) HasIntersection(otherSet *NSet[T]) bool {

	for i := 0; i < len(n.Buckets); i++ {

		b1 := &n.Buckets[i]
		b2 := &otherSet.Buckets[i]

		for j := 0; j < len(b1.Data) && j < len(b2.Data); j++ {

			if b1.Data[j]&b2.Data[j] > 0 {
				return true
			}
		}
	}

	return false
}

//String returns a string of the storage as bytes separated by spaces. A comma is between each storage unit
func (n *NSet[T]) String() string {

	b := strings.Builder{}
	b.Grow(int(n.StorageUnitCount*StorageTypeBits + n.StorageUnitCount*2))

	for i := 0; i < len(n.Buckets); i++ {

		bucket := &n.Buckets[i]
		for j := 0; j < len(bucket.Data); j++ {

			x := bucket.Data[j]
			shiftAmount := StorageTypeBits - 8
			for shiftAmount >= 0 {

				byteToShow := uint8(x >> shiftAmount)
				if shiftAmount > 0 {
					b.WriteString(fmt.Sprintf("%08b ", byteToShow))
				} else {
					b.WriteString(fmt.Sprintf("%08b", byteToShow))
				}

				shiftAmount -= 8
			}
			b.WriteString(", ")
		}
	}

	return b.String()
}

func (n *NSet[T]) Copy() *NSet[T] {

	newSet := NewNSet[T]()
	for i := 0; i < len(n.Buckets); i++ {

		b := &n.Buckets[i]
		newB := &newSet.Buckets[i]

		newB.StorageUnitCount = b.StorageUnitCount
		newB.Data = make([]StorageType, len(b.Data))

		copy(newB.Data, b.Data)
	}

	newSet.StorageUnitCount = n.StorageUnitCount
	return newSet

}

func NewNSet[T IntsIf]() *NSet[T] {

	n := &NSet[T]{
		Buckets:          [BucketCount]Bucket{},
		StorageUnitCount: 0,
		//We use this to either extract or clear the top 'n' bits, as they are used to select the bucket
		shiftAmount: T(reflect.TypeOf(*new(T)).Bits()) - BucketIndexingBits,
	}

	for i := 0; i < len(n.Buckets); i++ {
		n.Buckets[i].Data = make([]StorageType, 0)
	}

	return n
}
