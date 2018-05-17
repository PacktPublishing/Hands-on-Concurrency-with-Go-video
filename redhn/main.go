package main

// In this video, we'll be building a Reddit and HackerNews client - specifically, the non-concurrent version
// ST
// Let's start off with requirements. We'll build a simple terminal application which will load the latest
// programming news from both HackerNews and Reddit, display it in the terminal, and save it to a file.
// ST
// Open up your editor, and let's get coding!

// We'll be using some external packages, so we need to "go get" them before running the program.
import (
	"fmt"
	"os"

	"github.com/caser/gophernews"
	"github.com/jzelinskie/geddit"
)

// We need a variable for our Reddit API object
var redditSession *geddit.LoginSession

// And one for our HackerNews API object.
var hackerNewsClient *gophernews.Client

// In the init() function, we'll log into the two services we're using
func init() {
	// HackerNews allows API use without authentication, so we don't need an account.
	// We can just create our client object and use it.
	hackerNewsClient = gophernews.NewClient()
	// Reddit, on the other hand, does require authentication. I set up an account, but you'll
	// need to set up your own. It's free.
	// Here, I pass in the username, password, and user agent string the API client will use.
	var err error
	redditSession, err = geddit.NewLoginSession("g_d_bot", "K417k4FTua52", "gdAgent v0")
	// In case of an error, we'll just exit the program.
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Now we need a structure to store information about a story
type Story struct {
	// It'll have the title
	title string
	// And the URL to which the story leads
	url string
	// And finally the author of the story, and the source
	author string
	source string
}

// Now we need a function to get new stories from HackerNews
func newHnStories() []Story {
	// First we need a buffer to hold stories
	var stories []Story
	// We can use the GetChanges function to get all the most recent objects
	// from HackerNews. These will just be integer IDs, and we'll need to make requests
	// for each one.
	changes, err := hackerNewsClient.GetChanges()
	// In case of an error, we'll print it and return nil
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// Now, we can loop over the IDs we got back and make a request for each one.
	for _, id := range changes.Items {
		// The GetStory method will return a Story struct, or an error if the requested object wasn't
		// as story.
		story, err := hackerNewsClient.GetStory(id)
		// In case it wasn't a story, just move on.
		if err != nil {
			continue
		}
		// Now we can construct a Story struct and put it in the list
		newStory := Story{
			title:  story.Title,
			url:    story.URL,
			author: story.By,
			source: "HackerNews",
		}
		stories = append(stories, newStory)
	}

	// Finally, after the loop completes, we'll just return the stories
	return stories
}

// We need a function to get stories from Reddit, as well.
func newRedditStories() []Story {
	// Again, we need a buffer to hold stories
	var stories []Story
	// First we decide on some options. We want the most recent posts
	sort := geddit.PopularitySort(geddit.NewSubmissions)
	// and we'll use the default listing options
	var listingOptions geddit.ListingOptions
	// Now we can call the Subreddit Submissions method to get the submissions
	submissions, err := redditSession.SubredditSubmissions("programming", sort, listingOptions)
	// In case of an error, we'll print it and return nil
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// Unlike HackerNews, this single network operation will give us all the data we need.
	for _, s := range submissions {
		// As we did for HackerNews, we'll create new stories and add them
		// to the list
		newStory := Story{
			title:  s.Title,
			url:    s.URL,
			author: s.Author,
			source: "Reddit /r/programming",
		}
		stories = append(stories, newStory)
	}

	// Then finally we'll return the list
	return stories
}

// Now, in the main function, we can simply call these two functions
func main() {
	// We place each set of stories in a new buffer
	hnStories := newHnStories()
	redditStories := newRedditStories()
	// And we need a buffer to contain all stories
	var stories []Story

	// Now we check that each source actually did return some stories
	if hnStories != nil {
		// If so, we'll append those to the list
		stories = append(stories, hnStories...)
	}
	if redditStories != nil {
		stories = append(stories, redditStories...)
	}

	// Now let's write these stories to a file, stories.txt
	// First we open the file
	file, err := os.Create("stories.txt")
	// If there's a problem opening the file, abort
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// And we'll defer closing the file
	defer file.Close()
	// Now we'll write the stories out to the file
	for _, s := range stories {
		fmt.Fprintf(file, "%s: %s\nby %s on %s\n\n", s.title, s.url, s.author, s.source)
	}

	// Finally, we'll print out all the stories we received
	for _, s := range stories {
		fmt.Printf("%s: %s\nby %s on %s\n\n", s.title, s.url, s.author, s.source)
	}
}

// Running the program with go run main.go, you should see that it takes a
// moment but eventually does print the list of stories. If we check stories.txt,
// all the stories will be there.
// ST
// Now, it's pretty clear that this implementation could be improved. In the next video,
// we'll talk about exactly how.
