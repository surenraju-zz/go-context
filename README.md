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

**Emitting a cancellation event**

If you have an operation that could be cancelled, you will have to emit a cancellation event through the context. This can be done using the *WithCancel* function in the context package, which returns a context object, and a function. This function takes no arguments, and does not return anything, and is called when you want to cancel the context.

Consider the case of 2 dependent operations. Here, “dependent” means if one fails, it doesn’t make sense for the other to complete. In this case, if we get to know early on that one of the operations failed, we would like to cancel all dependent operations.



**Creating context**

The context package allows creating and deriving context in following ways:

**context.Background() ctx Context**

This function returns an empty context. This should be only used at a high level (in main or the top level request handler). This can be used to derive other contexts that we discuss later.

```
ctx, cancel := context.Background()
```

**context.TODO() ctx Context**

This function also creates an empty context. This should also be only used at a high level or when you are not sure what context to use or if the function has not been updated to receive a context. Which means you (or the maintainer) plans to add context to the function in future.

```
ctx, cancel := context.TODO()
```

**context.WithValue(parent Context, key, val interface{}) (ctx Context, cancel CancelFunc)**

This function takes in a context and returns a derived context where value val is associated with key and flows through the context tree with the context. This means that once you get a context with value, any context that derives from this gets this value. It is not recommended to pass in critical parameters using context value, instead, functions should accept those values in the signature making it explicit.
```
ctx := context.WithValue(context.Background(), key, "test")
```

**context.WithCancel(parent Context) (ctx Context, cancel CancelFunc)**

This is where it starts to get a little interesting. This function creates a new context derived from the parent context that is passed in. The parent can be a background context or a context that was passed into the function.

This returns a derived context and the cancel function. Only the function that creates this should call the cancel function to cancel this context. You can pass around the cancel function if you wanted to, but, that is highly not recommended. This can lead to the invoker of cancel not realizing what the downstream impact of canceling the context may be. There may be other contexts that are derived from this which may cause the program to behave in an unexpected fashion. In short, NEVER pass around the cancel function.

```
ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2 * time.Second))
```

**context.WithDeadline(parent Context, d time.Time) (ctx Context, cancel CancelFunc)**

This function returns a derived context from its parent that gets cancelled when the deadline exceeds or cancel function is called. For example, you can create a context that will automatically get canceled at a certain time in future and pass that around in child functions. When that context gets canceled because of deadline running out, all the functions that got the context get notified to stop work and return.

```
ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2 * time.Second))
```

**context.WithTimeout(parent Context, timeout time.Duration) (ctx Context, cancel CancelFunc)**

This function is similar to context.WithDeadline. The difference is that it takes in time duration as an input instead of the time object. This function returns a derived context that gets canceled if the cancel function is called or the timeout duration is exceeded.

```
ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2 * time.Second))
```

References: 
- https://ekocaman.com/go-context-c44d681da2e8
- https://blog.golang.org/context
- https://www.sohamkamani.com/blog/golang/2018-06-17-golang-using-context-cancellation/
