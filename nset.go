package nset

import (
	"fmt"
	"strings"
)

var _ fmt.Stringer = &NSet[uint8]{}

type StorageType uint64

const StorageTypeBits = 64

//IntsIf is limited to uint32 because we can store ALL 4 Billion uint32 numbers
//in 256MB with NSet (instead of the normal 16GB for an array of all uint32s).
//But if we allow uint64 (or int, since int can be 64-bit) users can easily put a big 64-bit number and use more RAM than maybe Google and crash.
type IntsIf interface {
	uint8 | uint16 | uint32
}

type NSet[T IntsIf] struct {
	Data             []StorageType
	StorageUnitCount uint64
}

func (n *NSet[T]) Add(x T) {

	unitIndex := n.GetStorageUnitIndex(x)
	if unitIndex >= n.Size() {
		storageUnitsToAdd := unitIndex - n.Size() + 1
		n.Data = append(n.Data, make([]StorageType, storageUnitsToAdd)...)
		n.StorageUnitCount += storageUnitsToAdd
	}

	n.Data[unitIndex] |= 1 << (x % StorageTypeBits)
}

func (n *NSet[T]) Remove(x T) {

	unitIndex := n.GetStorageUnitIndex(x)
	if unitIndex >= n.Size() {
		return
	}

	n.Data[unitIndex] ^= 1 << (x % StorageTypeBits)
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
	unitIndex := n.GetStorageUnitIndex(x)
	return unitIndex < n.Size() && n.Data[unitIndex]&(1<<(x%StorageTypeBits)) != 0
}

func (n *NSet[T]) GetStorageUnitIndex(x T) uint64 {
	return uint64(x) / StorageTypeBits
}

func (n *NSet[T]) GetStorageUnit(x T) StorageType {
	return n.Data[x/StorageTypeBits]
}

//Size returns the number of storage units
func (n *NSet[T]) Size() uint64 {
	return n.StorageUnitCount
}

func (n *NSet[T]) ElementCap() uint64 {
	return uint64(len(n.Data) * StorageTypeBits)
}

//String returns a string of the storage as bytes separated by spaces. A comma is between each storage unit
func (n *NSet[T]) String() string {

	b := strings.Builder{}
	b.Grow(len(n.Data)*StorageTypeBits + len(n.Data)*2)

	for i := 0; i < len(n.Data); i++ {

		x := n.Data[i]
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

	return b.String()
}

func NewNSet[T IntsIf]() NSet[T] {

	return NSet[T]{
		Data:             make([]StorageType, 1),
		StorageUnitCount: 1,
	}
}

//NewNSetWithMax creates a set that already has capacity to hold till at least largestNum without resizing.
//Note that this is NOT the count of elements you want to store, instead you input the largest value you want to store. You can store larger values as well.
func NewNSetWithMax[T IntsIf](largestNum T) NSet[T] {
	return NSet[T]{
		Data:             make([]StorageType, largestNum/StorageTypeBits+1),
		StorageUnitCount: uint64(largestNum/StorageTypeBits + 1),
	}
}
