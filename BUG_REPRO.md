# Bug Reproduction

## What happens

Clock election shutdown and event subscription operations are not safe when probe goroutines overlap. Candidate labels and recorded history can also change after callers modify their original values.

## How to trigger it

Run the clockmesh concurrency and snapshot checks while two goroutines close an election or publish and close a subscriber at the same time. Record a candidate, then mutate the input or a returned snapshot before reading the history again.

## Observed error

The race-enabled run reports `WARNING: DATA RACE` in `internal/clockmesh.(*Election).Close`, including concurrent reads and writes of the closed state. The same run reports history labels changing to `snapshot-changed` and `original label = "changed"`.
