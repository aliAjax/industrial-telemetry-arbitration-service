package uplink

import (
	"context"
	"errors"
	"testing"
	"time"
)

type senderFunc func(context.Context, []byte) error

func (f senderFunc) Send(ctx context.Context, payload []byte) error { return f(ctx, payload) }

type transportFunc func(context.Context, []byte) error

func (f transportFunc) Transmit(ctx context.Context, payload []byte) error { return f(ctx, payload) }

func TestHandlePropagatesRequestCancellation(t *testing.T) {
	ctx := context.WithValue(context.Background(), "request-id", "fresh")
	called := false
	err := Handle(ctx, senderFunc(func(callCtx context.Context, _ []byte) error {
		called = true
		if callCtx.Value("request-id") != "fresh" {
			return errors.New("request context lost")
		}
		return context.Canceled
	}), Request{Payload: []byte("x")})
	if !errors.Is(err, context.Canceled) || !called {
		t.Fatalf("err=%v called=%v", err, called)
	}
}

func TestSessionDoesNotReuseCanceledContext(t *testing.T) {
	session := NewSession(transportFunc(func(ctx context.Context, _ []byte) error { return ctx.Err() }))
	first, cancel := context.WithCancel(context.Background())
	if err := session.Send(first, []byte("first")); err != nil {
		t.Fatal(err)
	}
	cancel()
	if err := session.Send(context.Background(), []byte("fresh")); err != nil {
		t.Fatalf("fresh request reused canceled context: %v", err)
	}
}

func TestRetryStopsAfterCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	calls := 0
	err := (Retryer{Attempts: 3, Delay: time.Millisecond}).Do(ctx, func(context.Context) error { calls++; return errors.New("retry") })
	if !errors.Is(err, context.Canceled) || calls != 0 {
		t.Fatalf("err=%v calls=%d", err, calls)
	}
}

func TestWorkerKeepsFreshRequestDeadline(t *testing.T) {
	deadlineSeen := false
	sender := senderFunc(func(ctx context.Context, _ []byte) error { _, deadlineSeen = ctx.Deadline(); return nil })
	worker := NewWorker(Retryer{Attempts: 1}, sender)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := worker.Process(ctx, []Request{{Payload: []byte("frame")}}); err != nil {
		t.Fatal(err)
	}
	if !deadlineSeen {
		t.Fatal("worker discarded request deadline")
	}
}
