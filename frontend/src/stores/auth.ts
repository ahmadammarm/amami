import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { User } from '../types/auth';
import apiClient from '../api/client';

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null);
  const token = ref<string | null>(localStorage.getItem('token'));
  const requiresPasswordChange = ref<boolean>(localStorage.getItem('requiresPasswordChange') === 'true');

  const isAuthenticated = computed(() => !!token.value);
  const userRole = computed(() => user.value?.role_name);

  async function fetchUser(): Promise<void> {
    if (!token.value) return;
    try {
      const response = await apiClient.get('/auth/me');
      // response.data.data because of the common response wrapper on backend
      user.value = response.data.data;
      
      // Sync requiresPasswordChange from user status if it's currently false
      if (user.value?.status === 'PENDING_PASSWORD_CHANGE') {
        requiresPasswordChange.value = true;
        localStorage.setItem('requiresPasswordChange', 'true');
      }
    } catch (error) {
      logout();
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
    try {
      const response = await apiClient.post('/auth/change-password', { new_password: newPassword });
      const { token: newToken } = response.data.data;
      setToken(newToken, false); // Clear the flag after success
      await fetchUser();
    } catch (error) {
      throw error;
    }
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
    fetchUser,
    setToken,
    changePassword,
    logout,
  };
});
