import { createEventHook, isClient, useFetch } from '@vueuse/core';
import { ref } from 'vue';

import { accessTokenStore } from '@/stores/user.js';

const loginRouteState = isClient ? window.btoa(window.location.origin + '/login') : '';
const originState = isClient ? window.btoa(window.location.origin) : '';

export const redirectToLogin = () => {
  if (isClient) {
    window.location.replace(`/api/auth?state=${originState})}`);
  }
};

export const redirectToDashboard = () => {
  window.location.replace('/app/dashboard');
};

export const refreshAccessToken = async (): Promise<Response | string> => {
  const res = await fetch('/api/auth/token', { method: 'post' });
  if (!res.ok) return res;

  const data = (await res.json()) as { accessToken: string };
  accessTokenStore.set(data.accessToken);
  return data.accessToken;
};

export const twitchLoginHook = (params: URLSearchParams) => {
  const code = params.get('code');
  const isError = ref<boolean>(params.get('error') !== null);
  const successLoginHook = createEventHook<Response>();
  const error = ref<any | null>(null);

  if (!code) isError.value = true;

  const searchParams = new URLSearchParams({
    code: code as string,
    state: loginRouteState,
  });

  const { onFetchResponse, isFetching, execute, data } = useFetch(
    '/api/auth/token?' + searchParams,
    {
      immediate: false,
    },
  ).json<{ accessToken?: string }>();

  if (!isError.value) {
    execute();

    onFetchResponse(async (response) => {
      if (!data.value || !response.ok) {
        isError.value = true;
        error.value = 'Unexpected error';
        return;
      }

      const token = data.value.accessToken;
      if (!token) {
        isError.value = true;
        error.value = 'AccessToken is not found';
        return;
      }

      accessTokenStore.set(token);
      successLoginHook.trigger(response);
    });
  }

  return { isError, error, onSuccessLogin: successLoginHook.on, isFetching };
};
