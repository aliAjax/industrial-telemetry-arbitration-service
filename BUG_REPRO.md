# Bug Reproduction

## Bug

Uplink request contexts are replaced or retained across requests. Cancellation is lost on one path and contaminates a later request on another.

## Trigger

Cancel a slow uplink request, then submit a fresh high-priority request through the same session. Also cancel a request while it is retrying.

## Observed Errors

```text
fresh request reused canceled context: context canceled
err=retry calls=3
worker discarded request deadline
```
