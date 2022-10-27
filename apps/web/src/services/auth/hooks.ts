import { useQuery } from '@tanstack/vue-query';

import { getProfile } from './api.js';
import { redirectToDashboard } from './locationHelpers.js';
import { accessTokenStore } from './token.js';
import { handleTwitchLoginCallback } from './twitch.js';

export const useUserProfile = () =>
  useQuery(['profile'], getProfile, {
    retry: false,
    refetchOnReconnect: 'always',
  });

export const useTwitchAuth = () =>
  useQuery(['twitchAuth'], handleTwitchLoginCallback, {
    onSuccess: (data) => {
      accessTokenStore.set(data.accessToken);
      redirectToDashboard();
    },
  });
