import axios from 'axios';

/** Why the app cannot talk to the server right now. */
export type ServiceIssue = 'maintenance' | 'unreachable' | 'offline';

const SERVICE_ISSUES: readonly ServiceIssue[] = ['maintenance', 'unreachable', 'offline'];

export const MAINTENANCE_ROUTE = 'Maintenance';

/**
 * Classifies a request failure: 503 means the server is deliberately down for
 * maintenance; no response at all means it is unreachable, or the device is offline.
 * Any other error is an ordinary request failure and returns null.
 */
export function classifyServiceError(error: unknown): ServiceIssue | null {
  if (!axios.isAxiosError(error)) return null;
  if (error.response?.status === 503) return 'maintenance';
  if (error.response) return null;
  if (error.code === 'ERR_CANCELED') return null;
  return typeof navigator !== 'undefined' && !navigator.onLine ? 'offline' : 'unreachable';
}

export function parseServiceIssue(value: unknown): ServiceIssue {
  return SERVICE_ISSUES.includes(value as ServiceIssue) ? (value as ServiceIssue) : 'unreachable';
}

/** Only same-app paths are accepted as a return target, never another origin. */
export function safeReturnPath(value: unknown): string {
  return typeof value === 'string' && value.startsWith('/') && !value.startsWith('//') ? value : '/';
}
