# go-context
In Go servers, each incoming request is handled in its own goroutine. Request handlers often start additional goroutines to access backends such as databases and RPC services. The set of goroutines working on a request typically needs access to request-specific values such as the identity of the end user, authorization tokens, and the request's deadline. When a request is canceled or times out, all the goroutines working on that request should exit quickly so the system can reclaim any resources they are using.

Context makes it easy to pass request-scoped values, cancelation signals, and deadlines across API boundaries to all the goroutines involved in handling a request.

It is simply request object that flows through your program and interact (sending signal like cancelation or timeout).

```
type Context interface {
  Deadline() (deadline time.Time, ok bool)
  Done() <-chan struct{}
  Err() error
  Value(key interface{}) interface{}
}
```

**Listening for the cancellation event**

The Context type provides a Done() method, which returns a channel that receives an empty *struct{}* type everytime the context receives a cancellation event. Listening for a cancellation event is as easy as waiting for *<- ctx.Done()*.

For example, lets consider an HTTP server that takes two seconds to process an event. If the request gets cancelled before that, we want to return immediately

```
func main() {
	// Create an HTTP server that listens on port 8000
	http.ListenAndServe(":8000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// This prints to STDOUT to show that processing has started
		fmt.Fprint(os.Stdout, "processing request\n")
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

```

References: 
- https://ekocaman.com/go-context-c44d681da2e8
- https://blog.golang.org/context
- https://www.sohamkamani.com/blog/golang/2018-06-17-golang-using-context-cancellation/
