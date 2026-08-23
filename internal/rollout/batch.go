package rollout

import "fmt"

type Handle interface {
	Apply() error
	Close() error
}

type Batch struct{ Open func(string) (Handle, error) }

func (b Batch) Apply(stations []string) error {
	for _, station := range stations {
		if err := b.applyOne(station); err != nil {
			return err
		}
	}
	return nil
}

func (b Batch) applyOne(station string) (err error) {
	handle, err := b.Open(station)
	if err != nil {
		return fmt.Errorf("open station %s: %w", station, err)
	}
	defer func() {
		if closeErr := handle.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()
	if err := handle.Apply(); err != nil {
		return fmt.Errorf("apply station %s: %w", station, err)
	}
	return nil
}
