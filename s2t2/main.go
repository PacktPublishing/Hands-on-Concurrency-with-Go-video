package main

import (
	"fmt"
	"time"
)

// In this video, we'll talk about Goroutines.

// ST

// Goroutines are pieces of code that run at the same time as other code.
// They run independently of each other, and of the main code.
// Every concurrent Go program is built on Goroutines.

// ST

// Open up your editor and lets get coding!

// First off, let's see exactly how to write a Goroutine.
// Create a main() function just like the regular Hello, World! program,
// but in this case, write the keyword go in front of the Println invocation.
func main1() {
	go fmt.Println("Hello, World!")
}

// Now run the program, and...
// nothing happpens!

// This is because all "go function()" does is create a goroutine. The main function execution continues
// immediately, not waiting for the goroutine to finish.
// To see the goroutine's result, we need to wait for it to finish. For now, we'll just wait a fixed time.

func main2() {
	go fmt.Println("Hello, World!")
	// We'll use the time.Sleep function to wait a specific time period.
	// In this case, 1 second.
	time.Sleep(1 * time.Second)
}

// Now when we run the program, it does what we'd expect.
// A better demonstration is printing a bunch of things. This will demonstrate how goroutines execute out of order.

func main() {
	// We'll spawn Goroutines in a for loop.
	for i := 0; i < 10; i++ {
		// Each Goroutine will just print out the iteration it was spawned on
		go fmt.Printf("Goroutine number %d\n", i)
	}
	// We can also print out a string to see when the loop finishes.
	fmt.Println("For loop done!")
	// As before, we'll wait one second for all the goroutines to complete.
	time.Sleep(1 * time.Second)
}

// Running this in the terminal, we'll see that the Goroutines don't run in order!
// In fact, many of the complete before the for loop is even finished.
// Each time we run the program, the order will be different.

// Back to slideshow / ST

// Goroutines are able to accomplish this behavior by being run logically - and sometimes
// physically - at the same time. They do this with what's called M to N concurrency.
//
// In many programming languages, 1 to 1 concurrency is the norm. Each concurrent piece of the
// program gets its own OS thread, meaning it has its own context and stack, often up to a
// megabyte in size, and is kept track of by the operating system.
//
// ST
//
// Go, however, uses M to N concurrency, in which multiple Goroutines are assigned to each thread.
// This has many advantages. Primarily, it makes the creation of Goroutines so cheap as to be
// effectively free. It's a common practice to treat spawning a Goroutine as no less efficient than
// simply calling a function.
//
// ST
//
// This is actually a little bit more complex. Since, for example, an image processing function will
// probably use its thread's execution time very greedily, while simply checking on whether a file
// read is done will only need to happen very infrequently, Go is smart about how many Goroutines
// are assigned to one thread. Processor-intensive tasks may get their own thread so they can run
// close to 100% of the time on one core, while many I/O tasks might get scheduled in a thread together.
//
// In the next video, we'll talk about how Goroutines can be used to solve more realistic problems by
// sharing memory.
