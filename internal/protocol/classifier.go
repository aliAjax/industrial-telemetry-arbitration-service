package protocol

type FailureClass string

const (
	FailurePermanent FailureClass = "permanent"
	FailureMissing   FailureClass = "missing"
	FailureTransient FailureClass = "transient"
)

func Classify(err error) FailureClass {
	switch {
	case err == ErrChecksum, err == ErrVersion, err == ErrMalformed:
		return FailurePermanent
	case err == ErrDeviceMissing:
		return FailureMissing
	default:
		return FailureTransient
	}
}

func Retryable(err error) bool { return Classify(err) == FailureTransient }
