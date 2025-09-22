import { Navigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { useEffect, useState } from "react";
import apiClient from "../api/apiClient";

// This component checks if the logged-in user has admin permissions.
const AdminRoute = ({ children }: { children: JSX.Element }) => {
  const { isAuthenticated } = useAuth();
  const [isAdmin, setIsAdmin] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const checkAdminStatus = async () => {
      try {
        // We'll leverage the admin endpoint. If it succeeds, the user is an admin.
        await apiClient.get("/admin/stats");
        setIsAdmin(true);
      } catch (error) {
        setIsAdmin(false);
      } finally {
        setIsLoading(false);
      }
    };

    if (isAuthenticated) {
      checkAdminStatus();
    } else {
      setIsLoading(false);
    }
  }, [isAuthenticated]);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-50 dark:bg-gray-900">
        <div className="text-gray-600 dark:text-gray-300">Loading...</div>
      </div>
    );
  }

  return isAuthenticated && isAdmin ? children : <Navigate to="/dashboard" />;
};

export default AdminRoute;
