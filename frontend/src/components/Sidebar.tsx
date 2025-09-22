import { NavLink, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { useTheme } from "../context/ThemeContext";

const Sidebar = () => {
  const { logout } = useAuth();
  const { theme, toggleTheme } = useTheme();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  return (
    <div className="w-64 h-screen bg-white dark:bg-gray-800 shadow-md flex flex-col">
      <div className="p-4 border-b dark:border-gray-700">
        <h1 className="text-2xl font-bold text-gray-800 dark:text-white">
          File Vault
        </h1>
      </div>
      <nav className="flex-grow p-4 space-y-2">
        <NavLink
          to="/dashboard"
          className={({ isActive }) =>
            `block px-4 py-2 rounded-md transition-colors ${
              isActive
                ? "bg-blue-500 text-white"
                : "text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
            }`
          }
        >
          📁 My Files
        </NavLink>
        <NavLink
          to="/stats"
          className={({ isActive }) =>
            `block px-4 py-2 rounded-md transition-colors ${
              isActive
                ? "bg-blue-500 text-white"
                : "text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
            }`
          }
        >
          📊 Statistics
        </NavLink>
        <NavLink
          to="/shared-with-me"
          className={({ isActive }) =>
            `block px-4 py-2 rounded-md transition-colors ${
              isActive
                ? "bg-blue-500 text-white"
                : "text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
            }`
          }
        >
          🤝 Shared With Me
        </NavLink>
        <NavLink
          to="/admin"
          className={({ isActive }) =>
            `block px-4 py-2 rounded-md transition-colors ${
              isActive
                ? "bg-blue-500 text-white"
                : "text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
            }`
          }
        >
          ⚙️ Admin
        </NavLink>
      </nav>
      <div className="p-4 border-t dark:border-gray-700 flex items-center justify-between">
        <button
          onClick={handleLogout}
          className="px-4 py-2 text-sm bg-red-500 text-white rounded hover:bg-red-600 transition-colors"
        >
          Logout
        </button>
        <button
          onClick={toggleTheme}
          className="p-2 rounded-full hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors"
        >
          {theme === "light" ? "🌙" : "☀️"}
        </button>
      </div>
    </div>
  );
};

export default Sidebar;
