// eslint-disable-next-line import/no-cycle
import { api } from '@/plugins/api';

export const refreshAccessToken = async () => {
  const refreshToken = localStorage.getItem('refreshToken');

  if (!refreshToken) {
    throw new Error('Refresh token is empty.');
  }

  try {
    const request = await api.post<{
      accessToken: string;
    }>('/auth/token', { refreshToken });
    const data = request.data;

    if (request.status === 401) {
      localStorage.removeItem('accessToken');
      localStorage.removeItem('refreshToken');
    } else {
      localStorage.setItem('accessToken', data.accessToken);
    }
  } catch (e) {
    localStorage.removeItem('accessToken');
    localStorage.removeItem('refreshToken');
  }
};
