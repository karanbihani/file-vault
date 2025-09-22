interface UserFile {
  ID: number;
  Filename: string;
  MimeType: string;
}

interface FilePreviewModalProps {
  file: UserFile | null;
  onClose: () => void;
}

const FilePreviewModal = ({ file, onClose }: FilePreviewModalProps) => {
  if (!file) return null;

  const isImage = file.MimeType.startsWith("image/");
  const previewUrl = `http://localhost:8080/api/v1/files/${file.ID}/download`;

  return (
    <div
      className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center"
      onClick={onClose}
    >
      <div
        className="bg-white p-4 rounded-lg max-w-3xl max-h-[80vh]"
        onClick={(e) => e.stopPropagation()}
      >
        <h3 className="text-lg font-bold mb-2">{file.Filename}</h3>
        {isImage ? (
          <img
            src={previewUrl}
            alt={`Preview of ${file.Filename}`}
            className="max-w-full max-h-[70vh] object-contain"
          />
        ) : (
          <p>No preview available for this file type.</p>
        )}
        <button
          onClick={onClose}
          className="mt-4 px-4 py-2 bg-red-500 text-white rounded"
        >
          Close
        </button>
      </div>
    </div>
  );
};

export default FilePreviewModal;
