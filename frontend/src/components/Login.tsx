import { useState } from "react";
import { useAuth } from "../context/AuthContext";
import apiClient from "../api/apiClient";
import { useNavigate } from "react-router-dom";
import { AxiosError } from "axios";

const Login = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState(""); // State for error messages
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(""); // Clear previous errors
    try {
      const response = await apiClient.post("/login", { email, password });
      login(response.data.token);
      navigate("/dashboard");
    } catch (err) {
      console.error("Login failed", err);
      if (err instanceof AxiosError && err.response) {
        // Use the error message from the backend if available
        setError(
          err.response.data.error || "Invalid credentials. Please try again."
        );
      } else {
        setError("An unexpected error occurred. Please try again.");
      }
    }
  };

  return (
    <div className="flex items-center justify-center h-screen bg-gray-100">
      <form
        onSubmit={handleSubmit}
        className="p-8 bg-white rounded shadow-md w-96"
      >
        <h2 className="text-2xl font-bold mb-4">Login</h2>
        {error && <p className="text-red-500 mb-4">{error}</p>}
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="Email"
          className="w-full p-2 mb-4 border rounded"
        />
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder="Password"
          className="w-full p-2 mb-4 border rounded"
        />
        <button
          type="submit"
          className="w-full p-2 text-white bg-blue-500 rounded"
        >
          Login
        </button>
      </form>
    </div>
  );
};

export default Login;
