import { request } from '../client';
import type {
  AuditLogList,
  MosqueProfile,
  MosqueProfilePayload,
  SMTPConfig,
  SMTPConfigPayload,
  SystemHealth,
  TestSMTPPayload,
} from '../../types/settings';

export const getMosqueProfile = (): Promise<MosqueProfile> =>
  request<MosqueProfile>({ method: 'GET', url: '/settings/profile' });

export const updateMosqueProfile = (data: MosqueProfilePayload): Promise<void> =>
  request<void>({ method: 'PUT', url: '/settings/profile', data });

export const getSMTPConfig = (): Promise<SMTPConfig> =>
  request<SMTPConfig>({ method: 'GET', url: '/settings/smtp' });

export const updateSMTPConfig = (data: SMTPConfigPayload): Promise<void> =>
  request<void>({ method: 'PUT', url: '/settings/smtp', data });

export const testSMTPConnection = (data: TestSMTPPayload): Promise<void> =>
  request<void>({ method: 'POST', url: '/settings/smtp/test', data });

export const getSystemHealth = (): Promise<SystemHealth> =>
  request<SystemHealth>({ method: 'GET', url: '/settings/health' });

export const getAuditLogs = (page = 1, limit = 10): Promise<AuditLogList> =>
  request<AuditLogList>({ method: 'GET', url: '/settings/audit', params: { page, limit } });
