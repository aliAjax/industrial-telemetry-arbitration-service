package protocol

import (
	"errors"
	"fmt"
	"testing"
)

func TestDecodePreservesChecksumCause(t *testing.T) {
	_, err := Decode([]byte{1})
	if !errors.Is(err, ErrMalformed) {
		t.Fatalf("malformed cause lost: %v", err)
	}
}

func TestRepositoryPreservesDeviceCause(t *testing.T) {
	err := NewRepository().Store(Frame{})
	if !errors.Is(err, ErrDeviceMissing) {
		t.Fatalf("device cause lost: %v", err)
	}
}

func TestClassifierRecognizesProtocolCauses(t *testing.T) {
	err := fmt.Errorf("ingest layer: %w", ErrChecksum)
	if got := Classify(err); got != FailurePermanent {
		t.Fatalf("class = %s", got)
	}
}

func TestHandlerMapsPermanentFrameErrors(t *testing.T) {
	response := NewHandler(NewRepository()).Ingest([]byte{9, 0, 0, 0, 1, 0, 8})
	if response.Status != 422 || response.Retry {
		t.Fatalf("response = %#v", response)
	}
}
