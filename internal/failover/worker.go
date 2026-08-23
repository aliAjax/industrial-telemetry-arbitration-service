package failover

import "context"

type Probe interface {
	Stable(context.Context, string) (bool, error)
}

type Worker struct{ probe Probe }

func NewWorker(probe Probe) *Worker { return &Worker{probe: probe} }

func (w *Worker) Recover(ctx context.Context, link string, machine *Machine) error {
	if err := machine.Advance(StateStabilizing); err != nil {
		return err
	}
	stable, err := w.probe.Stable(ctx, link)
	if err != nil {
		_ = machine.Advance(StateDegraded)
		return err
	}
	if !stable {
		return machine.Advance(StateDegraded)
	}
	return nil
}
