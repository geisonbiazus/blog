title: Algorithms and Data Structures Series: Dynamic Arrays
author: Geison Biazus
description: This post is part of the algorithms and data structures series, a series of posts where I present the most common data structures and algorithms used in software engineering. In this post, I explain the basics of Arrays, the first data structure in the series.
image_path: /static/image/logo-small.png
time: 2022-06-22 09:00
--

This post is part of the algorithms and data structures series, a series of posts where I present the most common data structures and algorithms used in software engineering. In this post, I explain the basics of Arrays, the first data structure in the series.

## Contents

- [Data Structures](#data-structures)
- [Static Arrays](#static-arrays)
- [Dynamic Arrays](#dynamic-arrays)
- [Operations](#operations)
  - [Push](#push)
  - [Insert](#insert)
  - [Update](#update)
  - [Lookup](#lookup)
  - [Delete](#delete)
  - [Search](#search)
- [Complexity table](#complexity-table)
- [Choosing dynamic arrays over other data structures](#choosing-dynamic-arrays-over-other-data-structures)
  - [When to use](#when-to-use)
  - [When to avoid](#when-to-avoid)
- [Final toughts](#final-thoughts)
- [Full dynamic array implementation](#full-dynamic-array-implementation)
- [Sources](#sources)

## Data Structures

Before digging into Arrays, first, let's understand what is a data structure.

Everyone that writes code, uses data structures on a daily basis. [Wikipedia](https://en.wikipedia.org/wiki/Data_structure) states that _"...a data structure is a data organization, management, and storage format that enables efficient access and modification..."_. In other words, it is a collection of objects arranged in a way that facilitates certain operations based on the necessities of the algorithm using it.

Some examples of data structures are Arrays, Hash Tables, Linked Lists, Stacks, Queues, Trees, and Graphs. Each one is specialized for certain usages and operations. For example, Arrays have fast access by index, but inserting at the beginning of the array is slow. On the other hand, Linked Lists are fast to insert at the beginning of the list, but are slow to access by index.

The choice of the right data structure depends on the necessities of structure or performance of the algorithm making use of them. I'll be writing about each one of these data structures in future posts.

## Static Arrays

Most of the statically-typed languages provide a way of creating static arrays. Dynamically-typed languages usually only provide dynamic arrays, which we will be seeing later in this post.

A static array is a sequence of elements of the same type. This sequence contains a fixed size defined on its declaration. In the Go programming language, an array can be declared as follows:

```go
var ints [5]int
```

An array has a fixed size for its entire lifecycle. When it is created, the required memory to hold all of its values is allocated. The only operations supported by arrays are assigning to a position and reading from a position. These operations are referred as update and lookup.

```go
ints[2] = 3 // update
fmt.Println(ints[2]) //lookup
```

Any operation in a single index of an array is a direct access to the memory position of that index. Both the update and lookup operations have a time complexity of `O(1)`.

## Dynamic arrays

Although static arrays can give us a lot of possibilities, it is usually not enough for our needs. We don't always know the exact size of the array we need, and we also don't want to always keep track of the last index we inserted or how many positions of the array are filled. The data structure that adds these and other functionalities to the static arrays is the dynamic array.

Dynamic arrays are known in dynamic languages, such as Javascript or Ruby, simply as Arrays. Some other languages call them Lists. In Go, they are known as Slices.

To demonstrate how dynamic arrays are implemented, we are not going to use the various features of the Go Slices. Instead, we will make use of Slices only for the storage of our data and they will act as simple static arrays. Go's static arrays require their size to be pre-determined on their declaration. For our implementation, although the array size is pre-determined, it varies at runtime. That's why we'll be using Slices instead.

Let's start with the dynamic array struct and constructor function:

```go
type DynamicArray[T comparable] struct {
	data     []T
	capacity int
	Length   int
}

func NewDynamicArray[T comparable]() *DynamicArray[T] {
	return &DynamicArray[T]{
		data:     make([]T, 1),
		capacity: 1,
	}
}
```

The `DynamicArray` struct uses the (generics)[https://go.dev/doc/tutorial/generics] introduced on Go 1.18. This allows our dynamic array to store and retrieve any type determined in the array declaration at runtime and without type casting.

The `DynamicArray` contains three fields: `data` is a Slice that stores any item pushed to the array. `capacity` is an internal field that holds how many items the data slice can hold (assuming it is a static array). `Length` is the only public attribute, it contains the number of items this array currently has.

## Operations

Next, we are going to see the common operations in dynamic arrays with their implementation and complexity.

### Push

The `Push` method, is used to add new items to a dynamic array. This method is the main difference compared to a static array and it is what makes it "dynamic". Here is how it is implemented:

```go
func (a *DynamicArray[T]) Push(item T) {
	if a.Length == a.capacity {
		a.expandData()
	}

	a.data[a.Length] = item
	a.Length++
}

func (a *DynamicArray[T]) expandData() {
	newCapacity := a.capacity * 2
	newData := make([]T, newCapacity)

	for i, item := range a.data {
		newData[i] = item
	}

	a.capacity = newCapacity
	a.data = newData
}
```

If we ignore the conditional part at the beginning of the `Push` method for now, adding a new item to the array is done by simply setting a new value to the position of the `Length` property and incrementing this length. This operation has an `O(1)` complexity and this is the complexity we get most of the time when we push items to a dynamic array.

Now, back to the conditional path, to avoid setting a value in a position that goes beyond its capacity, the `data` slice needs to be expanded. That is what this conditional at the beginning of the method does. If the `data` reached its maximum capacity, a new slice having twice the previous capacity is created. Then all the values from the previous `data` are copied to the new extended slice that is replaced in the struct. This operation has an `O(n)` time complexity but subsequent pushes to the dynamic array will have an `O(1)` complexity until the `data` reaches its capacity again.

When the dynamic array length is small, we have frequent expansions. But as there are not much data in the array yet, these expansions do not take much time. As the array gets bigger, the frequency of expansions decreases. So we say that the push operation to a dynamic array has and average time complexity of `O(1)`, having a worst-case of `O(n)`.

### Insert

Dynamic arrays allow inserting data at any position. This operation not only inserts the new item at the given position, but also moves all subsequent items one position up:

```go
func (a *DynamicArray[T]) Insert(index int, item T) error {
	if index < 0 || index > a.Length {
		return ErrIndexOutOfBounds
	}

	if a.Length == a.capacity {
		a.expandData()
	}

	a.shiftDataUp(index)

	a.Length++
	a.data[index] = item

	return nil
}

func (a *DynamicArray[T]) shiftDataUp(index int) {
	for i := a.Length; i > index; i-- {
		a.data[i] = a.data[i-1]
	}
}
```

In the `Insert` method, we first check if the given index is valid. The `ErrIndexOutOfBounds` error is returned in case it is invalid. Then we check whether the array length is at its capacity and expand it as we did on the `Push` method. Next, the `shifDataUp` method is called. This method loops through the array moving every subsequent item one index up to give space to the newly inserted item. Last, the `Length` is incremented and the new value is inserted at the appropriated position.

The shifting of values makes the `Insert` method have a complexity of `O(n)`.

### Update

Updating a value of a dynamic array is no different than updating a value of a static array:

```go
func (a *DynamicArray[T]) Set(index int, value T) error {
	if !a.isIndexWithinBounds(index) {
		return ErrIndexOutOfBounds
	}

	a.data[index] = value
	return nil
}

func (a *DynamicArray[T]) isIndexWithinBounds(index int) bool {
	return index >= 0 && index < a.Length
}
```

The `Set` method first validates if the given index is valid, then it replaces the value directly in the internal `data` slice. This operation has an `O(1)` complexity.

### Lookup

Similarly to updating, getting an item from the array by its index, is no different from getting it from a static array:

```go
func (a *DynamicArray[T]) Get(index int) (T, error) {
	if !a.isIndexWithinBounds(index) {
		return a.emptyValue(), ErrIndexOutOfBounds
	}

	return a.data[index], nil
}

func (a *DynamicArray[T]) emptyValue() T {
	var value T
	return value
}
```

The `Get` method validates whether the index is within bounds. If the index is invalid, an empty value and an error are returned. Go has the concept of "zero values" where depending on the type of a variable its empty state is already a value of that type. As we're using generics in our dynamic array, the empty value resolution is delegated to the compiler. Other languages would return null or throw an error where no value is returned.

If the index is valid, it returns the appropriate item from the `data` slice. This operation has an `O(1)` complexity.

### Delete

Static arrays do not support deleting values, only replacing values with null or something similar can be done to "delete" a value, but that would keep the subsequent values on their original position. In a dynamic array, as we have the possibility of pushing new items, removing items is also desirable.

```go
func (a *DynamicArray[T]) Delete(index int) error {
	if !a.isIndexWithinBounds(index) {
		return ErrIndexOutOfBounds
	}

	a.shiftDataDown(index)
	a.Length--
	a.data[a.Length] = a.emptyValue()

	return nil
}

func (a *DynamicArray[T]) shiftDataDown(index int) {
	for i := index + 1; i < a.Length; i++ {
		a.data[i-1] = a.data[i]
	}
}
```

The `Delete` method starts with our habitual index validation. Next, it calls the `shiftDataDown` method. This method loops from the index being removed until the end of the array moving all of its elements one index down. In this process, the item being removed gets overriden with the next item in the array. To finish the process, the length gets decremented and the last item, which got duplicated after the shifiting, is overriden with an empty value.

The `Delete` operation of a dynamic array has an `O(n)` complexity, but there is a catch here. Although the method has an `O(n)` complexity, knowing that we always shift the elements from the given index until the end of the array, we can assume that removing the last element has a complexity of `O(1)`.

### Search

The last operation that a dynamic array might have is searching. This can be searching for a value, a prefix, or even a pattern. As an example we are going to implement a `Contains` method, but any other kind of search would behave in a similar way with a similar complexity.

```go
func (a *DynamicArray[T]) Contains(item T) bool {
	for i := 0; i < a.Length; i++ {
		if a.data[i] == item {
			return true
		}
	}

	return false
}
```

The `Contains` method loops through all the array items comparing them with the given target item. As soon as the item is found it returns `true`. If the item is not found, `false` is returned. This operation has an `O(n)` complexity.

## Complexity table

Here is the time complexity table for the operations we just saw:

<table class="table">
  <thead>
    <tr>
      <th scope="col">Operation</th>
      <th scope="col">Complexity</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>Push</td>
      <td>O(1)</td>
    </tr>
    <tr>
      <td>Insert</td>
      <td>O(n)</td>
    </tr>
    <tr>
      <td>Update</td>
      <td>O(1)</td>
    </tr>
    <tr>
      <td>Lookup</td>
      <td>O(1)</td>
    </tr>
    <tr>
      <td>Delete</td>
      <td>O(n)</td>
    </tr>
    <tr>
      <td>Search</td>
      <td>O(n)</td>
    </tr>
  </tbody>
</table>

By analysing this table, we can see that dynamic arrays are very fast for adding new items to the end and updating and retrieving items based on their indexes. However, inserting, deleting or searching for items are slow operations on this data structure.

## Choosing dynamic arrays over other data structures

Dynamic arrays should be preferred for the operations they are fast. If we are going to use an array, we can look at the complexity table and maximize the use of its operations with `O(1)` complexity. Other data structures are optimized for the operations that arrays are not. For example, if we need to lookup by value, a hash table should be preferred. if we need fast insertions, a linked list could be a good option. On these mentioned data structures, while they are fast on some operations where arrays are slow, they can be slow in the ones where arrays are fast. So we should choose wisely based on our requirements.

Here is what we should consider when choosing or avoiding arrays in our algorithms:

### When to use

- Items should be ordered
- We need fast lookups by index
- Items are constantly added to the end of the array
- Items are updated by their index
- The last item is constantly removed

### When to avoid

- Multiple searches by value are performed
- Items are frequently inserted at the beginning or random positions of the array
- Items are frequently removed from the beginning or random positions of the array

## Final thoughts

Dynamic arrays enrich static arrays with many useful features. They are probably the most used data structure in modern programming languages. By knowing their pros and cons we can take full advantage of their fast operations while avoiding the slow ones.

I hope you enjoyed this post and learned something new. The next article in the algorithm and data structure series will be about hash tables.

## Full dynamic array implementation

```go
var ErrIndexOutOfBounds = errors.New("array index out of bounds")

type DynamicArray[T comparable] struct {
	data     []T
	capacity int
	Length   int
}

func NewDynamicArray[T comparable]() *DynamicArray[T] {
	return &DynamicArray[T]{
		data:     make([]T, 1),
		capacity: 1,
	}
}

func (a *DynamicArray[T]) Push(item T) {
	if a.Length == a.capacity {
		a.expandData()
	}

	a.data[a.Length] = item
	a.Length++
}

func (a *DynamicArray[T]) Insert(index int, item T) error {
	if index < 0 || index > a.Length {
		return ErrIndexOutOfBounds
	}

	if a.Length == a.capacity {
		a.expandData()
	}

	a.shiftDataUp(index)

	a.Length++
	a.data[index] = item

	return nil
}

func (a *DynamicArray[T]) shiftDataUp(index int) {
	for i := a.Length; i > index; i-- {
		a.data[i] = a.data[i-1]
	}
}

func (a *DynamicArray[T]) expandData() {
	newCapacity := a.capacity * 2
	newData := make([]T, newCapacity)

	for i, item := range a.data {
		newData[i] = item
	}

	a.capacity = newCapacity
	a.data = newData
}

func (a *DynamicArray[T]) Get(index int) (T, error) {
	if !a.isIndexWithinBounds(index) {
		return a.emptyValue(), ErrIndexOutOfBounds
	}

	return a.data[index], nil
}

func (a *DynamicArray[T]) Set(index int, value T) error {
	if !a.isIndexWithinBounds(index) {
		return ErrIndexOutOfBounds
	}

	a.data[index] = value
	return nil
}

func (a *DynamicArray[T]) Delete(index int) error {
	if !a.isIndexWithinBounds(index) {
		return ErrIndexOutOfBounds
	}

	a.shiftDataDown(index)
	a.Length--
	a.data[a.Length] = a.emptyValue()

	return nil
}

func (a *DynamicArray[T]) shiftDataDown(index int) {
	for i := index + 1; i < a.Length; i++ {
		a.data[i-1] = a.data[i]
	}
}

func (a *DynamicArray[T]) isIndexWithinBounds(index int) bool {
	return index >= 0 && index < a.Length
}

func (a *DynamicArray[T]) emptyValue() T {
	var value T
	return value
}

func (a *DynamicArray[T]) Contains(item T) bool {
	for i := 0; i < a.Length; i++ {
		if a.data[i] == item {
			return true
		}
	}

	return false
}

```

## Sources

- [Data structure - Wikipedia](https://en.wikipedia.org/wiki/Data_structure)
- [Master the Coding Interview: Data Structures + Algorithms](https://www.udemy.com/course/master-the-coding-interview-data-structures-algorithms/)
