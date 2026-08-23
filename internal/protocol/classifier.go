package protocol

import "errors"

type FailureClass string

const (
	FailurePermanent FailureClass = "permanent"
	FailureMissing   FailureClass = "missing"
	FailureTransient FailureClass = "transient"
)

func Classify(err error) FailureClass {
	switch {
	case errors.Is(err, ErrChecksum), errors.Is(err, ErrVersion), errors.Is(err, ErrMalformed):
		return FailurePermanent
	case errors.Is(err, ErrDeviceMissing):
		return FailureMissing
	default:
		return FailureTransient
	}
}

func Retryable(err error) bool { return Classify(err) == FailureTransient }
