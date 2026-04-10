import React from "react";

function NavBar({ lastUpdated, loading, onRefresh }) {
  return (
    <header className="sticky top-0 z-10 border-b border-slate-800/60 bg-slate-950/80 backdrop-blur">
      <div className="mx-auto flex w-full max-w-7xl items-center justify-between px-4 py-4 sm:px-6 lg:px-8">
        <div>
          <p className="text-xs uppercase tracking-[0.2em] text-blue-400">Heimdall</p>
          <h1 className="text-xl font-semibold text-slate-100 sm:text-2xl">Network Operations Dashboard</h1>
        </div>
        <div className="flex items-center gap-3">
          <span className="rounded-full border border-slate-700 px-3 py-1 text-xs text-slate-300">
            {lastUpdated ? `Updated ${lastUpdated}` : "Waiting for data"}
          </span>
          <button
            type="button"
            onClick={onRefresh}
            disabled={loading}
            className="rounded-md bg-blue-600 px-3 py-2 text-sm font-medium text-white transition hover:bg-blue-500 disabled:cursor-not-allowed disabled:opacity-60"
          >
            {loading ? "Refreshing..." : "Refresh"}
          </button>
        </div>
      </div>
    </header>
  );
}

export default NavBar;
