import { useQuery } from '@tanstack/vue-query';
import { GetUsersRequest_Order, GetUsersRequest_SortBy } from '@twir/grpc/generated/api/api/community';
import { Ref, isRef } from 'vue';

import { protectedApiClient } from '@/api/twirp.js';

export const enum UsersOrder {
	Desc = 'desc',
	Asc = 'asc',
}

export const enum UsersSortBy {
	Watched = 'watched',
	Messages = 'messages',
	Emotes = 'emotes',
	UsedChannelPoints = 'used_channel_points',
}

export type GetCommunityUsersOpts = {
	limit: number;
	page: number;
	order: UsersOrder;
	sortBy: UsersSortBy;
}

export const useCommunityUsers = () => {
	return {
		getAll: (rawOpts: GetCommunityUsersOpts | Ref<GetCommunityUsersOpts>) => useQuery({
			queryKey: ['communityUsers', rawOpts],
			queryFn: async () => {
				const opts = isRef(rawOpts) ? rawOpts.value : rawOpts;

				const order = opts.order === UsersOrder.Desc
					? GetUsersRequest_Order.Desc
					: GetUsersRequest_Order.Asc;

				let sortBy = GetUsersRequest_SortBy.Watched;

				switch (opts.sortBy) {
					case UsersSortBy.Watched:
						sortBy = GetUsersRequest_SortBy.Watched;
						break;
					case UsersSortBy.Messages:
						sortBy = GetUsersRequest_SortBy.Messages;
						break;
					case UsersSortBy.Emotes:
						sortBy = GetUsersRequest_SortBy.Emotes;
						break;
					case UsersSortBy.UsedChannelPoints:
						sortBy = GetUsersRequest_SortBy.UsedChannelPoints;
						break;
				}

				const call = await protectedApiClient.communityGetUsers({
					limit: opts.limit,
					page: opts.page,
					order,
					sortBy,
				});
				return call.response;
			},
		}),
	};
};
