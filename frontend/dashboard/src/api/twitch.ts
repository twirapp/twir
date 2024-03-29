import { useMutation, useQuery } from '@tanstack/vue-query';
import type { GetResponse as RewardsResponse } from '@twir/api/messages/rewards/rewards';
import type {
	TwitchGetUsersResponse,
	TwitchSearchChannelsResponse,
} from '@twir/api/messages/twitch/twitch';
import { ComputedRef, isRef, MaybeRef, Ref } from 'vue';

import { protectedApiClient, unprotectedApiClient } from '@/api/twirp.js';


type TwitchIn = Ref<string[]> | Ref<string> | ComputedRef<string> | ComputedRef<string[]> | string[]
export const useTwitchGetUsers = (opts: {
	ids?: TwitchIn,
	names?: TwitchIn
}) => useQuery({
	queryKey: ['twitch', 'search', 'users', opts.ids, opts.names],
	queryFn: async (): Promise<TwitchGetUsersResponse> => {
		let ids = isRef(opts?.ids)
			? Array.isArray(opts.ids.value) ? opts.ids.value : [opts.ids.value]
			: opts?.ids ?? [''];
		let names = isRef(opts?.names)
			? Array.isArray(opts.names.value) ? opts.names.value : [opts.names.value]
			: opts?.names ?? [''];

		names = names.filter(n => n !== '');
		ids = ids.filter(n => n !== '');

		if (ids.length === 0 && names.length === 0) {
			return {
				users: [],
			};
		}

		const call = await unprotectedApiClient.twitchGetUsers({
			ids,
			names,
		});

		return call.response;
	},
});

export const useTwitchSearchChannels = (query: string | Ref<string>) => useQuery({
	queryKey: ['twitch', 'search', 'channels', query],
	queryFn: async (): Promise<TwitchSearchChannelsResponse> => {
		const rawQuery = isRef(query) ? query.value : query;

		if (!rawQuery) return {
			channels: [],
		};

		const call = await unprotectedApiClient.twitchSearchChannels({
			query: rawQuery,
		});

		return call.response;
	},
});

export const useTwitchRewards = () => useQuery({
	queryKey: ['twitchRewards'],
	queryFn: async (): Promise<RewardsResponse> => {
		const call = await protectedApiClient.rewardsGet({});
		return call.response;
	},
});

export const useTwitchSearchCategories = (query: string | Ref<string>) => useQuery({
	queryKey: ['twitchSearchCategories', query || ''],
	queryFn: async () => {
		const input = isRef(query) ? query.value : query;
		if (!input) return { categories: [] };

		const call = await protectedApiClient.twitchSearchCategories({ query: input });
		return call.response;
	},
});

export const useTwitchGetCategories = (ids: MaybeRef<string[]> | ComputedRef<string[]>) => useQuery({
	queryKey: ['twitchGetCategories', ids || ''],
	queryFn: async () => {
		const input = isRef(ids) ? ids.value : ids;
		if (!input) return { categories: [] };

		const call = await protectedApiClient.twitchGetCategories({ ids: input });
		return call.response;
	},
});

export const twitchSetChannelInformationMutation = () => useMutation({
	mutationKey: ['twitchSetChannelInformation'],
	mutationFn: async (req: { categoryId: string, title: string }) => {
		await protectedApiClient.twitchSetChannelInformation(req);
	},
});
