package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	// Create an HTTP server that listens on port 8000
	http.ListenAndServe(":8000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// This prints to STDOUT to show that processing has started
		fmt.Fprint(os.Stdout, "started processing request\n")
		// We use select to execute a peice of code depending on which channel receives a message first
		select {
		case <-time.After(2 * time.Second):
			// We use this section to simulate some useful work
			// If we receive a message after 2 seconds
			// that means the request has been processed
			// We then write this as the response
			w.Write([]byte("request processed"))
		case <-ctx.Done():
			// If the request gets cancelled before 2 seconds, log it to STDERR
			fmt.Fprint(os.Stderr, "request cancelled\n")
		}
	}))
}
