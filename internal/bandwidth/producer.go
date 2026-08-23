package bandwidth

import "errors"

var ErrLinkRejected = errors.New("bandwidth link rejected")

type Allocation struct {
	LinkID string
	Units  int
}

type Producer struct {
	LinkID string
	Units  []int
	Reject bool
}

func (p Producer) Stream(out chan<- Allocation, errs chan<- error) {
	defer close(out)
	if p.Reject {
		errs <- ErrLinkRejected
		return
	}
	for _, units := range p.Units {
		out <- Allocation{LinkID: p.LinkID, Units: units}
	}
}
