import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import {
	GetUsersRequest_Order,
	GetUsersRequest_SortBy,
	ResetStatsRequest_Field,
} from '@twir/grpc/generated/api/api/community';
import { Ref, isRef } from 'vue';

import { protectedApiClient, unprotectedApiClient } from '@/api/twirp.js';

export const enum ComminityOrder {
	Desc = 'desc',
	Asc = 'asc',
}

export const enum CommunitySortBy {
	Watched = 'watched',
	Messages = 'messages',
	Emotes = 'emotes',
	UsedChannelPoints = 'used_channel_points',
}

export type GetCommunityUsersOpts = {
	limit: number;
	page: number;
	order: ComminityOrder;
	sortBy: CommunitySortBy;
	channelId?: string
}

const sortBy = {
	[CommunitySortBy.Watched]: GetUsersRequest_SortBy.Watched,
	[CommunitySortBy.Messages]: GetUsersRequest_SortBy.Messages,
	[CommunitySortBy.Emotes]: GetUsersRequest_SortBy.Emotes,
	[CommunitySortBy.UsedChannelPoints]: GetUsersRequest_SortBy.UsedChannelPoints,
};

export const useCommunityUsers = () => {
	return {
		getAll: (opts: GetCommunityUsersOpts | Ref<GetCommunityUsersOpts>) => useQuery({
			queryKey: ['communityUsers', opts],
			queryFn: async () => {
				const rawOpts = isRef(opts) ? opts.value : opts;

				if (!rawOpts.channelId) return;

				const order = rawOpts.order === ComminityOrder.Desc
					? GetUsersRequest_Order.Desc
					: GetUsersRequest_Order.Asc;

				const call = await unprotectedApiClient.communityGetUsers({
					limit: rawOpts.limit,
					page: rawOpts.page,
					order,
					sortBy: sortBy[rawOpts.sortBy],
					channelId: rawOpts.channelId,
				});
				return call.response;
			},
		}),
	};
};

export const enum CommunityResetStatsField {
	Watched = 'watched',
	Messages = 'messages',
	Emotes = 'emotes',
	UsedChannelPoints = 'used_channel_points',
}

const resetFields = {
	[CommunityResetStatsField.Watched]: ResetStatsRequest_Field.Watched,
	[CommunityResetStatsField.Messages]: ResetStatsRequest_Field.Messages,
	[CommunityResetStatsField.Emotes]: ResetStatsRequest_Field.Emotes,
	[CommunityResetStatsField.UsedChannelPoints]: ResetStatsRequest_Field.UsedChannelsPoints,
};

export const useCommunityReset = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: async (field: CommunityResetStatsField) => {
			const call = await protectedApiClient.communityResetStats({
				field: resetFields[field],
			});
			return call.response;
		},
		onSuccess: () => {
			queryClient.invalidateQueries(['communityUsers']);
		},
	});
};
