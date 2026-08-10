/*
Complete the fizzbuzz function that prints the numbers 1 to 100 inclusive each on their own line,
but replace multiples of 3 with the text "fizz" and multiples of 5 with "buzz".
Print "fizzbuzz" for multiples of 3 AND 5.
*/
package main

import "fmt"

func main() {
	fizzbuzz()
}

func fizzbuzz() {
	for i := 1; i <= 100; i++ {
		isMultipleOf3 := i%3 == 0
		isMultipleOf5 := i%5 == 0

		if isMultipleOf3 && isMultipleOf5 {
			fmt.Println("fizzbuzz")
		} else if isMultipleOf3 {
			fmt.Println("fizz")
		} else if isMultipleOf5 {
			fmt.Println("buzz")
		} else {
			fmt.Println(i)
		}
	}
}
