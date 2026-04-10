# Heimdall — Bitcoin and Lightning Network Metrics

Heimdall is a Go service that collects and exposes node connectivity and network metrics from `btcd` and `lnd`, then stores historical data in SQLite via GORM.

## Requirements

- Go 1.19+
- Running `btcd` RPC endpoint
- Running `lnd` endpoint with TLS cert and admin macaroon

## Setup

1. Clone the repo and install dependencies:
```bash
git clone https://github.com/NonsoAmadi10/Heimdall
cd Heimdall
go mod download
```
2. Create environment variables:
```bash
cp .env.example .env
```
3. Build and run:
```bash
go build ./...
go run .
```

## Configuration

Environment variables used by the Bitcoin RPC client:

- `BTC_HOST` (example: `127.0.0.1:18334`)
- `BTC_USER`
- `BTC_PASS`

Lightning RPC certificate and macaroon paths are currently resolved from:

`$HOME/app_container/lightning/...`

## API Endpoints

- `GET /healthz` — service health check
- `GET /node-info` — aggregated Bitcoin + Lightning node information
- `GET /conn-metrics` — historical connection metrics from SQLite
- `GET /conn-metrics/analytics` — aggregated historical analytics with time buckets

### Analytics Query Parameters

- `from` (optional, RFC3339 timestamp, default: now - 24h)
- `to` (optional, RFC3339 timestamp, default: now)
- `interval_minutes` (optional, integer 1-1440, default: 60)

Example:

```bash
curl "http://localhost:1700/conn-metrics/analytics?from=2026-04-09T00:00:00Z&to=2026-04-10T00:00:00Z&interval_minutes=30"
```

## Frontend Dashboard

```bash
cd dashboard
yarn install
yarn dev
```

Open `http://localhost:3000`.

## Development

Run tests and build locally:

```bash
go test ./...
go build ./...
```

Continuous integration runs these checks on every push and pull request.

## Documentation

- [Historical Analytics API](docs/historical-analytics.md)

## License

MIT. See [LICENSE](LICENSE).
