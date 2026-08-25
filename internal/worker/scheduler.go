package worker

import (
	"context"
	"sync"
	"time"
)

type Scheduler struct {
	mu      sync.Mutex
	running bool
	cancel  context.CancelFunc
}

func (s *Scheduler) Start(parent context.Context, fn func(context.Context)) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running {
		return false
	}
	ctx, c := context.WithCancel(parent)
	s.cancel = c
	s.running = true
	go func() { fn(ctx); s.mu.Lock(); s.running = false; s.mu.Unlock() }()
	return true
}
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}
	s.running = false
}
func (s *Scheduler) Running() bool { s.mu.Lock(); defer s.mu.Unlock(); return s.running }
func Wait(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
