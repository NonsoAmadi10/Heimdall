import Conn from "@/components/Conn";
import NavBar from "@/components/NavBar";
import NodeInfo from "@/components/NodeInfo";
import { useCallback, useEffect, useMemo, useState } from "react";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:1700";
const POLL_INTERVAL_MS = 30000;

async function fetchJson(endpoint) {
  const response = await fetch(`${API_BASE_URL}${endpoint}`);
  if (!response.ok) {
    throw new Error(`Request failed (${response.status}) for ${endpoint}`);
  }
  return response.json();
}

function TopStat({ label, value }) {
  return (
    <div className="rounded-xl border border-slate-800 bg-slate-900/80 p-4">
      <p className="text-xs uppercase tracking-wide text-slate-400">{label}</p>
      <p className="mt-2 text-2xl font-semibold text-slate-100">{value}</p>
    </div>
  );
}

export default function Home() {
  const [nodeInfo, setNodeInfo] = useState(null);
  const [metrics, setMetrics] = useState([]);
  const [analytics, setAnalytics] = useState(null);
  const [alerts, setAlerts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [lastUpdated, setLastUpdated] = useState(null);

  const loadDashboard = useCallback(async () => {
    setLoading(true);
    setError("");

    try {
      const [nodeResponse, metricsResponse, analyticsResponse, alertsResponse] = await Promise.all([
        fetchJson("/node-info"),
        fetchJson("/conn-metrics"),
        fetchJson("/conn-metrics/analytics?interval_minutes=60"),
        fetchJson("/alerts?status=open"),
      ]);

      setNodeInfo(nodeResponse);
      setMetrics(metricsResponse?.data || []);
      setAnalytics(analyticsResponse?.data || null);
      setAlerts(alertsResponse?.data || []);
      setLastUpdated(new Date());
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load dashboard data.");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadDashboard();
    const timer = setInterval(loadDashboard, POLL_INTERVAL_MS);
    return () => clearInterval(timer);
  }, [loadDashboard]);

  const latestMetric = useMemo(() => (metrics.length > 0 ? metrics[metrics.length - 1] : null), [metrics]);
  const analyticsSummary = analytics?.summary;

  return (
    <main className="min-h-screen bg-slate-950">
      <NavBar
        loading={loading}
        onRefresh={loadDashboard}
        lastUpdated={lastUpdated ? lastUpdated.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit", second: "2-digit" }) : ""}
      />
      <div className="mx-auto w-full max-w-7xl space-y-6 px-4 py-6 sm:px-6 lg:px-8">
        {error ? (
          <div className="rounded-xl border border-red-900/80 bg-red-950/50 p-4 text-sm text-red-300">{error}</div>
        ) : null}

        <section className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
          <TopStat label="BTC Peers" value={latestMetric?.num_of_btc_peers ?? "N/A"} />
          <TopStat label="LND Peers" value={latestMetric?.num_lnd_peers ?? "N/A"} />
          <TopStat label="Sync Health (24h)" value={analyticsSummary ? `${analyticsSummary.sync_health_percent.toFixed(1)}%` : "N/A"} />
          <TopStat label="Open Alerts" value={alerts.length} />
        </section>

        <NodeInfo data={nodeInfo} />
        <Conn metrics={metrics} analytics={analytics} alerts={alerts} />
      </div>
    </main>
  );
}
