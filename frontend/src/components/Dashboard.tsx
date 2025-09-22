import { useEffect, useState, useCallback } from "react";
import apiClient from "../api/apiClient";
import ViewToggleButton from "./ViewToggleButton";
import FilePreviewModal from "./FilePreviewModal";
import UploadZone from "./UploadZone";
import SearchBar from "./SearchBar";
import TagManager from "./TagManager";
import FileActions from "./FileActions";
import Tooltip from "./Tooltip";
import AdvancedFilters, { FilterOptions } from "./AdvancedFilters";
import {
  FiFile,
  FiUpload,
  FiFilter,
  FiSearch,
  FiAlertCircle,
  FiLoader,
} from "react-icons/fi";

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
  const [error, setError] = useState("");
  const [view, setView] = useState<"grid" | "list">("grid");
  const [selectedFile, setSelectedFile] = useState<UserFile | null>(null);
  const [showUpload, setShowUpload] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [currentFilters, setCurrentFilters] = useState<FilterOptions>({});
  const [searchTerm, setSearchTerm] = useState("");

  // Debounced search function
  const fetchFilteredFiles = useCallback(
    async (filters: FilterOptions, search?: string) => {
      setIsLoading(true);
      setError("");

      try {
        // Combine search term with filename filter
        const finalFilters = {
          ...filters,
          filename: search || filters.filename || undefined,
        };

        // Build query parameters
        const params = new URLSearchParams();

        Object.entries(finalFilters).forEach(([key, value]) => {
          if (value !== undefined && value !== "" && value !== null) {
            params.append(key, String(value));
          }
        });

        let url = "/search";
        if (params.toString()) {
          url += `?${params.toString()}`;
        }

        const response = await apiClient.get(url);
        console.log("Search response:", response.data); // Debug log
        const filesData = Array.isArray(response.data) ? response.data : [];
        setFiles(filesData);
      } catch (error: any) {
        console.error("Error fetching files:", error);
        console.log("Error response:", error.response); // Debug log

        // Always set files to empty array when there's an error or no data
        setFiles([]);

        if (error.response?.status === 404) {
          // 404 likely means no files found, which is normal
          setError("");
        } else {
          setError(error.response?.data?.error || "Failed to fetch files");
        }
      } finally {
        setIsLoading(false);
      }
    },
    []
  );

  // Initial load - fetch all files
  useEffect(() => {
    fetchFilteredFiles({});
  }, [fetchFilteredFiles]);

  // Handle filter changes
  const handleFiltersChange = useCallback(
    (filters: FilterOptions) => {
      setCurrentFilters(filters);
      fetchFilteredFiles(filters, searchTerm);
    },
    [fetchFilteredFiles, searchTerm]
  );

  // Handle search changes
  const handleSearchChange = useCallback(
    (search: string) => {
      setSearchTerm(search);
      fetchFilteredFiles(currentFilters, search);
    },
    [fetchFilteredFiles, currentFilters]
  );

  // Reset all filters and search
  const handleReset = useCallback(() => {
    setCurrentFilters({});
    setSearchTerm("");
    fetchFilteredFiles({});
  }, [fetchFilteredFiles]);

  const handleDelete = async (fileId: number) => {
    if (window.confirm("Are you sure you want to delete this file?")) {
      try {
        await apiClient.delete(`/files/${fileId}`);
        // Refresh the current filtered view
        fetchFilteredFiles(currentFilters, searchTerm);
      } catch (error) {
        console.error("Failed to delete file", error);
        setError("Could not delete the file.");
      }
    }
  };

  const formatBytes = (bytes: number): string => {
    if (bytes === 0) return "0 Bytes";
    const k = 1024;
    const sizes = ["Bytes", "KB", "MB", "GB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
  };

  const handleAddTag = async (fileId: number, tag: string) => {
    try {
      await apiClient.post(`/files/${fileId}/tags`, { tag });
      // Refresh the current filtered view
      fetchFilteredFiles(currentFilters, searchTerm);
    } catch (error) {
      console.error("Failed to add tag", error);
    }
  };

  const handleRemoveTag = async (fileId: number, tag: string) => {
    try {
      await apiClient.delete(`/files/${fileId}/tags`, { data: { tag } });
      // Refresh the current filtered view
      fetchFilteredFiles(currentFilters, searchTerm);
    } catch (error) {
      console.error("Failed to remove tag", error);
    }
  };

  const refreshFiles = () => {
    fetchFilteredFiles(currentFilters, searchTerm);
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 transition-colors duration-200">
      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-between items-center mb-6">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">
            File Vault
          </h1>
          <button
            onClick={() => setShowUpload(true)}
            className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition-colors flex items-center"
          >
            <FiUpload className="mr-2" />
            Upload Files
          </button>
        </div>

        {/* Advanced Filters */}
        <AdvancedFilters
          onFiltersChange={handleFiltersChange}
          onReset={handleReset}
        />

        {/* Search and View Controls */}
        <div className="flex justify-between items-center mb-6">
          <SearchBar onSearch={handleSearchChange} />
          <ViewToggleButton view={view} setView={setView} />
        </div>

        {/* Loading and Error States */}
        {isLoading && (
          <div className="flex justify-center items-center py-8">
            <FiLoader className="animate-spin mr-2 text-gray-600 dark:text-gray-400" />
            <div className="text-gray-600 dark:text-gray-400">
              Loading files...
            </div>
          </div>
        )}

        {error && (
          <div className="bg-red-100 dark:bg-red-900/20 border border-red-400 dark:border-red-700 text-red-700 dark:text-red-400 px-4 py-3 rounded mb-4 flex items-center">
            <FiAlertCircle className="mr-2 flex-shrink-0" />
            {error}
          </div>
        )}

        {/* Results Count */}
        {!isLoading && (
          <div className="mb-4 text-sm text-gray-600 dark:text-gray-400">
            Showing {files.length} file{files.length !== 1 ? "s" : ""}
          </div>
        )}

        {/* File List View */}
        <div className={`space-y-2 ${view === "list" ? "block" : "hidden"}`}>
          {files.map((file) => (
            <div
              key={file.ID}
              className="flex items-center justify-between p-4 border dark:border-gray-700 rounded-lg hover:shadow-md bg-white dark:bg-gray-800 text-gray-800 dark:text-gray-200"
            >
              <div className="flex-1 min-w-0">
                <Tooltip content={file.Filename}>
                  <p
                    onClick={() => setSelectedFile(file)}
                    className="truncate hover:text-blue-600 dark:hover:text-blue-400 transition-colors cursor-pointer"
                  >
                    {file.Filename}
                  </p>
                </Tooltip>
                <div className="flex items-center gap-4 mt-1 text-sm text-gray-500 dark:text-gray-400">
                  <span>{formatBytes(file.SizeBytes)}</span>
                  <span>{file.MimeType}</span>
                  <span>{new Date(file.UploadDate).toLocaleDateString()}</span>
                </div>
                <div className="mt-1">
                  <TagManager
                    tags={file.Tags || []}
                    onAddTag={(tag) => handleAddTag(file.ID, tag)}
                    onRemoveTag={(tag) => handleRemoveTag(file.ID, tag)}
                  />
                </div>
              </div>
              <div className="ml-4">
                <FileActions file={file} onDelete={handleDelete} />
              </div>
            </div>
          ))}
        </div>

        {/* File Grid View */}
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
                <FiFile className="text-4xl text-gray-600 dark:text-gray-400" />
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

        {/* Empty State */}
        {!isLoading && files.length === 0 && (
          <div className="text-center py-12">
            <div className="flex flex-col items-center">
              {Object.values(currentFilters).some(
                (v) => v !== undefined && v !== ""
              ) || searchTerm ? (
                <>
                  <FiSearch className="text-6xl text-gray-400 dark:text-gray-500 mb-4" />
                  <div className="text-gray-500 dark:text-gray-400 text-lg mb-2">
                    No files match your current filters
                  </div>
                  <div className="text-gray-400 dark:text-gray-500 text-sm mb-4">
                    Try adjusting your search criteria or clearing filters
                  </div>
                  <button
                    onClick={handleReset}
                    className="px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-md transition-colors flex items-center"
                  >
                    <FiFilter className="mr-2" />
                    Clear all filters
                  </button>
                </>
              ) : (
                <>
                  <FiUpload className="text-6xl text-gray-400 dark:text-gray-500 mb-4" />
                  <div className="text-gray-500 dark:text-gray-400 text-lg mb-2">
                    No files uploaded yet
                  </div>
                  <div className="text-gray-400 dark:text-gray-500 text-sm mb-4">
                    Upload your first file to get started
                  </div>
                  <button
                    onClick={() => setShowUpload(true)}
                    className="px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-md transition-colors flex items-center"
                  >
                    <FiUpload className="mr-2" />
                    Upload Files
                  </button>
                </>
              )}
            </div>
          </div>
        )}

        {/* File Preview Modal */}
        {selectedFile && (
          <FilePreviewModal
            file={selectedFile}
            onClose={() => setSelectedFile(null)}
          />
        )}

        {/* Upload Modal */}
        {showUpload && (
          <UploadZone
            onUploadSuccess={refreshFiles}
            onClose={() => setShowUpload(false)}
          />
        )}
      </div>
    </div>
  );
};

export default Dashboard;
