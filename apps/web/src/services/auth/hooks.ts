import { useQuery } from '@tanstack/vue-query';

import { getProfile } from './api.js';
import { redirectToDashboard } from './locationHelpers.js';
import { accessTokenStore } from './token.js';
import { handleTwitchLoginCallback } from './twitch.js';

export const useUserProfile = () =>
  useQuery(['profile'], getProfile, {
    retry: false,
    refetchOnReconnect: true,
    refetchOnWindowFocus: false,
  });

export const useTwitchAuth = () =>
  useQuery(
    ['twitchAuth'],
    async () => {
      const accessToken = accessTokenStore.get();

      if (accessToken) {
        try {
          await getProfile();
          return { accessToken };
          // eslint-disable-next-line no-empty
        } catch (e) {}
      }

      return await handleTwitchLoginCallback();
    },
    {
      onSuccess: (data) => {
        accessTokenStore.set(data.accessToken);
        redirectToDashboard();
      },
      retry: false,
      refetchOnReconnect: true,
      refetchOnWindowFocus: false,
    },
  );
