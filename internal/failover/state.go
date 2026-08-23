package failover

type State string

const (
	StateActive      State = "active"
	StateDegraded    State = "degraded"
	StateStabilizing State = "stabilizing"
	StateFailed      State = "failed"
)

func (s State) CanTransition(next State) bool {
	allowed := map[State]map[State]bool{
		StateActive:   {StateDegraded: true, StateFailed: true},
		StateDegraded: {StateStabilizing: true, StateFailed: true},
		StateStabilizing: {
			StateActive:   true,
			StateDegraded: true,
			StateFailed:   true,
		},
		StateFailed: {StateStabilizing: true},
	}
	return allowed[s][next]
}

func (s State) InProgress() bool { return s == StateDegraded || s == StateStabilizing }
