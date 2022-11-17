// eslint-disable-next-line import/no-cycle
import axios from 'axios';

export const refreshAccessToken = async () => {
  const request = await axios.post<{
    accessToken: string;
  }>('/api/auth/token');
  const data = request.data;

  localStorage.setItem('accessToken', data.accessToken);
};
