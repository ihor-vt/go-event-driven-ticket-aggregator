# Project: Graceful Shutdown

Any production-grade application should handle shutdowns gracefully.
The process should wait for all requests to finish before exiting.

Your handlers should be ready for the service getting suddenly killed anyway: hardware failures or power outages can happen at any time.
One way to mitigate losing data this way is using database transactions.
Still, it's a good idea to let running requests finish. This way, your users don't notice

Shutting down the Router is easy: All you need to do is pass a context to the `Run` method.
Once the context is canceled, the Router stops accepting new requests and waits for the
running ones to finish.

To detect the application receiving an interrupt signal, use `signal.NotifyContext`, like this:

```go
ctx := context.Background()
ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
defer cancel()

err := router.Run(ctx)
if err != nil {
	panic(err)
}

<-ctx.Done()
```

The interrupt signal cancels the `ctx`, and the Router stops accepting new messages.
It waits for the running handlers to finish and then returns without errors.

Often, your applications will have more long-running goroutines ("daemons") besides the Router that need to be shut down gracefully.

You can use the `errgroup` package (`golang.org/x/sync/errgroup`), that allows running multiple goroutines and waiting for them to finish.
(The `golang.org/x/` packages don't have the same API stability guarantee as the standard library, but it's good enough for our use case.) 

```go
ctx := context.Background()
ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
defer cancel()

g, ctx := errgroup.WithContext(ctx)

g.Go(func() error {
	return router.Run(ctx)
})

g.Go(func() error {
	err := e.Start(":8080")
	if err != nil && !errors.Is(err, stdHTTP.ErrServerClosed) {
		return err
	}
	
	return nil
})

g.Go(func() error {
	// Shut down the HTTP server
    <-ctx.Done()
    return e.Shutdown(ctx)
})

// Will block until all goroutines finish
err := g.Wait()
if err != nil {
    panic(err)
}
```

Using `errgroup.WithContext` creates a new context that is canceled if any of the goroutines return an error or if the original context is canceled.
Once the Router or the HTTP server stops, the context gets canceled and other daemons are notified to shut down.

{{tip}}

If you need a refresher on how the `Context` works, the `ctx.Done()` channel is closed when the context is canceled.
So waiting for the context to be canceled is as simple as:

```go
<-ctx.Done()
```

{{endtip}}

To summarize how this works:

1. Create a new context, and pass it to `signal.NotifyContext`. The incoming interrupt signal will cancel the context.
2. Create a new `errgroup` and pass the context to it.
3. Start the Router in a new goroutine. It will stop accepting new requests once the context is canceled.
4. Start the HTTP server in a new goroutine.
5. Start a goroutine that will shut down the HTTP server once the context is canceled.
6. Wait for all goroutines to finish.

## Exercise

Exercise path: ./project

**Introduce graceful Router shutdown in your project using `errgroup` and `signal.NotifyContext`.**

Remember to add `golang.org/x/sync/errgroup` to your dependencies.
