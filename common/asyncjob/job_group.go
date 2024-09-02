package asyncjob

import (
	"context"
	"log"
	"sync"
)

type group struct {
	jobs         []Job
	isConcurrent bool
	wg           *sync.WaitGroup
}

func NewGroup(isConcurrent bool, jobs ...Job) *group {
	g := &group{
		isConcurrent: isConcurrent,
		jobs:         jobs,
		wg:           new(sync.WaitGroup),
	}
	return g
}

func (g *group) Run(ctx context.Context) error {
	if g.isConcurrent {
		g.wg.Add(len(g.jobs))
	}
	errCh := make(chan error, len(g.jobs))

	for i := range g.jobs {
		job := g.jobs[i]
		if g.isConcurrent {
			go func(aj Job) {
				errCh <- g.runJob(ctx, aj)
				g.wg.Done()
			}(job)
		} else {
			err := g.runJob(ctx, job)
			errCh <- err
		}
	}

	if g.isConcurrent {
		g.wg.Wait()
	}

	for i := 1; i <= len(g.jobs); i++ {
		if v := <-errCh; v != nil {
			return v
		}
	}
	return nil
}

// Retry if needed
func (g *group) runJob(ctx context.Context, j Job) error {
	if err := j.Execute(ctx); err != nil {
		for {
			log.Println(err)
			if j.State() == StateRetryFailed {
				return err
			}
			if j.Retry(ctx) == nil {
				return nil
			}
		}
	}
	return nil
}
