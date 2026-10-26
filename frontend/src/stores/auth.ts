import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { LoginResponse, User } from '../types/auth';
import type { Permission } from '../lib/permissions';
import { request } from '../api/client';
import { classifyServiceError } from '../lib/serviceStatus';

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null);
  const token = ref<string | null>(localStorage.getItem('token'));
  const requiresPasswordChange = ref<boolean>(localStorage.getItem('requiresPasswordChange') === 'true');

  const isAuthenticated = computed(() => !!token.value);
  const userRole = computed(() => user.value?.role_name);
  const permissions = computed(() => new Set<Permission>(user.value?.permissions ?? []));

  function can(permission: Permission): boolean {
    return permissions.value.has(permission);
  }

  function canAny(...required: Permission[]): boolean {
    return required.some(can);
  }

  async function fetchUser(): Promise<void> {
    if (!token.value) return;
    try {
      user.value = await request<User>({ method: 'GET', url: '/auth/me' });

      // Sync requiresPasswordChange from user status if it's currently false
      if (user.value.status === 'PENDING_PASSWORD_CHANGE') {
        requiresPasswordChange.value = true;
        localStorage.setItem('requiresPasswordChange', 'true');
      }
    } catch (error) {
      // Keep the session through an outage; only a rejected session signs out.
      if (!classifyServiceError(error)) logout();
      throw error;
    }
  }

  function setToken(newToken: string, changeRequired = false): void {
    token.value = newToken;
    requiresPasswordChange.value = changeRequired;
    localStorage.setItem('token', newToken);
    localStorage.setItem('requiresPasswordChange', String(changeRequired));
  }

  async function changePassword(newPassword: string): Promise<void> {
    const { token: newToken } = await request<LoginResponse>({
      method: 'POST',
      url: '/auth/change-password',
      data: { new_password: newPassword },
    });
    setToken(newToken, false);
    await fetchUser();
  }

  function logout(): void {
    user.value = null;
    token.value = null;
    requiresPasswordChange.value = false;
    localStorage.removeItem('token');
    localStorage.removeItem('requiresPasswordChange');
  }

  return {
    user,
    token,
    requiresPasswordChange,
    isAuthenticated,
    userRole,
    can,
    canAny,
    fetchUser,
    setToken,
    changePassword,
    logout,
  };
});
