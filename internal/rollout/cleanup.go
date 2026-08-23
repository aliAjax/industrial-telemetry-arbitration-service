package rollout

import "errors"

type Releaser interface{ Release(string) error }

type Cleanup struct{ releaser Releaser }

func NewCleanup(releaser Releaser) *Cleanup { return &Cleanup{releaser: releaser} }

func (c *Cleanup) Release(ids []string) error {
	var errs []error
	for _, id := range ids {
		if err := c.releaser.Release(id); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
