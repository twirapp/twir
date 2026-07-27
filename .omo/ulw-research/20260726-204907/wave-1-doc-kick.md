# Wave 1: Kick Contracts

## Findings
- Current API is OAuth 2.1/PKCE, REST on `api.kick.com`, signed webhooks, documented scopes.
- Chat write/delete is documented; no chat read/list endpoint was found.
- Reply uses optional `reply_to_message_id`; 429 is a documented error, not success.
- App-token event subscriptions require broadcaster ID; webhook failures can auto-unsubscribe after one day.
- Current livestream API is v2; v1 is deprecated.

## EXPAND
- Compare Twir token type/broadcaster ID and 429 behavior against docs.
- Audit webhook signature/timestamp/replay validation and subscription health.
- Verify use of current livestream v2 endpoints.
