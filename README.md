# NSet

NSet is a super fast and memory efficient set implementation built for unsigned integers up to and including uint32.

By 'set' we mean something like a hash map, but instead of key/value pairs there are only keys.
You can do the normal operations of add, check if item exists, and delete, but you can also do things like union sets and
get intersections.

**Contents**:

- [NSet](#nset)
  - [When to use NSet](#when-to-use-nset)
  - [Usage](#usage)
  - [Benchmarks](#benchmarks)

## When to use NSet

Maybe you want a set implementation? Then this is one, but there are other reasons.

If you are using your hash maps/arrays like sets or do a lot of checks to see if items exists in your hash maps then NSet might make sense.
In such cases NSet makes sense because it is both faster and more memory efficient. You can see more about this in the Benchmarks section.

Here are some examples where you might want to consider NSet:

``` go
//You might be using maps mostly for checking if things exist:

//This map is being used like a set. Some people might also do: make(map[uint32]bool, 0)
mapOfIds := make(map[uint32]struct{}, 0)

//Fill map here...

someId := 54312
if _, ok:= mapOfIds[someId]; ok {
    //Do something
} else {
    //Something else
}
```

```go
//You might be searching arrays a lot
func ExistsInArray(myArray []int, item int) bool {

    for i := 0; i < len(myArray); i++ {
        if myArray[i] == item {
            return true
        }
    }

    return false
}
```

## Usage

To install run `go get github.com/bloeys/nset`

Then usage is very simple:
```go

mySet := nset.NewNSet[uint32]()

mySet.Add(0)
mySet.Add(300)
mySet.Add(256)
mySet.Add(4)

if mySet.Contains(5) {
    panic("Oops I don't want 5!")
}

mySet.Remove(4)

```

## Benchmarks

NSet is faster than the built-in Go hash map in all operations (add, check, delete) by `1.6x to 64x` depending on the operation and data size.

Benchmark with 100 elements:

![Benchmark of 100 elements](./.res/bench-100.png)

Benchmark with 10,000,000 elements:

![Benchmark of 10,000,000 elements](./.res/bench-10-million.png)

As can be seen from the benchmarks, NSet has almost no change in its performance even with 10 million elements, while the
hash map slows down a lot as the size grows. NSet practically doesn't allocate at all. But it should be noted that
allocation can happen when adding a number bigger than all previously entered numbers.

Benchmarks that have 'Rand' in them mean that access patterns are randomized which can cause cache invalidation.
To make sure the test is fair the seed is the same for both Go Map and NSet. Here both suffer slowdowns but NSet remains faster.

Benchmarks that have `Presized` in them means that the data structure was fully allocated before usage, like:

```go
//This map already has space for ~100 elements and so doesn't need to resize, which is costly
myMap := make(map[uint16], 100)
```

Map benefits from sizing while NSet isn't affected, but in both cases NSet remains faster.
