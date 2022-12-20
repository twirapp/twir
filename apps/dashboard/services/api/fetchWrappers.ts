import { printError } from './error';

// Local storage key
const ACCESS_TOKEN_KEY = 'access_token';

export const mutationOptions = {
  revalidate: false,
  throwOnError: false,
};

/**
 * Polifill for fetch function with bearer token authorization
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

export const swrAuthFetcher = async (url: RequestInfo | URL, options?: RequestInit) => {
  const req = await authFetch(url, options);

  const isJson = req.headers.get('content-type') === 'application/json';
  const response = isJson ? await req.json() : await req.text();
  if (!req.ok) {
    printError(response.messages || response.message);
    throw new Error(response);
  }

  return response;
};

export const swrFetcher = async (url: RequestInfo | URL, options?: RequestInit) => {
  const req = await fetch(url, options);

  const isJson = req.headers.get('content-type') === 'application/json';
  const response = isJson ? await req.json() : await req.text();
  if (!req.ok) {
    printError(response.messages || response.message);
    throw new Error(response);
  }

  return response;
};
