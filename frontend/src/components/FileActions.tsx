import { useState, useEffect } from "react";
import apiClient from "../api/apiClient";
import Tooltip from "./Tooltip";
import { FiDownload, FiShare2, FiTrash2, FiBarChart2 } from "react-icons/fi";

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

interface PublicShare {
  share_token: string;
  download_count: number;
}

const FileActions = ({ file, onDelete }: FileActionsProps) => {
  const [showManageShare, setShowManageShare] = useState(false);
  const [shareEmail, setShareEmail] = useState("");
  const [recipients, setRecipients] = useState<ShareRecipient[]>([]);
  const [publicShare, setPublicShare] = useState<PublicShare | null>(null);

  useEffect(() => {
    if (showManageShare) {
      const fetchShareData = async () => {
        try {
          // Fetch private share recipients
          const recipientsResponse = await apiClient.get(
            `/files/${file.ID}/shares`
          );
          setRecipients(recipientsResponse.data || []);

          // Fetch public share info
          const publicResponse = await apiClient.get(
            `/files/${file.ID}/public-share`
          );
          setPublicShare(publicResponse.data || null);
        } catch (error) {
          console.error("Failed to fetch share data", error);
        }
      };
      fetchShareData();
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
        <Tooltip content="Download file">
          <button
            onClick={handleDownload}
            className="p-2 text-blue-500 hover:bg-blue-50 dark:hover:bg-blue-900/20 rounded-md transition-colors"
          >
            <FiDownload className="w-5 h-5" />
          </button>
        </Tooltip>

        <Tooltip content="Share file">
          <button
            onClick={() => setShowManageShare(true)}
            className="p-2 text-green-500 hover:bg-green-50 dark:hover:bg-green-900/20 rounded-md transition-colors"
          >
            <FiShare2 className="w-5 h-5" />
          </button>
        </Tooltip>

        <Tooltip content="Delete file">
          <button
            onClick={() => onDelete(file.ID)}
            className="p-2 text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-md transition-colors"
          >
            <FiTrash2 className="w-5 h-5" />
          </button>
        </Tooltip>
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
                className="w-full px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 mb-2"
              >
                Get Public Link
              </button>

              {publicShare && (
                <div className="p-3 bg-blue-50 dark:bg-blue-900/20 rounded-md">
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm font-medium text-blue-800 dark:text-blue-200">
                      Public Link Active
                    </span>
                    <span className="text-xs bg-blue-200 dark:bg-blue-800 px-2 py-1 rounded flex items-center">
                      <FiBarChart2 className="mr-1" />
                      {publicShare.download_count} downloads
                    </span>
                  </div>
                  <div className="text-xs text-blue-600 dark:text-blue-300 break-all">
                    {`${window.location.origin}/share/${publicShare.share_token}`}
                  </div>
                  <button
                    onClick={handleRevokePublic}
                    className="mt-2 text-xs text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300"
                  >
                    Revoke Public Link
                  </button>
                </div>
              )}
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
