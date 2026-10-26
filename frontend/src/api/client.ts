import axios, { type AxiosInstance, type AxiosRequestConfig, type InternalAxiosRequestConfig } from 'axios';
import type { ApiResponse } from '../types/api';

const API_BASE_URL: string = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';
const LOGIN_PATH = '/auth/login';
const SESSION_PATH = '/auth/me';

/**
 * 401s that must not trigger the global redirect: a failed login is a form
 * error, and the session check is resolved by the router guard (public pages
 * such as the landing page stay open when a stored token has expired).
 */
const NO_REDIRECT_PATHS = new Set([LOGIN_PATH, SESSION_PATH]);

const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request Interceptor: Inject JWT Token
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('token');
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error: Error) => {
    return Promise.reject(error);
  }
);

// Response Interceptor: an expired session on any other request sends the user to login.
apiClient.interceptors.response.use(
  (response) => response,
  (error: Error) => {
    if (axios.isAxiosError(error) && error.response?.status === 401 && !NO_REDIRECT_PATHS.has(error.config?.url ?? '')) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

/** Sends a request and returns the unwrapped `data` of the response envelope. */
export async function request<T>(config: AxiosRequestConfig): Promise<T> {
  const response = await apiClient.request<ApiResponse<T>>(config);
  return response.data.data;
}

/** Extracts a user-facing message from any error thrown by `request`. */
export function getErrorMessage(error: unknown, fallback = 'Terjadi kesalahan, silakan coba lagi'): string {
  if (axios.isAxiosError<Partial<ApiResponse<unknown>>>(error)) {
    const body = error.response?.data;
    const firstFieldError = body?.errors ? Object.values(body.errors)[0] : undefined;
    return firstFieldError ?? body?.message ?? fallback;
  }
  return error instanceof Error ? error.message : fallback;
}

export function isNotFoundError(error: unknown): boolean {
  return axios.isAxiosError(error) && error.response?.status === 404;
}

/** Turns a server-relative path (e.g. "/uploads/...") into an absolute URL on the API host. */
export function assetUrl(path: string): string {
  return new URL(path, new URL(API_BASE_URL).origin).toString();
}

export default apiClient;
