import { useState, useCallback } from "react";
import apiClient from "../api/apiClient";

interface UploadZoneProps {
  onUploadSuccess: () => void; // A function to call to refresh the file list
}

interface FileUploadStatus {
  file: File;
  status: "pending" | "uploading" | "success" | "error";
  progress: number;
  error?: string;
}

const UploadZone = ({ onUploadSuccess }: UploadZoneProps) => {
  const [isDragging, setIsDragging] = useState(false);
  const [uploadStatuses, setUploadStatuses] = useState<FileUploadStatus[]>([]);

  const handleUpload = async (files: FileList) => {
    if (files.length === 0) return;

    // Initialize status for all files
    const initialStatuses = Array.from(files).map((file) => ({
      file,
      status: "pending" as const,
      progress: 0,
    }));
    setUploadStatuses(initialStatuses);

    // Upload files sequentially for better progress tracking
    for (let i = 0; i < initialStatuses.length; i++) {
      const formData = new FormData();
      formData.append("files", initialStatuses[i].file);

      try {
        // Update status to uploading
        setUploadStatuses((prev) =>
          prev.map((s, idx) => (idx === i ? { ...s, status: "uploading" } : s))
        );

        await apiClient.post("/files", formData, {
          headers: {
            "Content-Type": "multipart/form-data",
          },
          onUploadProgress: (progressEvent) => {
            const percentCompleted = Math.round(
              (progressEvent.loaded * 100) / (progressEvent.total || 1)
            );
            setUploadStatuses((prev) =>
              prev.map((s, idx) =>
                idx === i ? { ...s, progress: percentCompleted } : s
              )
            );
          },
        });

        // Update status to success
        setUploadStatuses((prev) =>
          prev.map((s, idx) =>
            idx === i ? { ...s, status: "success", progress: 100 } : s
          )
        );
      } catch (error: any) {
        const errorMessage = error.response?.data?.error || "Upload failed";
        setUploadStatuses((prev) =>
          prev.map((s, idx) =>
            idx === i ? { ...s, status: "error", error: errorMessage } : s
          )
        );
      }
    }

    onUploadSuccess(); // Refresh the file list in the parent component
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
    handleUpload(e.dataTransfer.files);
  }, []);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      handleUpload(e.target.files);
    }
  };

  return (
    <div className="mb-8">
      <label
        htmlFor="file-upload"
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
        className={`flex justify-center w-full h-32 px-4 transition bg-white border-2 ${
          isDragging ? "border-blue-500" : "border-gray-300"
        } border-dashed rounded-md appearance-none cursor-pointer hover:border-gray-400 focus:outline-none`}
      >
        <span className="flex items-center space-x-2">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="w-6 h-6 text-gray-600"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            strokeWidth="2"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
            />
          </svg>
          <span className="font-medium text-gray-600">
            Drop files to attach, or{" "}
            <span className="text-blue-600 underline">browse</span>
          </span>
        </span>
        <input
          id="file-upload"
          type="file"
          multiple
          className="hidden"
          onChange={handleFileChange}
        />
      </label>

      {/* Upload Status Display */}
      {uploadStatuses.length > 0 && (
        <div className="mt-4 space-y-3">
          <h4 className="font-medium text-gray-700 dark:text-gray-300">
            Upload Progress:
          </h4>
          {uploadStatuses.map((status, index) => (
            <div
              key={index}
              className="bg-gray-50 dark:bg-gray-700 p-3 rounded-lg"
            >
              <div className="flex justify-between items-center mb-2">
                <p className="text-sm font-medium text-gray-800 dark:text-gray-200 truncate">
                  {status.file.name}
                </p>
                <span
                  className={`text-xs px-2 py-1 rounded ${
                    status.status === "success"
                      ? "bg-green-100 text-green-800"
                      : status.status === "error"
                      ? "bg-red-100 text-red-800"
                      : status.status === "uploading"
                      ? "bg-blue-100 text-blue-800"
                      : "bg-gray-100 text-gray-800"
                  }`}
                >
                  {status.status === "pending" && "Pending"}
                  {status.status === "uploading" && "Uploading..."}
                  {status.status === "success" && "Success"}
                  {status.status === "error" && "Failed"}
                </span>
              </div>

              {status.status === "uploading" && (
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div
                    className="bg-blue-600 h-2 rounded-full transition-all duration-300"
                    style={{ width: `${status.progress}%` }}
                  ></div>
                </div>
              )}

              {status.status === "error" && status.error && (
                <p className="text-red-500 text-xs mt-1">{status.error}</p>
              )}

              {status.status === "success" && (
                <div className="w-full bg-green-200 rounded-full h-2">
                  <div className="bg-green-600 h-2 rounded-full w-full"></div>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default UploadZone;
