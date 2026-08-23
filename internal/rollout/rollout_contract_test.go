package rollout

import (
	"errors"
	"testing"
)

type trackingHandle struct{ tracker *handleTracker }
type handleTracker struct{ open, max int }

func (h *trackingHandle) Apply() error { return nil }
func (h *trackingHandle) Close() error { h.tracker.open--; return nil }

func TestBatchReleasesEachHandleImmediately(t *testing.T) {
	tracker := &handleTracker{}
	batch := Batch{Open: func(string) (Handle, error) {
		tracker.open++
		if tracker.open > tracker.max {
			tracker.max = tracker.open
		}
		return &trackingHandle{tracker}, nil
	}}
	if err := batch.Apply([]string{"a", "b", "c"}); err != nil {
		t.Fatal(err)
	}
	if tracker.max != 1 {
		t.Fatalf("maximum open handles = %d", tracker.max)
	}
}

type fakeTx struct{ committed, rolledBack bool }

func (t *fakeTx) Commit() error   { t.committed = true; return nil }
func (t *fakeTx) Rollback() error { t.rolledBack = true; return nil }

func TestRepositoryPreservesPrimaryFailure(t *testing.T) {
	primary := errors.New("station rejected")
	tx := &fakeTx{}
	repository := Repository{Begin: func() (Tx, error) { return tx, nil }}
	err := repository.WithTx(func(Tx) error { return primary })
	if !errors.Is(err, primary) || !tx.rolledBack || tx.committed {
		t.Fatalf("err=%v tx=%+v", err, tx)
	}
}

type fakeLeases struct{ removed []string }

func (*fakeLeases) Create(string) error      { return nil }
func (*fakeLeases) Promote(string) error     { return errors.New("promote failed") }
func (l *fakeLeases) Remove(id string) error { l.removed = append(l.removed, id); return nil }

func TestServiceRollsBackRejectedRollout(t *testing.T) {
	leases := &fakeLeases{}
	repository := Repository{Begin: func() (Tx, error) { return &fakeTx{}, nil }}
	service := NewService(repository, leases, Batch{Open: func(string) (Handle, error) { return &trackingHandle{&handleTracker{}}, nil }})
	if err := service.Rollout("rollout-1", []string{"a"}); err == nil {
		t.Fatal("expected rejection")
	}
	if len(leases.removed) != 1 {
		t.Fatalf("removed = %v", leases.removed)
	}
}

type releaseRecorder struct{ calls []string }

func (r *releaseRecorder) Release(id string) error {
	r.calls = append(r.calls, id)
	if id == "bad" {
		return errors.New("busy")
	}
	return nil
}

func TestCleanupRemovesTemporaryLeases(t *testing.T) {
	recorder := &releaseRecorder{}
	err := NewCleanup(recorder).Release([]string{"bad", "later"})
	if err == nil || len(recorder.calls) != 2 {
		t.Fatalf("err=%v calls=%v", err, recorder.calls)
	}
}
