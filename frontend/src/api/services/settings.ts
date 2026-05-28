import apiClient from '../client';

export const getMosqueProfile = async () => {
  const response = await apiClient.get('/settings/profile');
  return response.data.data;
};

export const updateMosqueProfile = async (data: any) => {
  const response = await apiClient.put('/settings/profile', data);
  return response.data;
};

export const getSMTPConfig = async () => {
  const response = await apiClient.get('/settings/smtp');
  return response.data.data;
};

export const updateSMTPConfig = async (data: any) => {
  const response = await apiClient.put('/settings/smtp', data);
  return response.data;
};

export const testSMTPConnection = async (data: any) => {
  const response = await apiClient.post('/settings/smtp/test', data);
  return response.data;
};

export const getSystemHealth = async () => {
  const response = await apiClient.get('/settings/health');
  return response.data.data;
};

export const getAuditLogs = async (page: number = 1, limit: number = 10) => {
  const response = await apiClient.get('/settings/audit', { params: { page, limit } });
  return response.data.data;
};
