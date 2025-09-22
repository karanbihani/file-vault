import { useState } from "react";

interface TagManagerProps {
  tags: string[];
  onAddTag: (tag: string) => void;
  onRemoveTag: (tag: string) => void;
}

const TagManager = ({ tags, onAddTag, onRemoveTag }: TagManagerProps) => {
  const [newTag, setNewTag] = useState("");

  const handleAdd = (e: React.FormEvent) => {
    e.preventDefault();
    if (newTag && !tags.includes(newTag)) {
      onAddTag(newTag);
      setNewTag("");
    }
  };

  return (
    <div>
      <div className="flex flex-wrap gap-2 mb-2">
        {tags &&
          tags.map((tag) => (
            <span
              key={tag}
              className="flex items-center bg-blue-100 text-blue-800 text-xs font-medium px-2.5 py-0.5 rounded dark:bg-blue-900 dark:text-blue-300"
            >
              {tag}
              <button
                onClick={() => onRemoveTag(tag)}
                className="ml-2 text-blue-400 hover:text-blue-600 font-bold"
              >
                &times;
              </button>
            </span>
          ))}
      </div>
      <form onSubmit={handleAdd} className="flex gap-2">
        <input
          type="text"
          value={newTag}
          onChange={(e) => setNewTag(e.target.value)}
          placeholder="Add a tag..."
          className="p-1 border rounded-md text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white w-full"
        />
        <button
          type="submit"
          className="px-3 py-1 bg-green-500 text-white text-sm rounded-md hover:bg-green-600"
        >
          Add
        </button>
      </form>
    </div>
  );
};

export default TagManager;
