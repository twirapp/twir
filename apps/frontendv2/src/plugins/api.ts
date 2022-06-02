import Axios, { AxiosError } from 'axios';

import { refreshAccessToken } from '@/functions/refreshAccessToken';

// eslint-disable-next-line import/no-cycle

export const api = Axios.create({
  baseURL: '/api',
});

api.interceptors.request.use(
  (config) => {
    const accessToken = localStorage.getItem('accessToken');

    return {
      ...config,
      headers: {
        ...config.headers,
        Authorization: `Bearer ${accessToken}`,
      },
    };
  },
  (error) => Promise.reject(error),
);

api.interceptors.response.use(
  (config) => config,
  async (error: AxiosError & { config: { __isRetryRequest: boolean } }) => {
    const response = error.response;

    if (response) {
      if (response.status === 401 && error.config && !error.config.__isRetryRequest) {
        try {
          await refreshAccessToken();
        } catch (authError) {
          return Promise.reject(error);
        }

        error.config.__isRetryRequest = true;
        return api(error.config);
      }
    }

    return Promise.reject(error);
  },
);
