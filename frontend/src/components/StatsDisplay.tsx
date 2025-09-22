import { useEffect, useState } from "react";
import apiClient from "../api/apiClient";

interface UserStats {
  deduplicated_storage_usage_bytes: number;
  original_storage_usage_bytes: number;
  storage_savings_bytes: number;
  storage_savings_percentage: number;
}

const StatsDisplay = () => {
  const [stats, setStats] = useState<UserStats | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const response = await apiClient.get("/stats");
        setStats(response.data);
      } catch (err) {
        console.error("Failed to fetch stats:", err);
        setError("Could not load statistics.");
      }
    };
    fetchStats();
  }, []);

  const formatBytes = (bytes: number, decimals = 2) => {
    if (bytes === 0) return "0 Bytes";
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ["Bytes", "KB", "MB", "GB", "TB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
  };

  if (error) {
    return (
      <div className="p-4 bg-red-100 text-red-700 rounded-md">{error}</div>
    );
  }

  if (!stats) {
    return <div>Loading stats...</div>;
  }

  return (
    <div className="p-4 bg-white dark:bg-gray-800 rounded-lg shadow-md mb-8">
      <h2 className="text-xl font-bold mb-4 text-gray-800 dark:text-white">
        Storage Statistics
      </h2>
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-center">
        <div>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Used (Actual)
          </p>
          <p className="text-2xl font-semibold text-gray-800 dark:text-white">
            {formatBytes(stats.deduplicated_storage_usage_bytes)}
          </p>
        </div>
        <div>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Original Total
          </p>
          <p className="text-2xl font-semibold text-gray-800 dark:text-white">
            {formatBytes(stats.original_storage_usage_bytes)}
          </p>
        </div>
        <div>
          <p className="text-sm text-gray-500 dark:text-gray-400">Savings</p>
          <p className="text-2xl font-semibold text-green-600 dark:text-green-400">
            {formatBytes(stats.storage_savings_bytes)}
          </p>
        </div>
        <div>
          <p className="text-sm text-gray-500 dark:text-gray-400">Efficiency</p>
          <p className="text-2xl font-semibold text-green-600 dark:text-green-400">
            {stats.storage_savings_percentage.toFixed(2)}%
          </p>
        </div>
      </div>
    </div>
  );
};

export default StatsDisplay;
