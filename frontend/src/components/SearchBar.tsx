import { useState, useEffect, useCallback } from "react";
import { debounce } from "../utils/rateLimiter";

interface SearchBarProps {
  onSearch: (searchTerm: string) => void;
}

const SearchBar = ({ onSearch }: SearchBarProps) => {
  const [searchTerm, setSearchTerm] = useState("");

  // Create a debounced search function with increased delay to respect rate limits
  const debouncedSearch = useCallback(
    debounce((term: string) => onSearch(term), 800), // 800ms delay
    [onSearch]
  );

  useEffect(() => {
    debouncedSearch(searchTerm);
  }, [searchTerm, debouncedSearch]);

  return (
    <div className="mb-6">
      <input
        type="text"
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
        placeholder="Search by filename, tags, date, or any property..."
        className="w-full p-3 border rounded-md dark:bg-gray-700 dark:border-gray-600 dark:text-white placeholder-gray-400 dark:placeholder-gray-300"
      />
      {searchTerm && (
        <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
          Searching for: "{searchTerm}"
        </p>
      )}
    </div>
  );
};

export default SearchBar;
