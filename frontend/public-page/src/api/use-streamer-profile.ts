import { useQuery } from '@urql/vue';
import { defineStore, storeToRefs } from 'pinia';
import { computed, unref, watch } from 'vue';
import { useRouter } from 'vue-router';

import { graphql } from '@/gql';
import { routeNames } from '@/router';

export const useStreamerProfile = defineStore('streamerProfile', () => {
	const router = useRouter();
	const streamerName = computed(() => {
		const name = router.currentRoute.value.params.channelName;
		if (typeof name !== 'string') {
			return '';
		}

		return name;
	});

	const request = useQuery({
		query: graphql(`
			query StreamerTwitchProfile($userName: String!) {
				twitchGetUserByName(name: $userName) {
					id
					profileImageUrl
					login
					description
					displayName
					notFound
				}
			}
		`),
		variables: {
			get userName() {
				return unref(streamerName);
			},
		},
	});

	watch([streamerName, request.error, request.data], ([name, error, data]) => {
		if (!name || error || data?.twitchGetUserByName?.notFound) {
			router.push({ name: routeNames.notFound });
		}
	});

	return request;
});

export const useStreamerPublicSettings = defineStore('streamerPublicSettings', () => {
	const { data } = storeToRefs(useStreamerProfile());

	return useQuery({
		query: graphql(`
			query StreamerPublicSettings($userId: String!) {
				userPublicSettings(userId: $userId) {
					description
					socialLinks {
						title
						href
					}
				}
			}
		`),
		variables: {
			get userId() {
				return data.value?.twitchGetUserByName?.id || '';
			},
		},
	});
});

// export const useStreamerPublicSettings = () => {
// 	const { isLoading, data } = useStreamerProfile();
//
// 	return useQuery({
// 		queryKey: ['usePublicSettings'],
// 		queryFn: async (): Promise<Settings | undefined> => {
// 			if (!data.value) return;
// 			const call = await unprotectedClient.getPublicSettings({
// 				channelId: data.value!.id,
// 			});
// 			return call.response;
// 		},
// 		enabled: () => !isLoading.value,
// 	});
// };
