/*
Arrays are fixed in size. Once you make an array like [5]int you can't add a 6th element.

A slice is a dynamically-sized, flexible view of the elements of an array.

The zero value of slice is nil.
see: https://go.dev/tour/moretypes/12

Non-nil slices always have an underlying array, though it isn't always specified explicitly. To explicitly create a slice on top of an array we can do:
*/
package main

import (
	"fmt"
	"reflect"
	"slices"
)

func main() {
	/*
		primes := [7]int{2, 3, 5, 7, 11, 13, 17}
		mySlice := primes[1:5]
		// mySlice = {3, 5, 7, 11}

		The syntax is:

		arrayName[low:high]
		arrayName[low:]
		arrayName[:high]
		arrayName[:]

		Where low is inclusive and high is exclusive.

		low, high, or both can be omitted to use the entire array on that side of the colon.

		the range is [low, high)
	*/
	primes := [7]int{2, 3, 5, 7, 11, 13, 17}

	mySlice := primes[:1] // prints [2]
	mySlice = primes[4:]  // prints [11, 13, 17]
	mySlice = primes[:]   // prints the whole thing
	mySlice = primes[2:4] // prints [5, 7]
	// mySlice = primes[4:2] // throws an error, invalid slice indices: 2 < 4

	fmt.Println(mySlice)

	fmt.Println("usingAppend(): ", usingAppend())
	fmt.Println("usingMake(): ", usingMake())
	fmt.Println("createMatrix(): ", createMatrix(3, 3))

	// COMPARING SLICESS
	fmt.Println("--------------------------------")
	fmt.Println("compareArrays(): ", compareArrays())
	fmt.Println("compareSlices(): ", compareSlices())
	fmt.Println("--------------------------------")

	// DELETING FROM SLICES
	fmt.Println("--------------------------------")
	fmt.Println("deletingFromSlices(): ", deletingFromSlices())
	fmt.Println("--------------------------------")
}

type cost struct {
	day   int
	value float64
}

func usingAppend() []float64 {
	res := []float64{}

	res = append(res, 1.0)
	res = append(res, 2.0)
	res = append(res, 3.0)

	return res
}

func usingMake() []float64 {
	/*
		make is a built-in function that creates a slice.
		It takes the type of the slice, the length of the slice, and the capacity of the slice.
		The capacity is the number of elements the slice can hold.
		The length is the number of elements the slice currently holds.
		The capacity is optional and defaults to the length.
	*/

	res := make([]float64, 3)
	matrix := [][]int{}

	fmt.Printf("matrix: %v\n", matrix)

	res[0] = 1.0
	res[1] = 2.0
	res[2] = 3.0

	return res
}

func createMatrix(rows, cols int) [][]int {
	matrix := [][]int{}

	for i := 0; i < rows; i++ {
		rowVals := []int{}

		for j := 0; j < cols; j++ {
			rowVals = append(rowVals, i*j)
		}

		matrix = append(matrix, rowVals)
	}

	return matrix
}

func compareArrays() bool {
	/*
		In Go, an array's size is part of its type definition.
		Because their sizes are fixed and known at compile time, Go allows you to use the == operator.
	*/

	array1 := [3]int{1, 2, 3}
	array2 := [3]int{1, 2, 3}

	fmt.Println("array1 == array2: ", array1 == array2)

	return array1 == array2 // true
}

func compareSlices() bool {
	/*
		Slices ([]int) do not have a fixed size at compile time.
		Because they are headers pointing to underlying arrays,
			Go explicitly forbids using == to compare two slices.
		If you try, the compiler will fail. To compare slices by value, you have three primary options:
	*/

	// OPTION A: Use slices.Equal (The Modern, Recommended Way)Starting in Go 1.21,
	/*
		the standard library includes the slices package.
		The slices.Equal function checks if the lengths are identical and if all elements match by value.
	*/
	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}

	fmt.Println("slices.Equal(s1, s2): ", slices.Equal(s1, s2))

	// OPTION B: Use slices.EqualFunc (for slices of Structs)
	/*
		If you have a slice of custom structs that cannot be compared with standard operators,
		you can use slices.EqualFunc. This lets you pass a custom comparison function.
	*/

	type User struct {
		ID   int
		Tags []string // Tags is a slice and makes the struct non-comparable with ==
	}

	u1 := []User{{ID: 1, Tags: []string{"tag1", "tag2"}}}
	u2 := []User{{ID: 1, Tags: []string{"tag1", "tag2"}}}

	// Custom rule: users are equal if their IDs match
	isEqual := slices.EqualFunc(u1, u2, func(a, b User) bool {
		return a.ID == b.ID
	})

	fmt.Println("slices.EqualFunc(u1, u2): ", isEqual)

	// OPTION C: Use reflect.DeepEqual (For Deeply Nested Data)
	/*
		If your slice contains deeply nested, complex data types (like maps inside slices),
		you can use the reflect package.
		However, avoid this unless necessary because it relies on runtime reflection,
		making it significantly slower than slices.Equal.
	*/

	r1 := []interface{}{"apple", map[string]int{"a": 1}}
	r2 := []interface{}{"apple", map[string]int{"a": 1}}

	fmt.Println("reflect.DeepEqual(r1, r2): ", reflect.DeepEqual(r1, r2))

	return true
}

func deletingFromSlices() bool {
	test1 := []int{1, 2, 3, 4, 5}
	test2 := []int{1, 2, 3, 4, 5}

	type TestStruct struct {
		id   int
		name string
	}

	test3 := []TestStruct{
		{id: 1, name: "John"},
		{id: 2, name: "Jane"},
		{id: 3, name: "Jim"},
	}

	test1 = append(test1[:2], test1[3:]...)
	fmt.Println("test1: ", test1)

	// using slices

	test2 = slices.Delete(test2, 2, 3)

	test3 = slices.DeleteFunc(test3, func(t TestStruct) bool {
		return t.id == 2
	})

	fmt.Println("test2: ", test2)
	fmt.Println("test3: ", test3)

	return true
}
