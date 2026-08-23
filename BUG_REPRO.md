# Bug Reproduction

## Bug

Protocol errors lose their wrapped causes, so permanent frame failures and missing devices are classified as transient and retried.

## Trigger

Decode a malformed or checksum-invalid frame, or ingest a frame from an unknown device, then classify the returned error through the handler.

## Observed Errors

```text
malformed cause lost: decode header: protocol frame malformed
device cause lost: store device: protocol device missing
response = protocol.Response{Status:503, Retry:true, Class:"transient"}
```
