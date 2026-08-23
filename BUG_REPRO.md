# Bug Reproduction

## Bug

Bandwidth producers, consumers, and the coordinator disagree about WaitGroup and channel ownership. A rejected link can hang collection, while other runs return early or drop errors.

## Trigger

Rebalance several links concurrently with one rejected producer, consume all result streams, and publish a buffered worker error under the race detector.

## Observed Errors

```text
rejected producer left output open
coordinator completed before producer output was consumed
allocations = []bandwidth.Allocation{{LinkID:"a", Units:0}}
buffered error was dropped
```
