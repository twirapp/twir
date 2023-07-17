import { useQuery } from '@tanstack/vue-query';
import { GetSearchRequest_Type } from '@twir/grpc/generated/api/api/modules_sr';
import { type MaybeRef, isRef, toRaw } from 'vue';

import { protectedApiClient } from '@/api/twirp.js';

export const enum YoutubeSearchType  {
	Channel = 'channel',
	Video = 'video',
}

const searchType = {
	[YoutubeSearchType.Channel]: GetSearchRequest_Type.CHANNEL,
	[YoutubeSearchType.Video]: GetSearchRequest_Type.VIDEO,
};

export const useYoutubeVideoOrChannelSearch = (
	query: MaybeRef<string | string[]>,
	type: YoutubeSearchType,
) => {
	return useQuery({
		queryKey: [query, type],
		queryFn: async () => {
			const q = isRef(query) ? toRaw(query.value) : query;
			const qArray = (Array.isArray(q) ? q : [q]).filter(Boolean);
			if (!qArray.length) {
				return {
					items: [],
				};
			}

			const call = await protectedApiClient.modulesSRSearchVideosOrChannels({
				type: searchType[type],
				query: qArray,
			});

			return call.response;
		},
	});
};
