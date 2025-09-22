import { Link } from "react-router-dom";

const NotFound = () => {
  return (
    <div className="flex flex-col items-center justify-center h-screen text-center bg-gray-100">
      <h1 className="text-6xl font-bold text-gray-800">404</h1>
      <p className="text-xl text-gray-600 mb-4">Page Not Found</p>
      <Link
        to="/"
        className="px-4 py-2 text-white bg-blue-500 rounded hover:bg-blue-600"
      >
        Go Home
      </Link>
    </div>
  );
};

export default NotFound;
