package main

// In this video, we'll modify the existing HackerNews and Reddit client to incorporate our web server.
// ST
// Pull up your editor and let's get coding.

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/caser/gophernews"
	"github.com/jzelinskie/geddit"
)

var redditSession *geddit.LoginSession
var hackerNewsClient *gophernews.Client

// We need to add in the stories slice here
var stories []Story

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

func getHnStoryDetails(id int, c chan<- Story, wg *sync.WaitGroup) {
	defer wg.Done()
	story, err := hackerNewsClient.GetStory(id)
	if err != nil {
		return
	}
	newStory := Story{
		title:  story.Title,
		url:    story.URL,
		author: story.By,
		source: "HackerNews",
	}
	c <- newStory
}

func newHnStories(c chan<- Story) {
	defer close(c)
	changes, err := hackerNewsClient.GetChanges()
	if err != nil {
		fmt.Println(err)
		return
	}
	var wg sync.WaitGroup
	for _, id := range changes.Items {
		wg.Add(1)
		go getHnStoryDetails(id, c, &wg)
	}
	wg.Wait()
}

func newRedditStories(c chan<- Story) {
	defer close(c)
	sort := geddit.PopularitySort(geddit.NewSubmissions)
	var listingOptions geddit.ListingOptions
	submissions, err := redditSession.SubredditSubmissions("programming", sort, listingOptions)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, s := range submissions {
		newStory := Story{
			title:  s.Title,
			url:    s.URL,
			author: s.Author,
			source: "Reddit /r/programming",
		}
		c <- newStory
	}
}

// We don't need either output function
/*
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
}*/

// We do need a function that saves stories to the Stories list
func outputToStories(c <-chan Story) {
	for {
		s := <-c
		stories = append(stories, s)
	}
}

// Now we can just paste in our existing stories search function
func searchStories(query string) []Story {
	var foundStories []Story
	for _, story := range stories {
		if strings.Contains(strings.ToUpper(story.title), strings.ToUpper(query)) {
			foundStories = append(foundStories, story)
		}
	}
	return foundStories
}

// We also need to paste in our existing search and top ten functions
func search(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	if query == "" {
		http.Error(w, "Search parameter q is required to search.", http.StatusNotAcceptable)
		return
	}

	w.Write([]byte("<html><body>"))
	s := searchStories(query)
	if len(s) == 0 {
		w.Write([]byte(fmt.Sprintf("No results for query '%s'.\n<br>", r.FormValue("q"))))
	} else {
		for _, story := range s {
			w.Write([]byte(fmt.Sprintf("<a href='%s'>%s</a><br>by %s on %s<br><br>", story.url, story.title, story.author, story.source)))
		}
	}

	w.Write([]byte("<a href='../'>Back</a>"))
	w.Write([]byte("</body></html>"))

}

func topTen(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<html><body>"))
	form := "<form action='search' method='post'>Search: <input type='text' name='q'> <input type='submit'></form>\n"
	w.Write([]byte(form))
	for i := len(stories) - 1; i >= 0 && len(stories)-i < 10; i-- {
		story := stories[i]
		w.Write([]byte(fmt.Sprintf("<a href='%s'>%s</a><br>by %s on %s<br><br>", story.url, story.title, story.author, story.source)))
	}
	w.Write([]byte("</body></html>"))
}

func main() {

	// We can remove all the file handling logic as well.
	/*toFile := make(chan Story, 8)
	toPrint := make(chan Story, 8)
		file, err := os.Create("stories.txt")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		//go outputToConsole(toPrint)
		//go outputToFile(toFile, file)
	*/

	// We'll wrap the hackernews and reddit pulling logic in an anonymous function
	go func() {
		// We'll fetch stories in an infinite loop
		for {
			// We'll report that we're fetching stories
			fmt.Println("Fetching new stories...")

			// We'll create three channels - channels from our inputs and just one channel to our output
			fromHn := make(chan Story, 8)
			fromReddit := make(chan Story, 8)
			toList := make(chan Story, 8)
			// Now we'll spin off all three goroutines
			go outputToStories(toList)
			go newHnStories(fromHn)
			go newRedditStories(fromReddit)

			// The connector here is pretty similar...
			hnOpen := true
			redditOpen := true

			for hnOpen || redditOpen {
				select {
				case story, open := <-fromHn:
					if open {
						// Except that we only put things into
						// the toList channel
						toList <- story
					} else {
						hnOpen = false
					}
				case story, open := <-fromReddit:
					if open {
						toList <- story
					} else {
						redditOpen = false
					}
				}
			}
			// Now we'll report that we're finished and
			// wait a bit before getting new stories.
			fmt.Println("Done fetching new stories.")
			time.Sleep(30 * time.Second)
		}
	}()

	// Now we'll just start up the
	http.HandleFunc("/", topTen)
	http.HandleFunc("/search", search)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
