import { useState, useCallback } from "react";
import apiClient from "../api/apiClient";
import { AxiosError } from "axios";
import {
  FiUpload,
  FiX,
  FiCheck,
  FiAlertTriangle,
  FiTrash2,
} from "react-icons/fi";

interface FileUploadStatus {
  name: string;
  status: "pending" | "uploading" | "success" | "error";
  progress: number;
  error?: string;
}

const UploadZone = ({
  onUploadSuccess,
  onClose,
}: {
  onUploadSuccess: () => void;
  onClose?: () => void;
}) => {
  const [isDragging, setIsDragging] = useState(false);
  const [statuses, setStatuses] = useState<FileUploadStatus[]>([]);

  const uploadFile = async (file: File, index: number) => {
    const formData = new FormData();
    formData.append("files", file);

    try {
      setStatuses((prev) =>
        prev.map((s, i) => (i === index ? { ...s, status: "uploading" } : s))
      );

      await apiClient.upload("/files", formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
        onUploadProgress: (progressEvent) => {
          const percent = Math.round(
            (progressEvent.loaded * 100) / (progressEvent.total || 1)
          );
          setStatuses((prev) =>
            prev.map((s, i) => (i === index ? { ...s, progress: percent } : s))
          );
        },
      });

      setStatuses((prev) =>
        prev.map((s, i) =>
          i === index ? { ...s, status: "success", progress: 100 } : s
        )
      );
    } catch (err) {
      const axiosError = err as AxiosError<{ error: string }>;
      let errorMessage = "Upload failed - unknown error";

      if (axiosError.response) {
        errorMessage =
          axiosError.response.data?.error ||
          `Server error: ${axiosError.response.status}`;
      } else if (axiosError.request) {
        errorMessage = "Network error - please check your connection";
      } else {
        errorMessage = axiosError.message || "Upload failed";
      }

      setStatuses((prev) =>
        prev.map((s, i) =>
          i === index ? { ...s, status: "error", error: errorMessage } : s
        )
      );
    }
  };

  const handleFiles = async (files: FileList) => {
    if (files.length === 0) return;

    const initialStatuses = Array.from(files).map((file) => ({
      name: file.name,
      status: "pending" as "pending",
      progress: 0,
    }));
    setStatuses(initialStatuses);

    // Upload files sequentially - rate limiter handles delays automatically
    for (let i = 0; i < files.length; i++) {
      await uploadFile(files[i], i);
    }

    onUploadSuccess();
  };

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  }, []);

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  }, []);

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    handleFiles(e.dataTransfer.files);
  }, []);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      handleFiles(e.target.files);
    }
  };

  const clearStatuses = () => {
    setStatuses([]);
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white dark:bg-gray-800 rounded-lg p-6 max-w-2xl w-full mx-4 max-h-[80vh] overflow-y-auto">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100">
            Upload Files
          </h2>
          {onClose && (
            <button
              onClick={onClose}
              className="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 p-1 rounded-md hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
            >
              <FiX className="text-xl" />
            </button>
          )}
        </div>

        <div className="mb-8 p-4 border rounded-lg dark:border-gray-700">
          <div
            className={`p-8 border-2 border-dashed rounded-lg text-center transition-colors ${
              isDragging
                ? "border-blue-500 bg-blue-50 dark:bg-blue-900/20"
                : "border-gray-300 dark:border-gray-600 hover:border-gray-400 dark:hover:border-gray-500"
            }`}
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            onDrop={handleDrop}
          >
            <input
              type="file"
              multiple
              onChange={handleFileChange}
              className="hidden"
              id="file-input"
            />
            <label
              htmlFor="file-input"
              className="cursor-pointer flex flex-col items-center space-y-2"
            >
              <FiUpload className="text-4xl text-gray-600 dark:text-gray-400" />
              <div className="text-lg font-medium text-gray-700 dark:text-gray-300">
                Click to upload or drag and drop files
              </div>
              <div className="text-sm text-gray-500 dark:text-gray-400">
                Supports multiple files
              </div>
            </label>
          </div>

          {statuses.length > 0 && (
            <div className="mt-4 space-y-3">
              <div className="flex justify-between items-center">
                <h3 className="text-lg font-medium text-gray-700 dark:text-gray-300">
                  Upload Progress
                </h3>
                <button
                  onClick={clearStatuses}
                  className="text-sm text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 flex items-center px-2 py-1 rounded hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                >
                  <FiTrash2 className="mr-1" />
                  Clear
                </button>
              </div>
              {statuses.map((status, index) => (
                <div
                  key={index}
                  className="border rounded-lg p-3 dark:border-gray-600"
                >
                  <div className="flex justify-between items-start mb-2">
                    <p className="text-sm font-medium text-gray-700 dark:text-gray-300 truncate">
                      {status.name}
                    </p>
                    <span
                      className={`text-xs px-2 py-1 rounded ${
                        status.status === "success"
                          ? "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200"
                          : status.status === "error"
                          ? "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200"
                          : status.status === "uploading"
                          ? "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200"
                          : "bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200"
                      }`}
                    >
                      {status.status}
                    </span>
                  </div>

                  {status.status === "uploading" && (
                    <div className="w-full bg-gray-200 rounded-full h-2.5 dark:bg-gray-700">
                      <div
                        className="bg-blue-600 h-2.5 rounded-full transition-all duration-300"
                        style={{ width: `${status.progress}%` }}
                      />
                    </div>
                  )}

                  {status.status === "success" && (
                    <p className="text-green-600 dark:text-green-400 text-sm flex items-center">
                      <FiCheck className="mr-2" />
                      Upload completed successfully!
                    </p>
                  )}

                  {status.status === "error" && (
                    <p className="text-red-600 dark:text-red-400 text-sm flex items-center">
                      <FiAlertTriangle className="mr-2" />
                      {status.error}
                    </p>
                  )}
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default UploadZone;
