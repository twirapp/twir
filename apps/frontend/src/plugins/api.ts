import Axios, { AxiosError } from 'axios';
import { useToast } from 'vue-toastification';

import { router } from './router';

import { refreshAccessToken } from '@/functions/refreshAccessToken';

// eslint-disable-next-line import/no-cycle

export const api = Axios.create({
  baseURL: '/api',
});

const toast = useToast();

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
          localStorage.removeItem('accessToken');

          if (router.currentRoute.value.fullPath.startsWith('/dashboard')) {
            router.push('/');
          }
          return Promise.reject(error);
        }

        error.config.__isRetryRequest = true;
        return api(error.config);
      }

      const data = response.data as any;
      if (data?.message) {
        if (data.message !== 'Invalid Token!') {
          toast.error(Array.isArray(data.message) ? data.message.join(', ') : data.message);
        }

        return Promise.reject(error);
      }
    }

    return Promise.reject(error);
  },
);
