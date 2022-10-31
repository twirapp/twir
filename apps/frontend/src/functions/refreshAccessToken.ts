// eslint-disable-next-line import/no-cycle
import axios from 'axios';

export const refreshAccessToken = async () => {
  const refreshToken = localStorage.getItem('refreshToken');

  if (!refreshToken) {
    throw new Error('Refresh token is empty.');
  }

  const request = await axios.post<{
    accessToken: string;
  }>('/api/auth/token', { refreshToken });
  const data = request.data;

  localStorage.setItem('accessToken', data.accessToken);
};
