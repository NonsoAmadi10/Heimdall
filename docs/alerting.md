# Alerting and Anomaly Detection

Heimdall now evaluates each collected metric sample against alerting rules and stores alert history in SQLite.

## Rules

1. `low_btc_peers` (`warning`)
   - Triggered when Bitcoin peer count drops below `ALERT_MIN_BTC_PEERS`.
2. `sync_stalled` (`critical`)
   - Triggered when `SyncedToChain` is `false`.
3. `bandwidth_spike` (`warning`)
   - Triggered when current total BTC bandwidth (`in + out`) exceeds:
     `average(last N samples) * ALERT_BANDWIDTH_SPIKE_MULTIPLIER`
   - `N` is controlled by `ALERT_LOOKBACK_SAMPLES`.

## Alert Lifecycle

- `open`: newly triggered and active
- `acknowledged`: acknowledged by an operator
- `resolved`: auto-resolved when condition clears, or manually resolved via API

## Endpoints

### List alerts

`GET /alerts`

Optional query:

- `status=open|acknowledged|resolved`

### Acknowledge alert

`PATCH /alerts/:id/ack`

### Resolve alert

`PATCH /alerts/:id/resolve`

## Example

```bash
curl "http://localhost:1700/alerts?status=open"
curl -X PATCH "http://localhost:1700/alerts/12/ack"
curl -X PATCH "http://localhost:1700/alerts/12/resolve"
```
