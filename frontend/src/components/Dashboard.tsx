import { useEffect, useState } from "react";
import apiClient from "../api/apiClient";
import ViewToggleButton from "./ViewToggleButton";
import FilePreviewModal from "./FilePreviewModal";
import UploadZone from "./UploadZone";
import SearchBar from "./SearchBar";
import TagManager from "./TagManager";
import FileActions from "./FileActions";
import Tooltip from "./Tooltip";

// ... (UserFile interface is the same)
interface UserFile {
  ID: number;
  Filename: string;
  MimeType: string;
  UploadDate: string;
  SizeBytes: number;
  Tags: string[];
}

const Dashboard = () => {
  const [files, setFiles] = useState<UserFile[]>([]);
  const [allFiles, setAllFiles] = useState<UserFile[]>([]); // Store all files for client-side filtering
  const [error, setError] = useState("");
  const [view, setView] = useState<"grid" | "list">("grid");
  const [selectedFile, setSelectedFile] = useState<UserFile | null>(null);
  const [showUpload, setShowUpload] = useState(false);

  const fetchAllFiles = async () => {
    try {
      // Fetch all files without search params
      const response = await apiClient.get("/search");
      const filesData = Array.isArray(response.data) ? response.data : [];
      console.log("Fetched files data:", filesData); // Debug log
      setAllFiles(filesData);
      setFiles(filesData); // Initially show all files
    } catch (error) {
      console.error("Failed to fetch files", error);
      setError("Could not load your files.");
    }
  };

  useEffect(() => {
    fetchAllFiles(); // Initial fetch
  }, []);

  const handleSearch = (searchTerm: string) => {
    if (!searchTerm.trim()) {
      setFiles(allFiles); // Show all files if search is empty
      return;
    }

    const searchLower = searchTerm.toLowerCase();
    const filteredFiles = allFiles.filter((file) => {
      // Search in filename
      const filenameMatch = file.Filename?.toLowerCase().includes(searchLower);

      // Search in tags
      const tagsMatch = file.Tags?.some((tag) =>
        tag.toLowerCase().includes(searchLower)
      );

      // Search in date (formatted)
      const dateMatch =
        file.UploadDate &&
        new Date(file.UploadDate)
          .toLocaleDateString()
          .toLowerCase()
          .includes(searchLower);

      // Search in file size (converted to string)
      const sizeMatch =
        file.SizeBytes && file.SizeBytes.toString().includes(searchLower);

      // Search in mime type
      const mimeMatch = file.MimeType?.toLowerCase().includes(searchLower);

      return filenameMatch || tagsMatch || dateMatch || sizeMatch || mimeMatch;
    });

    setFiles(filteredFiles);
  };

  // ... (handleDelete and formatBytes are the same)
  const handleDelete = async (fileId: number) => {
    if (window.confirm("Are you sure you want to delete this file?")) {
      try {
        await apiClient.delete(`/files/${fileId}`);
        fetchAllFiles();
      } catch (error) {
        console.error("Failed to delete file", error);
        setError("Could not delete the file.");
      }
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

  // Handle adding a tag via API
  const handleAddTag = async (fileId: number, tag: string) => {
    try {
      await apiClient.post(`/files/${fileId}/tags`, { tag });
      fetchAllFiles(); // Refresh file list to show new tag
    } catch (error) {
      console.error("Failed to add tag", error);
      alert("Failed to add tag.");
    }
  };

  // Handle removing a tag via API
  const handleRemoveTag = async (fileId: number, tag: string) => {
    try {
      // For DELETE requests with a body, Axios requires the body to be in the 'data' property
      await apiClient.delete(`/files/${fileId}/tags`, { data: { tag } });
      fetchAllFiles(); // Refresh file list to show tag has been removed
    } catch (error) {
      console.error("Failed to remove tag", error);
      alert("Failed to remove tag.");
    }
  };

  return (
    <div className="container mx-auto p-8 bg-gray-50 dark:bg-gray-900 min-h-screen">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-gray-800 dark:text-white">
          Your Files
        </h1>
        <div>
          <button
            onClick={() => setShowUpload(!showUpload)}
            className="bg-blue-500 text-white px-4 py-2 rounded-md mr-4 hover:bg-blue-600"
          >
            {showUpload ? "Cancel Upload" : "Upload File"}
          </button>
          <ViewToggleButton view={view} setView={setView} />
        </div>
      </div>

      <SearchBar onSearch={handleSearch} />

      {showUpload && (
        <UploadZone
          onUploadSuccess={() => {
            setShowUpload(false);
            fetchAllFiles();
          }}
        />
      )}

      {error && (
        <p className="text-red-500 bg-red-100 p-3 rounded mb-4">{error}</p>
      )}

      {/* List View */}
      <div
        className={`bg-white dark:bg-gray-800 rounded shadow ${
          view === "list" ? "block" : "hidden"
        }`}
      >
        <table className="w-full text-left text-gray-800 dark:text-gray-200">
          <thead className="border-b border-gray-200 dark:border-gray-700">
            <tr>
              <th className="p-3">Name</th>
              <th className="p-3">Tags</th>
              <th className="p-3">Size</th>
              <th className="p-3">Date Added</th>
              <th className="p-3">Actions</th>
            </tr>
          </thead>
          <tbody>
            {files.map((file) => (
              <tr
                key={file.ID}
                className="border-b border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700"
              >
                <td className="p-3 cursor-pointer max-w-xs">
                  <Tooltip content={file.Filename}>
                    <div
                      className="truncate hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
                      onClick={() => setSelectedFile(file)}
                    >
                      {file.Filename}
                    </div>
                  </Tooltip>
                </td>
                <td className="p-3 max-w-xs">
                  <div className="max-w-full">
                    <TagManager
                      tags={file.Tags || []}
                      onAddTag={(tag) => handleAddTag(file.ID, tag)}
                      onRemoveTag={(tag) => handleRemoveTag(file.ID, tag)}
                    />
                  </div>
                </td>
                <td className="p-3">
                  <Tooltip
                    content={`${formatBytes(file.SizeBytes || 0)} (${
                      file.SizeBytes
                    } bytes)`}
                  >
                    <span>{formatBytes(file.SizeBytes || 0)}</span>
                  </Tooltip>
                </td>
                <td className="p-3">
                  <Tooltip content={new Date(file.UploadDate).toLocaleString()}>
                    <span>
                      {new Date(file.UploadDate).toLocaleDateString()}
                    </span>
                  </Tooltip>
                </td>
                <td className="p-3">
                  <FileActions file={file} onDelete={handleDelete} />
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Grid View */}
      <div
        className={`grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4 ${
          view === "grid" ? "grid" : "hidden"
        }`}
      >
        {files.map((file) => (
          <div
            key={file.ID}
            className="border dark:border-gray-700 rounded-lg p-2 text-center shadow hover:shadow-lg bg-white dark:bg-gray-800 text-gray-800 dark:text-gray-200 min-w-0 flex flex-col"
          >
            <div
              onClick={() => setSelectedFile(file)}
              className="flex items-center justify-center h-24 bg-gray-100 dark:bg-gray-700 rounded cursor-pointer"
            >
              <span className="text-4xl">📄</span>
            </div>
            <Tooltip content={file.Filename}>
              <p className="font-semibold truncate mt-2 cursor-default min-w-0">
                {file.Filename}
              </p>
            </Tooltip>
            <div className="mt-2">
              <TagManager
                tags={file.Tags || []}
                onAddTag={(tag) => handleAddTag(file.ID, tag)}
                onRemoveTag={(tag) => handleRemoveTag(file.ID, tag)}
              />
            </div>
            <div className="mt-2">
              <FileActions file={file} onDelete={handleDelete} />
            </div>
          </div>
        ))}
      </div>

      <FilePreviewModal
        file={selectedFile}
        onClose={() => setSelectedFile(null)}
      />
    </div>
  );
};

export default Dashboard;
