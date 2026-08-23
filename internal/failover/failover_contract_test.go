package failover

import (
	"context"
	"errors"
	"testing"
)

type stableProbe bool

func (p stableProbe) Stable(context.Context, string) (bool, error) { return bool(p), nil }

func TestStateAllowsStabilizingToActive(t *testing.T) {
	if !StateStabilizing.CanTransition(StateActive) {
		t.Fatal("stabilizing cannot reach active")
	}
}

func TestMachineRejectsIllegalRecoveryJump(t *testing.T) {
	machine := NewMachine(StateDegraded)
	if err := machine.Advance(StateActive); !errors.Is(err, ErrIllegalTransition) {
		t.Fatalf("error = %v", err)
	}
}

func TestWorkerPublishesRecoveredTerminalState(t *testing.T) {
	machine := NewMachine(StateDegraded)
	if err := NewWorker(stableProbe(true)).Recover(context.Background(), "link", machine); err != nil {
		t.Fatal(err)
	}
	if got := machine.State(); got != StateActive {
		t.Fatalf("state = %s", got)
	}
}

func TestQueryIncludesStabilizingLinks(t *testing.T) {
	links := NewQuery([]Link{{ID: "a", State: StateStabilizing}, {ID: "b", State: StateDegraded}}).InProgress()
	if len(links) != 2 {
		t.Fatalf("in progress = %#v", links)
	}
}
