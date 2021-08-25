# errgroupsem

This package is an extended errgroup package to fix the lack of CPU bound support when errgroup is used in conjuction with semaphore.

The semaphore package could not be used with errgroup since the error in semaphore is internally managed by errgroup.

See: https://github.com/golang/go/issues/27837 (create another package https://github.com/golang/go/issues/27837#issuecomment-633659652)

## How to install
```bash
go get github.com/cboudereau/errgroupsem
```

## Example
See the unit tests for little and channel based demo

```go
// Fan-in / Fan-out example
ctx := context.Background()

numCPU := runtime.NumCPU()

// one main errgroupsem g instance
g, ctx := errgroupsem.WithContext(ctx, numCPU)

producer := func(size int) <-chan string {
	output := make(chan string)
	g.Go(ctx, func() error {
		defer close(output)
		// another degree of parallelism with another errgroupsem instance
		wg, ctx := errgroupsem.WithContext(ctx, numCPU)

		for i := 0; i < size; i++ {
			i := i //golang closure issue
			wg.Go(ctx, func() error {
				s := int64(rand.Intn(100))
				time.Sleep(time.Millisecond * time.Duration(s))
				output <- fmt.Sprintf("%v/%vms", i, s)
				return nil
			})
		}
		return wg.Wait()
	})

	return output
}

consumer := func(input <-chan string) {
	g.Go(ctx, func() error {
		for x := range input {
			fmt.Println("consumer", x)
		}
		return nil
	})
}

consumer(producer(100))

g.Wait()
```