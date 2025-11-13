package app

import (
	"context"
	"time"

	"github.com/flohansen/coffee-table/pkg/logging"
	"golang.org/x/sync/errgroup"
)

var _ context.Context = &Stack{}

type Stack struct {
	ctx context.Context
	g   *errgroup.Group
}

// Deadline implements context.Context.
func (s *Stack) Deadline() (deadline time.Time, ok bool) {
	return s.ctx.Deadline()
}

// Done implements context.Context.
func (s *Stack) Done() <-chan struct{} {
	return s.ctx.Done()
}

// Err implements context.Context.
func (s *Stack) Err() error {
	return s.ctx.Err()
}

// Value implements context.Context.
func (s *Stack) Value(key any) any {
	return s.ctx.Value(key)
}

func (s *Stack) Go(f func() error) {
	s.g.Go(f)
}

func (s *Stack) Wait() error {
	return s.g.Wait()
}

func NewStack(opts ...StackOption) *Stack {
	g, ctx := errgroup.WithContext(SignalContext())
	s := &Stack{
		ctx: ctx,
		g:   g,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type StackOption func(*Stack)

func WithLogger(logger logging.Logger) StackOption {
	return func(s *Stack) {
		s.ctx = logging.WithContext(s.ctx, logger)
	}
}
