import { printError } from './error';

// Local storage key
const ACCESS_TOKEN_KEY = 'access_token';

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

export class FetcherError extends Error {
  messages?: string;

  constructor(data: string | Record<string, any>) {
    super(typeof data === 'string' ? data : 'Query error');
    if (typeof data === 'object') {
      this.messages = data.messages;
    }
  }
}

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

    const isJson = res.headers.get('content-type') == 'application/json';
    const data = isJson ? await res.json() : await res.text();

    if (!res.ok) {
      throw new FetcherError(data);
    }

    return data as T;
  };
};

export const fetcher = createFetcher();
export const authFetcher = createFetcher(authFetch);
