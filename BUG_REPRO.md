# Bug Reproduction

## Bug

The failover state machine omits a legal stabilization transition, permits an illegal recovery jump, and leaves recovered links in the wrong query state.

## Trigger

Move a degraded link into stabilization, report a successful probe, then inspect its terminal state and in-progress query results. Also attempt a failed-to-active jump.

## Observed Errors

```text
stabilizing cannot reach active
error = <nil>
state = stabilizing
in progress = []failover.Link{{ID:"b", State:"degraded"}}
```
