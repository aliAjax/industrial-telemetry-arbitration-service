# Bug Reproduction

## What happens

Concurrent acknowledgement fanout can lose slow gateway results, close result channels before workers finish, panic while publishing, and let later metadata edits rewrite stored history.

## How to trigger it

Start several fanout workers with one deliberately delayed processor, collect the result channel to completion, and race a publisher against close. Then append an acknowledgement with metadata, mutate the input and a returned snapshot, and read the store again.

## Observed error

The race-enabled run reports `panic: send on closed channel` from `internal/ackgroup.(*Publisher).Publish` at `publisher.go:22`. It also reports `collected 1 acknowledgements` for a three-item input and `stored metadata = "snapshot-changed"` after a snapshot was edited.
