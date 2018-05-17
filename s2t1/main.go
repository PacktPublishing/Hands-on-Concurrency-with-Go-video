// Every Go program must begin with a package name. If it's an executable, that has to be main.

package main

// Imports are generally written at the top of the program.
// We're only importing "fmt" because that's all we need for this simple program.
import "fmt"

// Now, the actual code goes in the main function.
func main1() {
	// We'll start out by just printing Hello, world!
	fmt.Println("Hello, world!")
}

// Go ahead and run the program with go run main.go.
// You should see the output Hello, world!
// If you don't, there's a problem with your Go configuration.

// Let's write something more substantial. First of all, recall how to define a function.
// We're going to write a recursive Lucasoid sequence generator - this is like Fibonacci, but a little different.
// We start with the func keyword, followed by the name of the function and it's arguments.
// In this case, we take three integer arguments: two starting numbers and the number to generate.
// Finally, there is the return type; integer in this case.

func lucasoid(a, b, n int) int {
	// Now we have the two base cases, using if statements.
	// In Go, if statement conditions don't require parentheses.
	if n == 0 {
		return a
	}
	if n == 1 {
		return b
	}

	// If neither of these conditions happen,
	// the function must recurse.
	return lucasoid(a, b, n-1) + lucasoid(a, b, n-2)
}

// Now in the main function, let's call the new function we created in a loop, to get a good look at the sequence.
// Recall that Go has only for loops.

func main() {
	// We'll print the first ten numbers in the Fibonacci and Lucas sequences using the lucasoid function.
	// So we create a loop from 0 to 9. We use the colon-equals syntax to define the iteration variable
	// then the loop condition
	// then the action (increment)
	for i := 0; i < 10; i++ {
		// Within the loop, we'll call lucasoid twice
		// First for the Fibonacci numbers (starting with 0 and 1)
		fib := lucasoid(0, 1, i)
		// Then for the Lucas numbers (starting with 2 and 1)
		luc := lucasoid(2, 1, i)
		// Then print out the result using the percent-d placeholder to print a digit.
		fmt.Printf("I: %d FIB: %d LUC: %d\n", i, fib, luc)
	}
}

// There you go - a simple linear Go program. I hope this has sufficiently refreshed your memory on Go.
// In the next video, we'll look at what Goroutines are, and how they can be used to parallelize this algorithm.
