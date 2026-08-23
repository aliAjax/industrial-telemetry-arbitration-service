# Industrial Telemetry Arbitration Service

This Go service coordinates telemetry uplink and downlink work across industrial
edge gateways. It owns link arbitration, protocol classification, adapter
validation, staged configuration rollout, waveform ownership, failover state,
and bandwidth allocation.

## Layout

- `cmd/telemetryd`: HTTP service entry point.
- `internal/arbitration`: terminal lease and subscriber lifecycle.
- `internal/uplink`: request-scoped uplink execution.
- `internal/protocol`: frame decoding and error classification.
- `internal/adapter`: optional vendor adapter construction and validation.
- `internal/rollout`: transactional station configuration rollout.
- `internal/waveform`: waveform filtering and snapshot ownership.
- `internal/failover`: link recovery state machine.
- `internal/bandwidth`: concurrent allocation collection.
- `internal/downlink`: request-scoped command dispatch.

## Run

```sh
go run ./cmd/telemetryd
curl http://127.0.0.1:18141/readyz
```

Set `TELEMETRY_ADDR` to override the listen address.

## Test

```sh
go test ./...
go test -race ./...
```
