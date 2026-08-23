package protocol

import "testing"

func TestClassifyUnknownErrorAsTransient(t *testing.T) {
	if got := Classify(assertionError{}); got != FailureTransient {
		t.Fatalf("class = %s", got)
	}
}

type assertionError struct{}

func (assertionError) Error() string { return "temporary" }
