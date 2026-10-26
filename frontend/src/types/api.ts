/** Envelope every backend endpoint responds with (see backend pkg/utils/response.go). */
export interface ApiResponse<T> {
  success: boolean;
  message?: string;
  data: T;
  /** Field name to message, present on validation failures. */
  errors?: Record<string, string>;
}

/** Paginated list payload (backend utils.Page). */
export interface Page<T> {
  items: T[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

/** ISO-8601 timestamp string as serialized by the backend. */
export type ISODateString = string;
