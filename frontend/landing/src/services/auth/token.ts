import { persistentAtom } from '@nanostores/persistent';

import { logout, postRefreshToken } from './api.js';

export const accessTokenStore = persistentAtom('access_token');

/**
 * @returns 
 * Access token on success and Reponse object on error
 */
export const refreshAccessToken = async (): Promise<Response | string> => {
  const res = await postRefreshToken();
  if (!res.ok) return res;

  const { accessToken } = (await res.json()) as { accessToken: string };
  if (!accessToken) return res;
  accessTokenStore.set(accessToken);
  return accessToken;
};

export const logoutAndRemoveToken = async () => {
  const isLoggedOut = await logout();
  if (isLoggedOut) {
    accessTokenStore.set(undefined);
  }
  return isLoggedOut;
};
