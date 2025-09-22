import { useEffect, useState } from "react";
import apiClient from "../api/apiClient";

interface SystemStats {
  total_users: number;
  total_files: number;
  total_storage_used: number;
  total_downloads: number;
}

interface AllFiles {
  id: number;
  filename: string;
  owner_email: string;
  size_bytes: number;
  upload_date: string;
  mime_type?: string;
}

const AdminDashboard = () => {
  const [stats, setStats] = useState<SystemStats | null>(null);
  const [files, setFiles] = useState<AllFiles[]>([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const [statsRes, filesRes] = await Promise.all([
          apiClient.get("/admin/stats"),
          apiClient.get("/admin/files"),
        ]);
        setStats(statsRes.data);
        setFiles(Array.isArray(filesRes.data) ? filesRes.data : []);
      } catch (err) {
        console.error("Failed to fetch admin data", err);
        setError(
          "Could not load admin dashboard. You may not have admin permissions."
        );
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return "0 Bytes";
    const k = 1024;
    const sizes = ["Bytes", "KB", "MB", "GB", "TB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-50 dark:bg-gray-900">
        <div className="text-gray-600 dark:text-gray-300">
          Loading admin dashboard...
        </div>
      </div>
    );
  }

  return (
    <div className="bg-gray-50 dark:bg-gray-900 min-h-screen">
      <div className="container mx-auto p-8">
        <h1 className="text-3xl font-bold mb-6 text-gray-800 dark:text-white">
          Admin Dashboard
        </h1>

        {error && (
          <p className="text-red-500 bg-red-100 p-3 rounded mb-4">{error}</p>
        )}

        {stats && (
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
            <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
              <div className="flex items-center">
                <div className="p-3 rounded-full bg-blue-500 bg-opacity-75">
                  <svg
                    className="h-8 w-8 text-white"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
                    />
                  </svg>
                </div>
                <div className="mx-5">
                  <h4 className="text-2xl font-semibold text-gray-700 dark:text-gray-200">
                    {stats.total_users}
                  </h4>
                  <div className="text-gray-500 dark:text-gray-400">
                    Total Users
                  </div>
                </div>
              </div>
            </div>

            <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
              <div className="flex items-center">
                <div className="p-3 rounded-full bg-green-500 bg-opacity-75">
                  <svg
                    className="h-8 w-8 text-white"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                    />
                  </svg>
                </div>
                <div className="mx-5">
                  <h4 className="text-2xl font-semibold text-gray-700 dark:text-gray-200">
                    {stats.total_files}
                  </h4>
                  <div className="text-gray-500 dark:text-gray-400">
                    Total Files
                  </div>
                </div>
              </div>
            </div>

            <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
              <div className="flex items-center">
                <div className="p-3 rounded-full bg-yellow-500 bg-opacity-75">
                  <svg
                    className="h-8 w-8 text-white"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4"
                    />
                  </svg>
                </div>
                <div className="mx-5">
                  <h4 className="text-2xl font-semibold text-gray-700 dark:text-gray-200">
                    {formatBytes(stats.total_storage_used)}
                  </h4>
                  <div className="text-gray-500 dark:text-gray-400">
                    Storage Used
                  </div>
                </div>
              </div>
            </div>

            <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
              <div className="flex items-center">
                <div className="p-3 rounded-full bg-purple-500 bg-opacity-75">
                  <svg
                    className="h-8 w-8 text-white"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                    />
                  </svg>
                </div>
                <div className="mx-5">
                  <h4 className="text-2xl font-semibold text-gray-700 dark:text-gray-200">
                    {stats.total_downloads}
                  </h4>
                  <div className="text-gray-500 dark:text-gray-400">
                    Total Downloads
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}

        <div className="bg-white dark:bg-gray-800 rounded shadow">
          <h2 className="text-xl font-bold p-4 border-b border-gray-200 dark:border-gray-700 text-gray-800 dark:text-white">
            All Files in System ({files.length})
          </h2>
          {files.length === 0 ? (
            <div className="p-8 text-center text-gray-500 dark:text-gray-400">
              No files found in the system.
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-left text-gray-800 dark:text-gray-200">
                <thead className="border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-700">
                  <tr>
                    <th className="p-3">Filename</th>
                    <th className="p-3">Owner</th>
                    <th className="p-3">Size</th>
                    <th className="p-3">Type</th>
                    <th className="p-3">Upload Date</th>
                  </tr>
                </thead>
                <tbody>
                  {files.map((file) => (
                    <tr
                      key={file.id}
                      className="border-b border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700"
                    >
                      <td className="p-3 font-medium">{file.filename}</td>
                      <td className="p-3">{file.owner_email}</td>
                      <td className="p-3">{formatBytes(file.size_bytes)}</td>
                      <td className="p-3 text-sm text-gray-600 dark:text-gray-400">
                        {file.mime_type || "Unknown"}
                      </td>
                      <td className="p-3">
                        {new Date(file.upload_date).toLocaleDateString()}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default AdminDashboard;
