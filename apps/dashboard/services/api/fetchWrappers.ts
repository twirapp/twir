import { printError } from './error';

// Local storage key
const ACCESS_TOKEN_KEY = 'access_token';

export const mutationOptions = {
  revalidate: false,
  throwOnError: false,
};

/**
 * Wrapper for fetch function with bearer token authorization
 */
export const authFetch = async (
  url: RequestInfo | URL,
  options?: RequestInit,
): Promise<Response> => {
  // Boolean value reflecting whether there was an attempt
  // to refresh the token
  let isTryiedRefresh = false;

  let accessToken = localStorage.getItem(ACCESS_TOKEN_KEY);
  if (accessToken == null) {
    const result = await refreshAccessToken();
    if (typeof result != 'string') return result;

    accessToken = result;
    isTryiedRefresh = true;
  }

  const execute = async (token: string) => {
    return await fetch(url, {
      ...options,
      headers: {
        ...options?.headers,
        Authorization: `Bearer ${token}`,
      },
    });
  };

  let response = await execute(accessToken);

  if (response.status == 401 && !isTryiedRefresh) {
    const result = await refreshAccessToken();
    if (typeof result != 'string') return result;

    accessToken = result;
    response = await execute(accessToken);
  }

  return response;
};

/**
 * @returns Access token on success or Reponse object on error
 */
const refreshAccessToken = async (): Promise<Response | string> => {
  const res = await fetch('/api/auth/token', { method: 'POST' });

  if (!res.ok) {
    localStorage.removeItem(ACCESS_TOKEN_KEY);
    return res;
  }

  const accessToken = ((await res.json()) as { accessToken: string }).accessToken;
  localStorage.setItem(ACCESS_TOKEN_KEY, accessToken);
  return accessToken;
};

const createFetcher = (fetchFn = fetch) => {
  return async <T = any>(url: RequestInfo | URL, options?: RequestInit) => {
    const res = await fetchFn(url, options);

    if (res.headers.get('content-type') !== 'application/json') {
      throw new Error('Invalid data format, expected application/json');
    }

    const data = (await res.json()) as T;

    if (!res.ok) {
      printError((data as any).messages || (data as any).message);
      throw new Error(JSON.stringify(data));
    }

    return data;
  };
};

export const swrFetcher = createFetcher();
export const swrAuthFetcher = createFetcher(authFetch);
