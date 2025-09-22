import { useState } from "react";
import apiClient from "../api/apiClient";
import { useNavigate } from "react-router-dom";
import { AxiosError } from "axios";

const Register = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState(""); // State for error messages
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    try {
      await apiClient.post("/register", { email, password });
      alert("Registration successful! Please log in.");
      navigate("/login");
    } catch (err) {
      console.error("Registration failed", err);
      if (err instanceof AxiosError && err.response) {
        setError(
          err.response.data.error || "Registration failed. Please try again."
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
        <h2 className="text-2xl font-bold mb-4">Register</h2>
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
          className="w-full p-2 text-white bg-green-500 rounded"
        >
          Register
        </button>
      </form>
    </div>
  );
};

export default Register;
