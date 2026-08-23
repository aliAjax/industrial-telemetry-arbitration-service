package waveform

import "testing"

func TestFilterStableKeepsOrder(t *testing.T) {
	in := []Sample{{Sequence: 1, Quality: 1}, {Sequence: 2, Quality: 0}, {Sequence: 3, Quality: 1}}
	out := FilterStable(in, 1)
	if len(out) != 2 || out[0].Sequence != 1 || out[1].Sequence != 3 {
		t.Fatalf("unexpected filter: %#v", out)
	}
}
