import React from "react";

function StatCard({ label, value, monospace = false }) {
  return (
    <div className="rounded-xl border border-slate-800 bg-slate-900/80 p-4">
      <p className="text-xs uppercase tracking-wide text-slate-400">{label}</p>
      <p className={`mt-2 text-lg font-semibold text-slate-100 ${monospace ? "font-mono text-sm sm:text-base" : ""}`}>{value}</p>
    </div>
  );
}

function formatHashrate(value) {
  if (!Number.isFinite(value)) return "N/A";
  if (value >= 1e12) return `${(value / 1e12).toFixed(2)} TH/s`;
  if (value >= 1e9) return `${(value / 1e9).toFixed(2)} GH/s`;
  if (value >= 1e6) return `${(value / 1e6).toFixed(2)} MH/s`;
  return `${value.toFixed(0)} H/s`;
}

function NodeInfo({ data }) {
  const bitcoin = data?.bitcoin;
  const lightning = data?.lightning;

  return (
    <section className="space-y-4">
      <h2 className="text-lg font-semibold text-slate-100 sm:text-xl">Node Overview</h2>
      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <StatCard label="Network" value={bitcoin?.chain || "N/A"} />
        <StatCard label="Bitcoin Blocks" value={bitcoin?.no_of_blocks ?? "N/A"} />
        <StatCard label="Lightning Alias" value={lightning?.alias || "N/A"} />
        <StatCard label="Network Capacity" value={lightning?.network_capacity ?? "N/A"} />
        <StatCard label="Node Hashrate" value={formatHashrate(bitcoin?.hash_rate)} />
        <StatCard label="Block Propagation" value={Number.isFinite(bitcoin?.block_propagation) ? `${bitcoin.block_propagation.toFixed(2)} mins` : "N/A"} />
        <StatCard label="Bitcoin User Agent" value={bitcoin?.user_agent || "N/A"} monospace />
        <StatCard label="Lightning PubKey" value={lightning?.pub_key || "N/A"} monospace />
      </div>
    </section>
  );
}

export default NodeInfo;
