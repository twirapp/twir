import axios from 'axios';

import { redirectToLogin } from './redirectToLogin';

import { api } from '@/plugins/api';
import { setUser, User } from '@/stores/userStore';

export const fetchAndSetUser = async () => {
  const accessToken = localStorage.getItem('accessToken');
  const refreshToken = localStorage.getItem('refreshToken');

  if (!accessToken || !refreshToken) {
    return;
  }

  try {
    const profile = await api.get<User>('/auth/profile');
    setUser(profile.data);
  } catch (e) {
    if (axios.isAxiosError(e) && e.response?.status === 401 && (e.response.data as Record<string, any>).message === 'Missed scopes') {
      redirectToLogin();
    } else {
      console.error(e);
    }
  }
};
