import { request } from '../client';
import type { Page } from '../../types/api';
import type { InviteUserRequest, InviteUserResponse, ManageableUserStatus, Role, UserResponse } from '../../types/user';

export const userService = {
  getAllUsers(page = 1, limit = 100): Promise<Page<UserResponse>> {
    return request<Page<UserResponse>>({ method: 'GET', url: '/users', params: { page, limit } });
  },

  getRoles(): Promise<Role[]> {
    return request<Role[]>({ method: 'GET', url: '/users/roles' });
  },

  inviteUser(data: InviteUserRequest): Promise<InviteUserResponse> {
    return request<InviteUserResponse>({ method: 'POST', url: '/users/invite', data });
  },

  updateUserStatus(id: string, status: ManageableUserStatus): Promise<void> {
    return request<void>({ method: 'PATCH', url: `/users/${id}/status`, data: { status } });
  },

  deleteUser(id: string): Promise<void> {
    return request<void>({ method: 'DELETE', url: `/users/${id}` });
  },
};
