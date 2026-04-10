# Heimdall Dashboard

Professional monitoring UI for the Heimdall backend APIs.

## Features

- Node overview for Bitcoin and Lightning
- Real-time operational KPIs (peers, sync health, alert count)
- Bandwidth and peer trend visualizations
- Open alert table for active incidents
- Auto-refresh every 30 seconds plus manual refresh

## Configuration

Set the backend URL in `dashboard/.env.local`:

```bash
NEXT_PUBLIC_API_BASE_URL=http://localhost:1700
```

If omitted, it defaults to `http://localhost:1700`.

## Getting started

Install dependencies and run the app:

```bash
yarn dev
```

Open `http://localhost:3000`.

## Build and lint

```bash
yarn lint
yarn build
```

## API dependencies

- `GET /node-info`
- `GET /conn-metrics`
- `GET /conn-metrics/analytics?interval_minutes=60`
- `GET /alerts?status=open`

## Learn More

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

You can check out [the Next.js GitHub repository](https://github.com/vercel/next.js/) - your feedback and contributions are welcome!

## Deploy on Vercel

The easiest way to deploy your Next.js app is to use the [Vercel Platform](https://vercel.com/new?utm_medium=default-template&filter=next.js&utm_source=create-next-app&utm_campaign=create-next-app-readme) from the creators of Next.js.

Check out our [Next.js deployment documentation](https://nextjs.org/docs/deployment) for more details.
