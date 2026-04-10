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

Alerting configuration:

- `ALERT_MIN_BTC_PEERS` (default: `3`)
- `ALERT_LOOKBACK_SAMPLES` (default: `10`)
- `ALERT_BANDWIDTH_SPIKE_MULTIPLIER` (default: `2.5`)

## API Endpoints

- `GET /healthz` — service health check
- `GET /node-info` — aggregated Bitcoin + Lightning node information
- `GET /conn-metrics` — historical connection metrics from SQLite
- `GET /alerts` — list current and historical alerts (optional `status=open|acknowledged|resolved`)
- `PATCH /alerts/:id/ack` — acknowledge an alert
- `PATCH /alerts/:id/resolve` — resolve an alert manually

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

- [Alerting and Anomaly Detection](docs/alerting.md)

## License

MIT. See [LICENSE](LICENSE).
