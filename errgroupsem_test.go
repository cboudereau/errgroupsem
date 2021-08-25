package errgroupsem_test

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"

	"github.com/cboudereau/errgroupsem"
	"github.com/stretchr/testify/require"
)

func TestErrGroupSem(t *testing.T) {
	g, ctx := errgroupsem.WithContext(context.Background(), 10)
	start := time.Now()
	for i := 0; i < 100; i++ {
		g.Go(ctx, func() error {
			time.Sleep(1 * time.Millisecond)
			return nil
		})
	}
	err := g.Wait()
	elapsed := time.Since(start)
	require.NoError(t, err)
	require.LessOrEqual(t, elapsed, 20*time.Millisecond)
	require.GreaterOrEqual(t, elapsed, 5*time.Millisecond)
}

func TestFanInFanOutExample(t *testing.T) {
	ctx := context.Background()

	numCPU := runtime.NumCPU()
	g, ctx := errgroupsem.WithContext(ctx, numCPU)

	producer := func(size int) <-chan string {
		output := make(chan string)
		g.Go(ctx, func() error {
			defer close(output)
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
}
