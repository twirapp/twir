import { useMutation, useQuery } from '@tanstack/react-query';
import {
	GetUsersRequest_Order,
	GetUsersRequest_SortBy,
	ResetStatsRequest_Field,
} from '@twir/grpc/generated/api/api/community';

import { protectedApiClient } from '@/services/api/twirp';

export type SortByField = 'messages' | 'watched' | 'emotes' | 'usedChannelPoints'

export const useCommunityUsers = (
	limit= 50,
	page = 1,
	sortBy: SortByField,
	order: 'asc' | 'desc' = 'desc',
) => useQuery<ReturnType<typeof protectedApiClient.communityGetUsers>['response']>({
	queryKey: ['communityUsers'],
	queryFn: async () => {
		let sort: GetUsersRequest_SortBy;
		switch (sortBy) {
			case 'emotes':
				sort = GetUsersRequest_SortBy.Emotes;
				break;
			case 'messages':
				sort = GetUsersRequest_SortBy.Messages;
				break;
			case 'usedChannelPoints':
				sort = GetUsersRequest_SortBy.UsedChannelPoints;
				break;
			case 'watched':
				sort = GetUsersRequest_SortBy.Watched;
				break;
			default:
				sort = GetUsersRequest_SortBy.Watched;
		}

		const call = await protectedApiClient.communityGetUsers({
			limit,
			page,
			sortBy: sort,
			order: order === 'asc' ? GetUsersRequest_Order.Asc : GetUsersRequest_Order.Desc,
		});

		return call.response;
	},
	retry: false,
	refetchInterval: 1000,
});

export const useCommunityReset = useMutation({
	mutationFn: async (field: SortByField) => {
		let resetField: ResetStatsRequest_Field;
		switch (field) {
			case 'emotes':
				resetField = ResetStatsRequest_Field.Emotes;
				break;
			case 'messages':
				resetField = ResetStatsRequest_Field.Messages;
				break;
			case 'usedChannelPoints':
				resetField = ResetStatsRequest_Field.UsedChannelsPoints;
				break;
			case 'watched':
				resetField = ResetStatsRequest_Field.Watched;
		}

		await protectedApiClient.communityResetStats({
			field: resetField,
		});
	},
});
