package downlink

import (
	"context"
	"errors"
	"testing"
	"time"
)

type commandSenderFunc func(context.Context, Command) error

func (f commandSenderFunc) Send(ctx context.Context, command Command) error { return f(ctx, command) }

type wireFunc func(context.Context, Command) error

func (f wireFunc) Write(ctx context.Context, command Command) error { return f(ctx, command) }

func TestDownlinkHandlePreservesDeadline(t *testing.T) {
	seen := false
	sender := commandSenderFunc(func(ctx context.Context, _ Command) error { _, seen = ctx.Deadline(); return nil })
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := Handle(ctx, sender, Command{}); err != nil {
		t.Fatal(err)
	}
	if !seen {
		t.Fatal("downlink deadline discarded")
	}
}

func TestClientUsesCurrentCommandScope(t *testing.T) {
	client := NewClient(wireFunc(func(ctx context.Context, _ Command) error { return ctx.Err() }))
	first, cancel := context.WithCancel(context.Background())
	if err := client.Send(first, Command{}); err != nil {
		t.Fatal(err)
	}
	cancel()
	if err := client.Send(context.Background(), Command{}); err != nil {
		t.Fatalf("fresh command reused canceled context: %v", err)
	}
}

func TestBackoffStopsWhenCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	started := time.Now()
	err := (Backoff{Delay: 200 * time.Millisecond}).Wait(ctx)
	if !errors.Is(err, context.Canceled) || time.Since(started) > 100*time.Millisecond {
		t.Fatalf("err=%v elapsed=%s", err, time.Since(started))
	}
}

func TestCanceledDispatcherSkipsWire(t *testing.T) {
	ctx := context.WithValue(context.Background(), "request-id", "urgent")
	calls := 0
	dispatcher := NewDispatcher(commandSenderFunc(func(callCtx context.Context, _ Command) error {
		calls++
		if callCtx.Value("request-id") != "urgent" {
			return errors.New("request context lost")
		}
		return nil
	}), Backoff{}, 2)
	err := dispatcher.Dispatch(ctx, Command{})
	if err != nil || calls != 1 {
		t.Fatalf("err=%v calls=%d", err, calls)
	}
}
