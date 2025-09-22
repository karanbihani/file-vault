import { useState, useEffect } from "react";
import apiClient from "../api/apiClient";

interface UserFile {
  ID: number;
  Filename: string;
}

interface FileActionsProps {
  file: UserFile;
  onDelete: (fileId: number) => void;
}

interface ShareRecipient {
  id: number;
  email: string;
}

const FileActions = ({ file, onDelete }: FileActionsProps) => {
  const [showManageShare, setShowManageShare] = useState(false);
  const [shareEmail, setShareEmail] = useState("");
  const [recipients, setRecipients] = useState<ShareRecipient[]>([]);

  useEffect(() => {
    if (showManageShare) {
      const fetchRecipients = async () => {
        try {
          const response = await apiClient.get(`/files/${file.ID}/shares`);
          setRecipients(response.data || []);
        } catch (error) {
          console.error("Failed to fetch share recipients", error);
        }
      };
      fetchRecipients();
    }
  }, [showManageShare, file.ID]);

  const handleDownload = async () => {
    try {
      const response = await apiClient.get(`/files/${file.ID}/download`, {
        responseType: "blob", // Important: Tell Axios to expect binary data
      });
      // Create a temporary URL from the blob data
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement("a");
      link.href = url;
      link.setAttribute("download", file.Filename); // Set the filename
      document.body.appendChild(link);
      link.click();
      link.remove(); // Clean up the temporary link
      window.URL.revokeObjectURL(url); // Clean up the temporary URL
    } catch (error) {
      console.error("Download failed", error);
      alert("Could not download the file.");
    }
  };

  const handleSharePublic = async () => {
    try {
      const response = await apiClient.post(`/files/${file.ID}/share`);
      prompt("Share this public link:", response.data.share_url);
    } catch (error) {
      console.error("Failed to create public share link", error);
      alert("Could not create a public share link.");
    }
  };

  const handleShareWithUser = async () => {
    if (!shareEmail) {
      alert("Please enter an email address.");
      return;
    }
    try {
      await apiClient.post(`/files/${file.ID}/share-to-user`, {
        email: shareEmail,
      });
      alert(`File shared successfully with ${shareEmail}`);
      setShowManageShare(false);
      setShareEmail("");
      // Refresh recipients list
      const response = await apiClient.get(`/files/${file.ID}/shares`);
      setRecipients(response.data || []);
    } catch (error) {
      console.error("Failed to share with user", error);
      alert("Failed to share file. Please check the email and try again.");
    }
  };

  const handleUnshare = async (recipientId: number) => {
    try {
      await apiClient.delete(`/files/${file.ID}/share-to-user`, {
        data: { recipient_id: recipientId },
      });
      setRecipients(recipients.filter((r) => r.id !== recipientId)); // Update UI immediately
      alert("Successfully unshared file.");
    } catch (error) {
      console.error("Failed to unshare file", error);
      alert("Failed to unshare file.");
    }
  };

  const handleRevokePublic = async () => {
    try {
      await apiClient.delete(`/files/${file.ID}/share`);
      alert("All public links for this file have been revoked.");
    } catch (error) {
      console.error("Failed to revoke public links", error);
      alert("Failed to revoke public links.");
    }
  };

  return (
    <>
      <div className="flex gap-2">
        <button
          onClick={handleDownload}
          className="px-3 py-1 text-sm bg-blue-500 text-white rounded hover:bg-blue-600"
        >
          Download
        </button>
        <button
          onClick={() => setShowManageShare(true)}
          className="px-3 py-1 text-sm bg-green-500 text-white rounded hover:bg-green-600"
        >
          Share
        </button>
        <button
          onClick={() => onDelete(file.ID)}
          className="px-3 py-1 text-sm bg-red-500 text-white rounded hover:bg-red-600"
        >
          Delete
        </button>
      </div>

      {/* Manage Share Modal */}
      {showManageShare && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-xl max-w-md w-full mx-4">
            <h3 className="text-lg font-bold mb-4 text-gray-800 dark:text-white">
              Manage Sharing for "{file.Filename}"
            </h3>

            <div className="mb-4">
              <button
                onClick={handleSharePublic}
                className="w-full px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
              >
                Get Public Link
              </button>
            </div>

            <div className="border-t pt-4">
              <label className="block text-sm font-medium mb-1 text-gray-700 dark:text-gray-300">
                Share with a specific user:
              </label>
              <div className="flex gap-2">
                <input
                  type="email"
                  value={shareEmail}
                  onChange={(e) => setShareEmail(e.target.value)}
                  placeholder="user@example.com"
                  className="flex-grow p-2 border rounded-md dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                />
                <button
                  onClick={handleShareWithUser}
                  className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
                >
                  Share
                </button>
              </div>
            </div>

            {recipients.length > 0 && (
              <div className="border-t pt-4 mt-4">
                <h4 className="font-semibold mb-2 text-gray-800 dark:text-white">
                  Shared With:
                </h4>
                <ul className="space-y-2">
                  {recipients.map((r) => (
                    <li
                      key={r.id}
                      className="flex justify-between items-center p-2 bg-gray-50 dark:bg-gray-700 rounded"
                    >
                      <span className="text-gray-800 dark:text-gray-200">
                        {r.email}
                      </span>
                      <button
                        onClick={() => handleUnshare(r.id)}
                        className="px-2 py-1 text-xs bg-red-500 text-white rounded hover:bg-red-600"
                      >
                        Unshare
                      </button>
                    </li>
                  ))}
                </ul>
              </div>
            )}

            <div className="border-t pt-4 mt-4">
              <button
                onClick={handleRevokePublic}
                className="w-full px-4 py-2 bg-yellow-500 text-white rounded hover:bg-yellow-600 mb-2"
              >
                Revoke All Public Links
              </button>
              <button
                onClick={() => setShowManageShare(false)}
                className="w-full px-4 py-2 bg-gray-300 dark:bg-gray-600 text-gray-800 dark:text-white rounded hover:bg-gray-400 dark:hover:bg-gray-500"
              >
                Close
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
};

export default FileActions;
