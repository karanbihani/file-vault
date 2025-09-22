import React, { useState } from "react";
import { FiFilter, FiChevronDown, FiChevronUp } from "react-icons/fi";

export interface FilterOptions {
  filename?: string;
  mime_type?: string;
  min_size?: number;
  max_size?: number;
  start_date?: string;
  end_date?: string;
  tags?: string;
  uploader_email?: string;
}

interface AdvancedFiltersProps {
  onFiltersChange: (filters: FilterOptions) => void;
  onReset: () => void;
}

const AdvancedFilters: React.FC<AdvancedFiltersProps> = ({
  onFiltersChange,
  onReset,
}) => {
  const [filters, setFilters] = useState<FilterOptions>({});
  const [isExpanded, setIsExpanded] = useState(false);

  const commonMimeTypes = [
    { value: "", label: "All Types" },
    { value: "image/jpeg", label: "JPEG Images" },
    { value: "image/png", label: "PNG Images" },
    { value: "image/gif", label: "GIF Images" },
    { value: "application/pdf", label: "PDF Documents" },
    { value: "text/plain", label: "Text Files" },
    { value: "application/zip", label: "ZIP Archives" },
    { value: "video/mp4", label: "MP4 Videos" },
    { value: "audio/mpeg", label: "MP3 Audio" },
    { value: "application/json", label: "JSON Files" },
    { value: "text/csv", label: "CSV Files" },
    {
      value:
        "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
      label: "Word Documents",
    },
    {
      value:
        "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
      label: "Excel Spreadsheets",
    },
  ];

  const fileSizeOptions = [
    { value: "", label: "Any Size" },
    { value: "1024", label: "1 KB" },
    { value: "1048576", label: "1 MB" },
    { value: "10485760", label: "10 MB" },
    { value: "104857600", label: "100 MB" },
    { value: "1073741824", label: "1 GB" },
  ];

  const handleFilterChange = (
    key: keyof FilterOptions,
    value: string | number
  ) => {
    const newFilters = {
      ...filters,
      [key]: value === "" ? undefined : value,
    };
    setFilters(newFilters);
    onFiltersChange(newFilters);
  };

  const handleReset = () => {
    setFilters({});
    onReset();
  };

  const formatDateForInput = (date?: string) => {
    if (!date) return "";
    return new Date(date).toISOString().split("T")[0];
  };

  const hasActiveFilters = Object.values(filters).some(
    (value) => value !== undefined && value !== ""
  );

  return (
    <div className="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded-lg p-4 mb-6">
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 flex items-center">
          <FiFilter className="mr-2" />
          Advanced Filters
        </h3>
        <div className="flex items-center gap-2">
          {hasActiveFilters && (
            <span className="text-xs bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 px-2 py-1 rounded">
              {
                Object.values(filters).filter(
                  (v) => v !== undefined && v !== ""
                ).length
              }{" "}
              active
            </span>
          )}
          <button
            onClick={() => setIsExpanded(!isExpanded)}
            className="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-200 transition-colors flex items-center"
          >
            {isExpanded ? (
              <>
                <FiChevronUp className="mr-1" />
                Collapse
              </>
            ) : (
              <>
                <FiChevronDown className="mr-1" />
                Expand
              </>
            )}
          </button>
        </div>
      </div>

      {isExpanded && (
        <div className="space-y-4">
          {/* Filename Search */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Filename
            </label>
            <input
              type="text"
              value={filters.filename || ""}
              onChange={(e) => handleFilterChange("filename", e.target.value)}
              placeholder="Search by filename..."
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
            />
          </div>

          {/* MIME Type */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              File Type
            </label>
            <select
              value={filters.mime_type || ""}
              onChange={(e) => handleFilterChange("mime_type", e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
            >
              {commonMimeTypes.map((type) => (
                <option key={type.value} value={type.value}>
                  {type.label}
                </option>
              ))}
            </select>
          </div>

          {/* File Size Range */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Min Size
              </label>
              <select
                value={filters.min_size || ""}
                onChange={(e) =>
                  handleFilterChange(
                    "min_size",
                    e.target.value ? parseInt(e.target.value) : ""
                  )
                }
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
              >
                {fileSizeOptions.map((size) => (
                  <option key={`min-${size.value}`} value={size.value}>
                    {size.label}
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Max Size
              </label>
              <select
                value={filters.max_size || ""}
                onChange={(e) =>
                  handleFilterChange(
                    "max_size",
                    e.target.value ? parseInt(e.target.value) : ""
                  )
                }
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
              >
                {fileSizeOptions.map((size) => (
                  <option key={`max-${size.value}`} value={size.value}>
                    {size.label}
                  </option>
                ))}
              </select>
            </div>
          </div>

          {/* Date Range */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Upload Date From
              </label>
              <input
                type="date"
                value={formatDateForInput(filters.start_date)}
                onChange={(e) =>
                  handleFilterChange(
                    "start_date",
                    e.target.value ? new Date(e.target.value).toISOString() : ""
                  )
                }
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Upload Date To
              </label>
              <input
                type="date"
                value={formatDateForInput(filters.end_date)}
                onChange={(e) =>
                  handleFilterChange(
                    "end_date",
                    e.target.value ? new Date(e.target.value).toISOString() : ""
                  )
                }
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
              />
            </div>
          </div>

          {/* Tags */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Tags
            </label>
            <input
              type="text"
              value={filters.tags || ""}
              onChange={(e) => handleFilterChange("tags", e.target.value)}
              placeholder="Enter tags separated by commas (e.g., work, personal, important)"
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
            />
            <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
              Files must have ALL specified tags
            </p>
          </div>

          {/* Uploader Email */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Uploader Email
            </label>
            <input
              type="email"
              value={filters.uploader_email || ""}
              onChange={(e) =>
                handleFilterChange("uploader_email", e.target.value)
              }
              placeholder="Filter by uploader's email..."
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
            />
          </div>

          {/* Action Buttons */}
          <div className="flex gap-2 pt-2">
            <button
              onClick={handleReset}
              className="px-4 py-2 bg-gray-500 hover:bg-gray-600 text-white rounded-md transition-colors"
            >
              Clear All Filters
            </button>
            {hasActiveFilters && (
              <span className="flex items-center text-sm text-gray-600 dark:text-gray-400">
                Filters are applied automatically
              </span>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default AdvancedFilters;
