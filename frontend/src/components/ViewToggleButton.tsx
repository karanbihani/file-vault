interface ViewToggleButtonProps {
  view: "grid" | "list";
  setView: (view: "grid" | "list") => void;
}

const ViewToggleButton = ({ view, setView }: ViewToggleButtonProps) => {
  return (
    <div className="flex border rounded-md">
      <button
        onClick={() => setView("list")}
        className={`px-3 py-1 ${
          view === "list" ? "bg-blue-500 text-white" : "bg-white"
        }`}
      >
        List
      </button>
      <button
        onClick={() => setView("grid")}
        className={`px-3 py-1 ${
          view === "grid" ? "bg-blue-500 text-white" : "bg-white"
        }`}
      >
        Grid
      </button>
    </div>
  );
};

export default ViewToggleButton;
