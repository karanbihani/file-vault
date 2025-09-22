import { useEffect, useState } from "react";
import apiClient from "../api/apiClient";

interface SharedFile {
  ID: number;
  Filename: string;
  MimeType: string;
  UploadDate: string;
  SizeBytes: number;
  Tags: string[];
}

const SharedWithMe = () => {
  const [files, setFiles] = useState<SharedFile[]>([]);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchSharedFiles = async () => {
      try {
        const response = await apiClient.get("/files/shared-with-me");
        setFiles(response.data || []);
      } catch (err) {
        console.error("Failed to fetch shared files", err);
        setError("Could not load files shared with you.");
      }
    };
    fetchSharedFiles();
  }, []);

  const handleDownloadShared = async (fileId: number, filename: string) => {
    try {
      const response = await apiClient.get(`/files/${fileId}/download`, {
        responseType: "blob", // Important: Tell Axios to expect binary data
      });
      // Create a temporary URL from the blob data
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement("a");
      link.href = url;
      link.setAttribute("download", filename); // Set the filename
      document.body.appendChild(link);
      link.click();
      link.remove(); // Clean up the temporary link
      window.URL.revokeObjectURL(url); // Clean up the temporary URL
    } catch (error) {
      console.error("Download failed", error);
      alert("Could not download the file.");
    }
  };

  const formatBytes = (bytes: number, decimals = 2) => {
    if (bytes === 0) return "0 Bytes";
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ["Bytes", "KB", "MB", "GB", "TB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
  };

  return (
    <div className="bg-gray-50 dark:bg-gray-900 min-h-screen">
      <div className="container mx-auto p-8">
        <h1 className="text-3xl font-bold mb-6 text-gray-800 dark:text-white">
          Files Shared With Me
        </h1>
        {error && (
          <p className="text-red-500 bg-red-100 p-3 rounded mb-4">{error}</p>
        )}

        {files.length === 0 && !error ? (
          <div className="bg-white dark:bg-gray-800 rounded shadow p-8 text-center">
            <p className="text-gray-600 dark:text-gray-400">
              No files have been shared with you yet.
            </p>
          </div>
        ) : (
          <div className="bg-white dark:bg-gray-800 rounded shadow">
            <table className="w-full text-left text-gray-800 dark:text-gray-200">
              <thead className="border-b border-gray-200 dark:border-gray-700">
                <tr>
                  <th className="p-3">Filename</th>
                  <th className="p-3">Size</th>
                  <th className="p-3">Date Shared</th>
                  <th className="p-3">Actions</th>
                </tr>
              </thead>
              <tbody>
                {files.map((file) => (
                  <tr
                    key={file.ID}
                    className="border-b border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700"
                  >
                    <td className="p-3">{file.Filename}</td>
                    <td className="p-3">{formatBytes(file.SizeBytes || 0)}</td>
                    <td className="p-3">
                      {new Date(file.UploadDate).toLocaleDateString()}
                    </td>
                    <td className="p-3">
                      <button
                        onClick={() =>
                          handleDownloadShared(file.ID, file.Filename)
                        }
                        className="text-blue-500 hover:underline"
                      >
                        Download
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
};

export default SharedWithMe;
