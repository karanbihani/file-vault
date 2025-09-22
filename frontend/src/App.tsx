import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import { AuthProvider } from "./context/AuthContext";
import Login from "./components/Login";
import Register from "./components/Register";
import Dashboard from "./components/Dashboard";
import SharedWithMe from "./components/SharedWithMe";
import AdminDashboard from "./components/AdminDashboard";
import ProtectedRoute from "./components/ProtectedRoute";
import AdminRoute from "./components/AdminRoute";
import PublicRoute from "./components/PublicRoute"; // <-- Import PublicRoute
import Navbar from "./components/Navbar";
import NotFound from "./components/NotFound"; // <-- Import NotFound
import { ThemeProvider } from "./context/ThemeContext.tsx";

function App() {
  return (
    <AuthProvider>
      <ThemeProvider>
        <Router>
          <Navbar />
          <Routes>
            <Route
              path="/login"
              element={
                <PublicRoute>
                  <Login />
                </PublicRoute>
              }
            />
            <Route
              path="/register"
              element={
                <PublicRoute>
                  <Register />
                </PublicRoute>
              }
            />
            <Route
              path="/dashboard"
              element={
                <ProtectedRoute>
                  <Dashboard />
                </ProtectedRoute>
              }
            />
            <Route
              path="/shared-with-me"
              element={
                <ProtectedRoute>
                  <SharedWithMe />
                </ProtectedRoute>
              }
            />
            <Route
              path="/admin"
              element={
                <AdminRoute>
                  <AdminDashboard />
                </AdminRoute>
              }
            />
            <Route path="/" element={<Login />} />
            <Route path="*" element={<NotFound />} />{" "}
            {/* <-- Catch-all for 404 */}
          </Routes>
        </Router>
      </ThemeProvider>
    </AuthProvider>
  );
}

export default App;
