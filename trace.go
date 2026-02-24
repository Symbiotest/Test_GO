package main

import (
	"context"
	"net/http"
	"sync"
)

type traceKey struct{}

type RequestTrace struct {
	mu    sync.Mutex
	steps []string
}

func newRequestTrace() *RequestTrace {
	return &RequestTrace{steps: make([]string, 0, 8)}
}

func (t *RequestTrace) Add(step string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.steps = append(t.steps, step)
}

func (t *RequestTrace) Steps() []string {
	t.mu.Lock()
	defer t.mu.Unlock()
	cp := make([]string, len(t.steps))
	copy(cp, t.steps)
	return cp
}

func withTrace(ctx context.Context) (context.Context, *RequestTrace) {
	if t, ok := ctx.Value(traceKey{}).(*RequestTrace); ok && t != nil {
		return ctx, t
	}
	t := newRequestTrace()
	return context.WithValue(ctx, traceKey{}, t), t
}

func ensureTrace(r *http.Request) *http.Request {
	ctx, _ := withTrace(r.Context())
	if ctx == r.Context() {
		return r
	}
	return r.WithContext(ctx)
}

func addTrace(ctx context.Context, step string) {
	if t, ok := ctx.Value(traceKey{}).(*RequestTrace); ok && t != nil {
		t.Add(step)
	}
}

func traceSteps(ctx context.Context) []string {
	if t, ok := ctx.Value(traceKey{}).(*RequestTrace); ok && t != nil {
		return t.Steps()
	}
	return nil
}
