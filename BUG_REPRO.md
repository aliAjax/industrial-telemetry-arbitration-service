# Bug Reproduction

## Bug

Downlink request contexts are discarded or cached across commands. A canceled ordinary command contaminates a fresh emergency command, while canceled dispatches continue sending.

## Trigger

Cancel one downlink command, immediately send another command through the same client, and cancel a dispatcher during backoff.

## Observed Errors

```text
downlink deadline discarded
fresh command reused canceled context: context canceled
err=<nil> elapsed=202.037084ms
err=request context lost calls=2
```
