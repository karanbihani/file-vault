import { Link } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { useTheme } from "../context/ThemeContext";

const Navbar = () => {
  const { isAuthenticated, logout } = useAuth();
  const { theme, toggleTheme } = useTheme();

  return (
    <nav className="bg-gray-800 p-4 text-white">
      <div className="container mx-auto flex justify-between">
        <Link to="/" className="font-bold">
          File Vault
        </Link>
        <div>
          {isAuthenticated ? (
            <>
              <Link to="/dashboard" className="ml-4">
                My Files
              </Link>
              <Link to="/shared-with-me" className="ml-4">
                Shared With Me
              </Link>
              <Link
                to="/admin"
                className="ml-4 text-yellow-400 hover:text-yellow-300"
              >
                Admin
              </Link>
              <button onClick={logout} className="ml-4">
                Logout
              </button>
            </>
          ) : (
            <>
              <Link to="/login" className="ml-4">
                Login
              </Link>
              <Link to="/register" className="ml-4">
                Register
              </Link>
            </>
          )}
          <button
            onClick={toggleTheme}
            className="ml-4 p-2 rounded-full hover:bg-gray-700"
          >
            {theme === "light" ? "🌙" : "☀️"}
          </button>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
