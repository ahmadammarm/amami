export interface UserResponse {
  id: string;
  username: string;
  email: string;
  full_name: string;
  role_name: string;
  status: 'ACTIVE' | 'SUSPENDED' | 'PENDING_PASSWORD_CHANGE';
}

export interface InviteUserRequest {
  username: string;
  email: string;
  full_name: string;
  role_id: number;
  password: string;
}

export interface UpdateStatusRequest {
  status: 'ACTIVE' | 'SUSPENDED';
}

export interface Role {
  id: number;
  name: string;
}
