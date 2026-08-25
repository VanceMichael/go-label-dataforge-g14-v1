package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/domain"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/repository"
	"sync"
)

type BatchResult struct {
	ID      string
	Success bool
	Error   string
}
type BatchRunner struct {
	Store    repository.Store
	Parallel int
}

func (b BatchRunner) Run(ctx context.Context, items []domain.Authorization) ([]BatchResult, error) {
	if len(items) == 0 {
		return []BatchResult{}, nil
	}
	workers := b.Parallel
	if workers < 1 {
		workers = 1
	}
	if workers > len(items) {
		workers = len(items)
	}
	jobs := make(chan domain.Authorization)
	out := make(chan BatchResult, len(items))
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for a := range jobs {
				if e := ctx.Err(); e != nil {
					out <- BatchResult{ID: a.ID, Error: e.Error()}
					continue
				}
				e := b.Store.WithTx(ctx, func(tx *sql.Tx) error { return b.Store.CreateAuthorizationTx(ctx, tx, a) })
				if e != nil {
					out <- BatchResult{ID: a.ID, Error: e.Error()}
				} else {
					out <- BatchResult{ID: a.ID, Success: true}
				}
			}
		}()
	}
	go func() {
		defer close(jobs)
		for _, a := range items {
			select {
			case jobs <- a:
			case <-ctx.Done():
				return
			}
		}
	}()
	wg.Wait()
	close(out)
	res := make([]BatchResult, 0, len(items))
	for x := range out {
		res = append(res, x)
	}
	return res, nil
}
func Summarize(results []BatchResult) (ok, failed int, detail string) {
	for _, r := range results {
		if r.Success {
			ok++
		} else {
			failed++
			detail += fmt.Sprintf("%s:%s;", r.ID, r.Error)
		}
	}
	return
}
