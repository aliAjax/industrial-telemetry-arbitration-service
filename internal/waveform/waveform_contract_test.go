package waveform

import "testing"

func TestFilterDoesNotRewriteInputWindow(t *testing.T) {
	input := []Sample{{Sequence: 1, Quality: 0}, {Sequence: 2, Quality: 2}}
	before := Clone(input)
	_ = FilterStable(input, 1)
	if input[0] != before[0] || input[1] != before[1] {
		t.Fatalf("input changed: %#v", input)
	}
}

func TestArchiveOwnsStoredSamples(t *testing.T) {
	archive := NewArchive()
	input := []Sample{{Sequence: 1, Value: 10}}
	archive.Store("window", input)
	input[0].Value = 99
	if got := archive.Load("window")[0].Value; got != 10 {
		t.Fatalf("archived value = %v", got)
	}
}

func TestAssemblerUsesReturnedSlice(t *testing.T) {
	assembler := NewAssembler(4)
	result := assembler.Append([]Sample{{Sequence: 1, Value: 10}})
	result[0].Value = 99
	if got := assembler.Samples()[0].Value; got != 10 {
		t.Fatalf("assembler retained caller mutation: %v", got)
	}
}

func TestCacheReturnsIndependentWindow(t *testing.T) {
	cache := NewCache()
	cache.Put("window", []Sample{{Sequence: 1, Value: 10}})
	snapshot := cache.Snapshot()
	snapshot["window"][0].Value = 99
	if got := cache.Snapshot()["window"][0].Value; got != 10 {
		t.Fatalf("cache value = %v", got)
	}
}
