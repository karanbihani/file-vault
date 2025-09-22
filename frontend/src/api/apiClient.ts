import axios, { AxiosRequestConfig } from "axios";
import { rateLimiter } from "../utils/rateLimiter";

const baseClient = axios.create({
  baseURL: "http://localhost:8080/api/v1", // Your backend API URL
});

// Interceptor to add the auth token to every request
baseClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("authToken");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Rate-limited API client wrapper
const apiClient = {
  get: <T = any>(url: string, config?: AxiosRequestConfig) =>
    rateLimiter.execute(() => baseClient.get<T>(url, config)),

  post: <T = any>(url: string, data?: any, config?: AxiosRequestConfig) =>
    rateLimiter.execute(() => baseClient.post<T>(url, data, config)),

  // Special upload method with enhanced rate limiting
  upload: <T = any>(url: string, data?: any, config?: AxiosRequestConfig) =>
    rateLimiter.executeUpload(() => baseClient.post<T>(url, data, config)),

  put: <T = any>(url: string, data?: any, config?: AxiosRequestConfig) =>
    rateLimiter.execute(() => baseClient.put<T>(url, data, config)),

  delete: <T = any>(url: string, config?: AxiosRequestConfig) =>
    rateLimiter.execute(() => baseClient.delete<T>(url, config)),

  patch: <T = any>(url: string, data?: any, config?: AxiosRequestConfig) =>
    rateLimiter.execute(() => baseClient.patch<T>(url, data, config)),
};

export default apiClient;
