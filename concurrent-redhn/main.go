package main

// In this video, we'll make our Reddit and HackerNews client concurrent.
// ST
// Open up your editor, and let's get coding!

import (
	"fmt"
	"os"
	"sync"

	"github.com/caser/gophernews"
	"github.com/jzelinskie/geddit"
)

var redditSession *geddit.LoginSession
var hackerNewsClient *gophernews.Client

func init() {
	hackerNewsClient = gophernews.NewClient()
	var err error
	redditSession, err = geddit.NewLoginSession("g_d_bot", "K417k4FTua52", "gdAgent v0")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Story struct {
	title  string
	url    string
	author string
	source string
}

// We need to change the new stories functions to work concurrently.
// First, we'll make a function that gets details on individual HackerNews
// stories.
// It'll take an ID, a channel to send the Story through, and a WaitGroup pointer.
func getHnStoryDetails(id int, c chan<- Story, wg *sync.WaitGroup) {
	// First off, we defer the finishing of the waitgroup.
	defer wg.Done()
	// Then, we do the actual work, getting the story
	story, err := hackerNewsClient.GetStory(id)
	// And simply giving up on an error.
	if err != nil {
		return
	}
	// We'll create the new Story struct as before
	newStory := Story{
		title:  story.Title,
		url:    story.URL,
		author: story.By,
		source: "HackerNews",
	}
	// And send it through the channel.
	c <- newStory
}

// Now we'll modifiy the newHnStories function. It'll take a channel instead of
// returning a slice, and close that channel when it's done.
func newHnStories(c chan<- Story) {
	// We'll defer closing the channel, so no matter how the function exits it'll get done.
	defer close(c)
	// Now we'll pull back the list of changes
	changes, err := hackerNewsClient.GetChanges()
	// And handle errors in the usual way
	if err != nil {
		fmt.Println(err)
		return
	}
	// Once we get the list of changes, we know how many we'll need
	// to perform, so we can create a waitgroup
	var wg sync.WaitGroup
	// Now we loop over the items
	for _, id := range changes.Items {
		// For each item we'll add one to the waitgroup
		wg.Add(1)
		// And spin off a Goroutine to get its details
		go getHnStoryDetails(id, c, &wg)
	}
	// All that's left now is to wait on that waitgroup
	wg.Wait()
}

// For the Reddit client function, we essentially only need to convert it to
// use a channel rather than a slice.
func newRedditStories(c chan<- Story) {
	// As before, we'll defer closing the channel
	defer close(c)
	sort := geddit.PopularitySort(geddit.NewSubmissions)
	var listingOptions geddit.ListingOptions
	submissions, err := redditSession.SubredditSubmissions("programming", sort, listingOptions)
	if err != nil {
		fmt.Println(err)
		// And we need to remove the "nil" here; there's no return type anymore
		return
	}
	for _, s := range submissions {
		newStory := Story{
			title:  s.Title,
			url:    s.URL,
			author: s.Author,
			source: "Reddit /r/programming",
		}
		// Rather than appending to the slice, we'll just send it through the channel
		c <- newStory
	}
	// And there's no need to return the slice anymore.
}

// Now we need two relatively simple functions.
// They take a channel and just output what they get from it, either to a file
// or to the console
func outputToConsole(c <-chan Story) {
	for {
		s := <-c
		fmt.Printf("%s: %s\nby %s on %s\n\n", s.title, s.url, s.author, s.source)
	}
}

func outputToFile(c <-chan Story, file *os.File) {
	for {
		s := <-c
		fmt.Fprintf(file, "%s: %s\nby %s on %s\n\n", s.title, s.url, s.author, s.source)
	}
}

// The main function requires a lot of changes.
func main() {
	// First off, we need four channels.
	fromHn := make(chan Story, 8)
	fromReddit := make(chan Story, 8)
	toFile := make(chan Story, 8)
	toPrint := make(chan Story, 8)
	// Now, we'll pass two of those channels to the two functions we created
	// and spin them off as goroutines
	go newHnStories(fromHn)
	go newRedditStories(fromReddit)

	// We'll start opening the output file while those operations
	// are working (remember, network is slower than disk.)
	file, err := os.Create("stories.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Once the file is open, we'll spin off the output functions.
	go outputToConsole(toPrint)
	go outputToFile(toFile, file)

	// Now, we'll connect the channels. We'll use a select statement to recieve
	// any stories that are available.
	// We'll use two boolean variables to track the status of the channels.
	hnOpen := true
	redditOpen := true
	// As long as one is open, there are more stories to recieve.
	for hnOpen || redditOpen {
		select {
		// Both cases will actually do the same thing
		case story, open := <-fromHn:
			if open {
				toFile <- story
				toPrint <- story
			} else {
				hnOpen = false
			}
		case story, open := <-fromReddit:
			if open {
				toFile <- story
				toPrint <- story
			} else {
				redditOpen = false
			}
		}
	}
}

// Running this modified program, you should see that it runs much, much faster.
// For me, it takes about 1.5 seconds - less than half the time of the previous
// implementation. I/O concurrency is one of the most effective uses of the technique.
// ST
// Now that you've seen the performance benefits of I/O concurrency, we'll use
// the next video to talk about some other benefits.
