# Wave 2 OSS Comparison

SHA-pinned references:

- `caesarakalaeii/all-chat@3831acf8f060ad31da7b0049177b2ba0e0a3cada`
- `elythi0n/virta@027d0779365793781dff9d6f53374fbd9f6b3d07`
- `maxjschiller/crossfeed@45749339700be2706bb5f91697fdc5ef6a66bd54`
- `Fixlation/Mergerino@dd1864960b42bf61eb1cb409f6563052b41aff5c`
- `NamoVize/universal-chat-aggregator@29f2b30933171394a8fb02c2b5487ef7fa25abee`

## Applicable Patterns

- Virta: centralized typed platform contract, explicit capabilities, `ErrUnsupported`, centralized outbound dispatch, and contract/golden/replay tests.
- All-Chat: fanout outcome reconciliation, Redis-backed replay/dedup, and bounded history.
- Crossfeed: one normalized message model, source lifecycle, persisted source/session state, and explicitly labeled degraded modes.
- Mergerino: snapshot regressions for provider parsing and normalization.

## Counterexample

- Universal Chat Aggregator uses string-keyed service maps and simulation stubs without typed contracts or capability enforcement. Twir should not move further toward string arrays and implicit fallbacks.
