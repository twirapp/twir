import type { GetPublicUserInfoResponse } from '@twir/api/messages/auth/auth';
import type { TwitchUser } from '@twir/api/messages/twitch/twitch';
import { defineStore } from 'pinia';
import { ref } from 'vue';

import { unprotectedClient } from '@/api/twirp';

export type StreamerProfile = TwitchUser & GetPublicUserInfoResponse;

export const useStreamerProfile = defineStore('streamerProfile', () => {
	const profile = ref<StreamerProfile | null>();
	const isLoading = ref(true);

	async function fetchProfile(streamerName: string) {
		if (profile.value) return;

		const twitchUserCall = await unprotectedClient.twitchGetUsers({
			names: [streamerName],
			ids: [],
		});

		const twitchUser = twitchUserCall.response.users[0];
		if (!twitchUser) return;

		const internalUserCall = await unprotectedClient.getPublicUserInfo({
			userId: twitchUser.id,
		});

		const internalUser = internalUserCall.response;
		if (!internalUser) return;

		profile.value = {
			...twitchUser,
			...internalUser,
		};
		isLoading.value = false;
	}

	return {
		profile,
		fetchProfile,
		isLoading,
	};
});
