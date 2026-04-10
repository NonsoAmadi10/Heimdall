import React from "react";
import { Doughnut, Line, Bar } from "react-chartjs-2";
import { ArcElement, BarElement, CategoryScale, Chart as ChartJS, Filler, Legend, LineElement, LinearScale, PointElement, Tooltip } from "chart.js";

ChartJS.register(ArcElement, CategoryScale, LinearScale, PointElement, LineElement, BarElement, Tooltip, Legend, Filler);

function number(value) {
  return Number.isFinite(value) ? value : 0;
}

function formatBytes(value) {
  if (!Number.isFinite(value)) return "N/A";
  const units = ["B", "KB", "MB", "GB", "TB"];
  let index = 0;
  let num = value;
  while (num >= 1024 && index < units.length - 1) {
    num /= 1024;
    index += 1;
  }
  return `${num.toFixed(index === 0 ? 0 : 2)} ${units[index]}`;
}

function Conn({ metrics, analytics, alerts }) {
  const latest = metrics.length > 0 ? metrics[metrics.length - 1] : null;
  const recentMetrics = metrics.slice(-12);
  const analyticsPoints = analytics?.points || [];

  const peerBreakdownData = {
    labels: ["Bitcoin Peers", "Lightning Peers"],
    datasets: [
      {
        data: [number(latest?.num_of_btc_peers), number(latest?.num_lnd_peers)],
        backgroundColor: ["rgba(59,130,246,0.75)", "rgba(16,185,129,0.75)"],
        borderColor: ["rgba(59,130,246,1)", "rgba(16,185,129,1)"],
        borderWidth: 1,
      },
    ],
  };

  const bandwidthTrendData = {
    labels: recentMetrics.map((item) => new Date(item.timestamp).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })),
    datasets: [
      {
        label: "Bandwidth In",
        data: recentMetrics.map((item) => number(item.btc_bandwidth_in)),
        borderColor: "rgba(59,130,246,1)",
        backgroundColor: "rgba(59,130,246,0.16)",
        fill: true,
        tension: 0.35,
      },
      {
        label: "Bandwidth Out",
        data: recentMetrics.map((item) => number(item.btc_bandwidth_out)),
        borderColor: "rgba(234,179,8,1)",
        backgroundColor: "rgba(234,179,8,0.16)",
        fill: true,
        tension: 0.35,
      },
    ],
  };

  const peerTrendData = {
    labels: analyticsPoints.map((item) => new Date(item.bucket_start).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })),
    datasets: [
      {
        label: "Avg BTC Peers",
        data: analyticsPoints.map((item) => number(item.avg_btc_peers)),
        backgroundColor: "rgba(59,130,246,0.85)",
      },
      {
        label: "Avg LND Peers",
        data: analyticsPoints.map((item) => number(item.avg_lnd_peers)),
        backgroundColor: "rgba(16,185,129,0.85)",
      },
    ],
  };

  return (
    <section className="space-y-4">
      <h2 className="text-lg font-semibold text-slate-100 sm:text-xl">Operations & Trends</h2>
      <div className="grid gap-4 xl:grid-cols-3">
        <div className="rounded-xl border border-slate-800 bg-slate-900/80 p-4">
          <p className="mb-3 text-sm text-slate-300">Peer Distribution</p>
          <Doughnut data={peerBreakdownData} />
        </div>
        <div className="rounded-xl border border-slate-800 bg-slate-900/80 p-4 xl:col-span-2">
          <p className="mb-3 text-sm text-slate-300">Recent Bandwidth (12 samples)</p>
          <Line data={bandwidthTrendData} options={{ responsive: true, maintainAspectRatio: false }} height={96} />
        </div>
      </div>

      <div className="grid gap-4 xl:grid-cols-3">
        <div className="rounded-xl border border-slate-800 bg-slate-900/80 p-4 xl:col-span-2">
          <p className="mb-3 text-sm text-slate-300">Analytics Peer Trend</p>
          <Bar data={peerTrendData} options={{ responsive: true, maintainAspectRatio: false }} height={96} />
        </div>
        <div className="rounded-xl border border-slate-800 bg-slate-900/80 p-4">
          <p className="text-sm text-slate-300">Current Snapshot</p>
          <dl className="mt-3 space-y-3 text-sm">
            <div className="flex items-center justify-between">
              <dt className="text-slate-400">BTC Inbound</dt>
              <dd className="font-medium text-slate-100">{formatBytes(latest?.btc_bandwidth_in)}</dd>
            </div>
            <div className="flex items-center justify-between">
              <dt className="text-slate-400">BTC Outbound</dt>
              <dd className="font-medium text-slate-100">{formatBytes(latest?.btc_bandwidth_out)}</dd>
            </div>
            <div className="flex items-center justify-between">
              <dt className="text-slate-400">Active Channels</dt>
              <dd className="font-medium text-slate-100">{latest?.num_active_channels ?? "N/A"}</dd>
            </div>
            <div className="flex items-center justify-between">
              <dt className="text-slate-400">Sync State</dt>
              <dd className={`font-medium ${latest?.synched_to_chain ? "text-emerald-400" : "text-red-400"}`}>
                {latest?.synched_to_chain ? "Synced" : "Unsynced"}
              </dd>
            </div>
            <div className="flex items-center justify-between">
              <dt className="text-slate-400">Open Alerts</dt>
              <dd className={`font-medium ${alerts.length > 0 ? "text-amber-400" : "text-emerald-400"}`}>{alerts.length}</dd>
            </div>
          </dl>
        </div>
      </div>

      <div className="rounded-xl border border-slate-800 bg-slate-900/80 p-4">
        <div className="mb-3 flex items-center justify-between">
          <p className="text-sm text-slate-300">Open Alerts</p>
          <span className="text-xs text-slate-500">{alerts.length} active</span>
        </div>
        {alerts.length === 0 ? (
          <p className="rounded-md bg-slate-800/70 px-3 py-2 text-sm text-emerald-300">No active alerts.</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full text-left text-sm">
              <thead className="text-slate-400">
                <tr>
                  <th className="px-3 py-2">Type</th>
                  <th className="px-3 py-2">Severity</th>
                  <th className="px-3 py-2">Value</th>
                  <th className="px-3 py-2">Threshold</th>
                  <th className="px-3 py-2">Triggered</th>
                </tr>
              </thead>
              <tbody>
                {alerts.map((alert) => (
                  <tr key={alert.id} className="border-t border-slate-800 text-slate-200">
                    <td className="px-3 py-2">{alert.type}</td>
                    <td className="px-3 py-2 capitalize">{alert.severity}</td>
                    <td className="px-3 py-2">{Number(alert.metric_value).toFixed(2)}</td>
                    <td className="px-3 py-2">{Number(alert.threshold).toFixed(2)}</td>
                    <td className="px-3 py-2">{new Date(alert.triggered_at).toLocaleString()}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </section>
  );
}

export default Conn;
