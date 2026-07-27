# Wave 1: OSS Architecture

## Findings
- Six repositories were cloned and pinned, but the worker returned an incomplete synthesis.
- Recurrent patterns identified: adapter contracts, explicit capability matrices, internal ID assignment, provider-ID registries, normalized event schemas, conformance/replay tests, opt-in persistence.
- Strong candidates: All-Chat (internal UUID registry + contract tests), Virta (capability matrix + ErrUnsupported + ULID), Crossfeed (adapter/hub/persistence), Multibot (tenant/provider split), Unified-chat (typed connector model).

## EXPAND
- Produce SHA-pinned comparative matrix with source permalinks.
- Counter-search whether internal ID generation breaks provider reply/delete semantics.
- Extract Twir-applicable conformance-test design.
