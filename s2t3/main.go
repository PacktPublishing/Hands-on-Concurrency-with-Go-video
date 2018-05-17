package main

import (
	"fmt"
	"sync"
)

// In this video, we'll discuss some fundamental components of sharing work between Goroutines.
// ST
// The simplest way to share data between goroutines is to pass a pointer to each goroutine which
// points to the same piece of memory. Each Goroutine can read and write to this memory whenever it
// wants to.
// This is easy but can have downsides.
// ST
// Let's jump into the editor and take a look.
// One important structure that is designed to be shared between Goroutines is the WaitGroup, a data
// structure that allows us to forgoe the explicit timeout we used in the previous examples.

// In this example, I'll just run a simple "even or odd" program.
// Let's write the actual computation as a function. It takes an integer argument by value,
// and a pointer to a WaitGroup. All the Goroutines will be using the same WaitGroup.
func printEven(x int, wg *sync.WaitGroup) {
	// The actual computation here is trivial.
	// If even, print even.
	// Otherwise, print odd.
	if x%2 == 0 {
		fmt.Printf("%d is even\n", x)
	} else {
		fmt.Printf("%d is odd\n", x)
	}

	// Finally, we'll inform the WaitGroup the function has completed.
	wg.Done()
}

// Now, in our main function, we'll use this printEven function asyncronously.
func main1() {
	// First, lets make a WaitGroup.
	// It has an internal counter, and it will initially be set to zero.
	var wg sync.WaitGroup
	// We can use a for loop to spawn a bunch of goroutines
	for i := 1; i < 10; i++ {
		// Before creating each Goroutine, we add to the WaitGroup's internal counter.
		wg.Add(1)
		// Then we pass the WaitGroup to the worker.
		// Remember, when the function is done, it calls Done on the WaitGroup.
		// Done decrements the counter by one.
		go printEven(i, &wg)
	}

	// Finally, we tell the program to wait until the WaitGroup is finished.
	// That is, this function blocks until the WaitGroup's counter is at zero.
	wg.Wait()
}

// Running this program, we can see that the WaitGroup allows us to wait only as long as
// needed for all Goroutines to complete.

// So, we've seen an example of how communicating by sharing memory works well. Now let's see
// how it can go wrong. Our function will be really similar, but instead of having an independent
// computation per Goroutine, we'll make all the Goroutines do the same thing to the same data.

// In this case, it's a simple increment. It takes a pointer to an integer and a pointer to
// a WaitGroup.
func increment(ptr *int, wg *sync.WaitGroup) {
	// To illustrate the bug, we'll first extract the value of the integer.
	i := *ptr
	// Then we'll do something that takes a little time, like print.
	fmt.Printf("value is %d\n", i)
	// Now we'll use the value we extracted, adding to it and putting it back in the pointer.
	*ptr = i + 1
	wg.Done()
}

// The main function is pretty much exactly the same.
func main() {
	// Make a WaitGroup, as above.
	var wg sync.WaitGroup
	// Val will be the value that the increment function operates on.
	val := 0
	// As before, we'll spawn goroutines in a loop.
	for i := 0; i < 10000; i++ {
		// Again, add to the WaitGroup every time.
		wg.Add(1)
		// Then run the actual operation.
		go increment(&val, &wg)
	}

	wg.Wait()
	// Once everything is done, we can look at the final value.
	// It should be 10000.
	fmt.Printf("Final value was %d\n", val)
}

// Running the program, you should see that it's not always 10000.
// That's because each Goroutine sets the value to whatever it saw plus one,
// potentially undoing the work of other Goroutines.
// ST
// In the next video, we'll see how to overcome this issue with one of Go's most
// powerful constructs; channels.
