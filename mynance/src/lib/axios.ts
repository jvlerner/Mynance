// lib/axios.ts
import axios from "axios";

let onUnauthorized: (() => void) | null = null;

export function setUnauthorizedHandler(callback: () => void) {
  onUnauthorized = callback;
}

const api = axios.create({
  baseURL: process.env.MYNANCE_API_HOST || "http://localhost:8080",
  withCredentials: true,
});

api.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401 && onUnauthorized) {
      onUnauthorized();
    }
    return Promise.reject(error);
  }
);

export default api;
