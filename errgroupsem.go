package errgroupsem

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

type ErrGroupSem struct {
	cancel func()

	wg sync.WaitGroup

	errOnce sync.Once
	err     error
	sem     *semaphore.Weighted
}

func WithContext(ctx context.Context, limit int) (*ErrGroupSem, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &ErrGroupSem{cancel: cancel, sem: semaphore.NewWeighted(int64(limit))}, ctx
}

func (g *ErrGroupSem) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}

func (g *ErrGroupSem) markFailed(err error) {
	g.errOnce.Do(func() {
		g.err = err
		if g.cancel != nil {
			g.cancel()
		}
	})
}

func (g *ErrGroupSem) Go(ctx context.Context, f func() error) {
	if err := g.sem.Acquire(ctx, 1); err != nil {
		g.markFailed(err)
		return
	}
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()
		defer g.sem.Release(1)

		if err := f(); err != nil {
			g.markFailed(err)
		}
	}()
}
