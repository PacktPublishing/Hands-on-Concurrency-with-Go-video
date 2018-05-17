package main

// In this video, we'll implement a search engine, but without any actual connection to Reddit or HackerNews.
// ST
// Pull up your editor, and let's get started!

// We need to import net/http which provides a task concurrent web server
import (
	"fmt"
	"net/http"
	"strings"
)

// We'll use the same Story struct as we have used in the past.

type Story struct {
	title  string
	url    string
	author string
	source string
}

// We'll set a pre-defined list of known stories - just a few.
var stories []Story

func init() {
	stories = append(stories,
		Story{"Go Language Stuff", "http://example.com", "leotindall", "fake"},
		Story{"Python Language Stuff", "http://example.com", "leotindall", "fake"},
		Story{"Rust Language Stuff", "http://example.com", "leotindall", "fake"},
		Story{"Programming Culture Stuff", "http://example.com", "leotindall", "fake"},
		Story{"Go Performance Stuff", "http://example.com", "leotindall", "fake"})
}

// Now we'll create a simple function to search all stories.
// This could be made data parallel, but in this simple example we won't do that.
func searchStories(query string) []Story {
	var foundStories []Story
	for _, story := range stories {
		if strings.Contains(strings.ToUpper(story.title), strings.ToUpper(query)) {
			foundStories = append(foundStories, story)
		}
	}
	return foundStories
}

// Now we can create a function that connects the search function to the route.
// It takes some specialized types from the http module.
// ResponseWriter gives us access to the client response, while Request gives us all the information
// the client provided to the server.
func search(w http.ResponseWriter, r *http.Request) {
	// This handler is only interested in the query value, which we'll call q.
	query := r.FormValue("q")
	// If there was no query, we can't do anything, so we'll return an error.
	if query == "" {
		http.Error(w, "Search parameter q is required to search.", http.StatusNotAcceptable)
		return
	}

	// If there was no error, we write the opening of an HTML document
	w.Write([]byte("<html><body>"))
	// We'll search for stories and get the list
	s := searchStories(query)
	// If there were no results we need to report that
	if len(s) == 0 {
		w.Write([]byte(fmt.Sprintf("No results for query '%s'.\n<br>", r.FormValue("q"))))
	} else {
		// Otherwise we'll loop over the stories and write them out as links
		for _, story := range s {
			w.Write([]byte(fmt.Sprintf("<a href='%s'>%s</a><br>by %s on %s<br><br>", story.url, story.title, story.author, story.source)))
		}
	}

	// Either way we need to write a back button and close the html document
	w.Write([]byte("<a href='../'>Back</a>"))
	w.Write([]byte("</body></html>"))

}

// Now we need a function that lists the first ten known stories
func topTen(w http.ResponseWriter, r *http.Request) {
	// As before we'll write the html document opening
	w.Write([]byte("<html><body>"))
	// Now we'll write a form that lets us search with the prior route
	form := "<form action='search' method='post'>Search: <input type='text' name='q'> <input type='submit'></form>\n"
	w.Write([]byte(form))
	// Next, we just need to loop over the stories we have, starting at the end and going for ten.
	for i := len(stories) - 1; i >= 0 && len(stories)-i < 10; i-- {
		story := stories[i]
		w.Write([]byte(fmt.Sprintf("<a href='%s'>%s</a><br>by %s on %s<br><br>", story.url, story.title, story.author, story.source)))
	}
	// Finally, we'll close the HTML document.
	w.Write([]byte("</body></html>"))
}

// The main function is really easy.
func main() {
	// We just add the two handlers we created
	http.HandleFunc("/", topTen)
	http.HandleFunc("/search", search)
	// And now we run the actual server.
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

// Congratulations! We've created a task parallel web server.
// We can pull it up in a web browser at localhost:8080 and see our top stories
// and we can also search for different keywords, like Rust or Go
// ST
// In the next video, we'll add in the actual Reddit and HackerNews client.
