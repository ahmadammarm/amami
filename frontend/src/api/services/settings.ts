import apiClient from '../client';

export const getMosqueProfile = async () => {
  const response = await apiClient.get('/settings/profile');
  return response.data.data;
};
