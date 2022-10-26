import { accessTokenStore, refreshAccessToken } from '@/services/auth/token.js';

/**
 * Polifill for fetch function with bearer token authorization
 */
export const authFetch = async (
  url: RequestInfo | URL,
  options: RequestInit = {},
): Promise<Response> => {
  let isTryiedFetchToken = false;

  let accessToken = accessTokenStore.get();
  if (!accessToken) {
    const result = await refreshAccessToken();
    if (typeof result !== 'string') {
      accessTokenStore.set(undefined);
      return result;
    }

    accessToken = result;
    isTryiedFetchToken = true;
  }

  const { headers, ...opts } = options;

  const execute = async (token: string) => {
    return await fetch(url, {
      headers: new Headers({
        ...headers,
        Authorization: `Bearer ${token}`,
      }),
      ...opts,
    });
  };

  let response = await execute(accessToken);

  if (response.status === 401 && !isTryiedFetchToken) {
    const result = await refreshAccessToken();
    if (typeof result !== 'string') {
      accessTokenStore.set(undefined);
      return result;
    }

    accessToken = result;
    response = await execute(accessToken);
  }

  return response;
};
