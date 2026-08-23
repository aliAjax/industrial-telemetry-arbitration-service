# Bug Reproduction

## Bug

Waveform filters, archives, assemblers, and caches share caller-owned slice backing arrays. Later mutation rewrites previously stored samples.

## Trigger

Archive a waveform window, filter or append another frame, then mutate the caller's original or returned slice and read the archived/cache values.

## Observed Errors

```text
input changed: []waveform.Sample{...}
archived value = 99
assembler retained caller mutation: 99
cache value = 99
```
