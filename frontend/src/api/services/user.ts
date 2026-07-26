import apiClient from '../client';
import type { InviteUserRequest } from '../../types/user';

export const userService = {
  getAllUsers: async (page = 1, limit = 100) => {
    const response = await apiClient.get(`/users?page=${page}&limit=${limit}`);
    return response.data.data;
  },

  inviteUser: async (data: InviteUserRequest) => {
    const response = await apiClient.post('/users/invite', data);
    return response.data.data;
  },

  updateUserStatus: async (id: string, status: string) => {
    const response = await apiClient.patch(`/users/${id}/status`, { status });
    return response.data.data;
  },

  deleteUser: async (id: string) => {
    const response = await apiClient.delete(`/users/${id}`);
    return response.data.data;
  },

  getRoles: async () => {
    // We don't have a dedicated roles endpoint yet, but for now we can hardcode 
    // or wait for a specific endpoint. Let's assume we might need one or just 
    // fetch them if we had one. Actually, let's just hardcode the IDs from seeder for now
    // or better, create a small endpoint if missing.
    // For this plan, I'll assume we can use the ID from seeder:
    // 1: SUPER_ADMIN, 2: BENDAHARA, 3: TAKMIR, 4: SEKRETARIS, 5: JAMAAH
    return [
      { id: 1, name: 'SUPER_ADMIN' },
      { id: 2, name: 'BENDAHARA' },
      { id: 3, name: 'TAKMIR' },
      { id: 4, name: 'SEKRETARIS' },
      { id: 5, name: 'JAMAAH' },
    ];
  }
};
