# Historical Analytics API

The historical analytics endpoint provides bucketed summaries of connection metrics across a time window.

## Endpoint

`GET /conn-metrics/analytics`

## Query Parameters

- `from` (optional): RFC3339 timestamp, default is current time minus 24 hours.
- `to` (optional): RFC3339 timestamp, default is current time.
- `interval_minutes` (optional): bucket size in minutes, accepted range is 1-1440, default is 60.

## Example Request

```bash
curl "http://localhost:1700/conn-metrics/analytics?from=2026-04-09T00:00:00Z&to=2026-04-10T00:00:00Z&interval_minutes=60"
```

## Response Shape

```json
{
  "success": true,
  "data": {
    "from": "2026-04-09T00:00:00Z",
    "to": "2026-04-10T00:00:00Z",
    "interval_minutes": 60,
    "points": [
      {
        "bucket_start": "2026-04-09T00:00:00Z",
        "bucket_end": "2026-04-09T01:00:00Z",
        "samples": 3,
        "avg_btc_peers": 8.67,
        "avg_lnd_peers": 6.33,
        "avg_bandwidth_in": 12000,
        "avg_bandwidth_out": 11400,
        "max_bandwidth_in": 14000,
        "max_bandwidth_out": 13000,
        "sync_health_percent": 100
      }
    ],
    "summary": {
      "total_samples": 3,
      "avg_btc_peers": 8.67,
      "avg_lnd_peers": 6.33,
      "avg_bandwidth_in": 12000,
      "avg_bandwidth_out": 11400,
      "sync_health_percent": 100
    }
  }
}
```

## Notes

- Empty windows return a successful response with an empty `points` array.
- `sync_health_percent` is calculated as: `(synced samples / total samples) * 100`.
