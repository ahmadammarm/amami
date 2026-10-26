import type { Permission } from '../lib/permissions';

export type UserRole = 'SUPER_ADMIN' | 'BENDAHARA' | 'TAKMIR' | 'SEKRETARIS' | 'JAMAAH';

export type UserStatus = 'ACTIVE' | 'SUSPENDED' | 'PENDING_PASSWORD_CHANGE';

export interface User {
  id: string;
  username: string;
  email: string;
  full_name: string;
  role_name: UserRole;
  status: UserStatus;
  permissions: Permission[];
}

export interface LoginResponse {
  token: string;
  requires_password_change: boolean;
}
