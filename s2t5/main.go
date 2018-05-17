package main

import (
	"fmt"
	"net/http"
	"sync"
)

// In this section, we'll talk about buffered channels and how they can be used to queue work.
// We'll process a bunch of web requests asyncronously, putting together all the techniques we just looked at.
// ST
// A buffered channel is like a regular channel, except that it can accept some number of messages that just
// get held in the channel until read. Unless full, writes don't block, and unless empty, reads don't block.
// ST
// Let's get coding!

// The first thing to do is to write the function that serves as the worker.
// It will accept an input channel as an argument.
// It also takes a WaitGroup, which we'll use slightly differently from before.
// It will return nothing.
// Here, the arrow syntax denotes that the channel is unidirectional.
func webGetWorker(in <-chan string, wg *sync.WaitGroup) {
	// This worker will run forever.
	for {
		// Accept a unit of work from the channel using arrow syntax.
		// This will block if the channel is empty.
		url := <-in
		// Perform the actual work, getting the web page. This is a
		// blocking call, but since it's running in a worker, it won't
		// block the main thread.
		res, err := http.Get(url)

		// Once complete, report the success or failure of the condition.
		if err != nil {
			// If there's an error, we'll print the error
			fmt.Println(err.Error())
		} else {
			// If the request was successful, report the status code
			fmt.Printf("GET %s: %d\n", url, res.StatusCode)
		}
		// We inform the WaitGroup of a "done" condition every time
		// a single unit of work is finished.
		wg.Done()
	}
}

// In the main function, we'll do the setup and teardown processing as well as queuing some work.
func main() {
	// Making a buffered channel is the same as a regular channel with a number of
	// buffer slots available.
	work := make(chan string, 1024)
	// Finally, create a WaitGroup
	var wg sync.WaitGroup

	// For ease of modificaiton, we'll put the number of workers in a variable.
	// Let's start off with 100 workers.
	numWorkers := 100
	// We'll spin off workers in a thread.
	for i := 0; i < numWorkers; i++ {
		// Within the loop, just spawn the Goroutine
		go webGetWorker(work, &wg)
	}

	// Now let's add some work to the channel. First, we'll put some URLs we want to get fetched in an array.
	urls := [6]string{"http://example.com", "http://packtpub.com", "http://reddit.com", "http://twitter.com", "http://facebook.com", "http://i.dont.exist"}
	// For testing purposes, we'll do each URL 100 times.
	for i := 0; i < 100; i++ {
		// Then loop over the URLs.
		for _, url := range urls {
			// Add one to the WaitGroup
			wg.Add(1)
			// Send the work into the channel. This won't block.
			work <- url
		}
	}

	// Finally, we just wait on the WaitGroup.
	// Every time a unit of work finishes, the WaitGroup gets decremented.
	wg.Wait()
}

// If we run this program, you can see that these GETs happen very very fast.

// Now, let's try converting the number of workers. Try setting it to 1; try setting it to 200.
// You'll notice a big difference. What's happening here is that, with just a few workers,
// the program spends a lot of time waiting. With lots of workers, it doesn't, so less time is wasted.
// ST
// In the next video, we'll look a little more in-depth at syncronization and non-buffered operations.
