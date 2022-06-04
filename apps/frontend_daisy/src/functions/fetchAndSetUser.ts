import { api } from '@/plugins/api';
import { setUser, User } from '@/stores/userStore';

export const fetchAndSetUser = async () => {
  const accessToken = localStorage.getItem('accessToken');
  const refreshToken = localStorage.getItem('refreshToken');

  if (!accessToken || !refreshToken) {
    return;
  }

  const profile = await api.get<User>('/auth/profile');
  setUser(profile.data);
};
