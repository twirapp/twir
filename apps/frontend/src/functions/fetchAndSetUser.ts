import { AuthUser } from '@tsuwari/shared';
import axios from 'axios';

import { redirectToLogin } from './redirectToLogin';

import { api } from '@/plugins/api';
import { setUser } from '@/stores/userStore';

export const fetchAndSetUser = async () => {
  const accessToken = localStorage.getItem('accessToken');

  if (!accessToken) {
    return;
  }
  const profile = await api.get<AuthUser>('/auth/profile');
  setUser(profile.data);
};
