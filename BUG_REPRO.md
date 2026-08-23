# Bug Reproduction

## Bug

Rollout handles and temporary leases are not cleaned up correctly when a device rejects a batch. Primary failures can also be hidden by transaction cleanup.

## Trigger

Run a multi-device rollout where one device rejects the configuration and cleanup encounters a busy temporary lease.

## Observed Errors

```text
maximum open handles = 3
err=<nil> tx=&{committed:true rolledBack:false}
expected rejection
err=busy calls=[bad]
```
