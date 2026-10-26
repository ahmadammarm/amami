import type { UserRole, UserStatus } from './auth';

export interface UserResponse {
  id: string;
  username: string;
  email: string;
  full_name: string;
  role_name: UserRole;
  status: UserStatus;
}

export interface InviteUserRequest {
  username: string;
  email: string;
  full_name: string;
  role_id: number;
  password: string;
}

export interface InviteUserResponse {
  user: UserResponse;
  temporary_password: string;
}

export type ManageableUserStatus = Extract<UserStatus, 'ACTIVE' | 'SUSPENDED'>;

export interface UpdateStatusRequest {
  status: ManageableUserStatus;
}

export interface Role {
  id: number;
  name: UserRole;
  description: string;
}
