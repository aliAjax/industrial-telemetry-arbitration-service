package failover

import "testing"

func TestActiveCanDegrade(t *testing.T) {
	if !StateActive.CanTransition(StateDegraded) {
		t.Fatal("active should degrade")
	}
}
