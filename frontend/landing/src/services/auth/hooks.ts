import { useQuery } from '@tanstack/vue-query';

import { getProfile } from './api.js';

import { redirectToDashboard } from '@/services/auth/locationHelpers.js';
import { handleTwitchLoginCallback } from '@/services/auth/twitch.js';

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
			await handleTwitchLoginCallback();
			redirectToDashboard();
		},
		{
			retry: false,
			refetchOnReconnect: true,
			refetchOnWindowFocus: false,
		},
	);
