import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query';
import type {
	CreateBadgeRequest,
	UpdateBadgeRequest,
} from '@twir/api/messages/admin_badges/admin_badges';
import type { GetBadgesResponse } from '@twir/api/messages/badges_unprotected/badges_unprotected';

import { adminApiClient, unprotectedApiClient } from './twirp';

const BADGES_QUERY_KEY = 'public/badges';

export const useAdminBadges = () => {
	const queryClient = useQueryClient();

	const useBadgesUpload = () => useMutation({
		mutationFn: async (data: CreateBadgeRequest) => {
			return await adminApiClient.badgesCreate(data);
		},
		onSuccess() {
			queryClient.invalidateQueries([BADGES_QUERY_KEY]);
		},
	});

	const useBadgesUpdate = () => useMutation({
		mutationFn: async (data: UpdateBadgeRequest) => {
			return await adminApiClient.badgesUpdate(data);
		},
		onSuccess() {
			queryClient.invalidateQueries([BADGES_QUERY_KEY]);
		},
	});

	const useBadgesDelete = () => useMutation({
		mutationFn: async (id: string) => {
			return await adminApiClient.badgesDelete({ id });
		},
		onSuccess(_, id) {
			queryClient.setQueriesData<GetBadgesResponse>([BADGES_QUERY_KEY], (data) => {
				if (!data) return data;

				return {
					badges: data.badges.filter((badge) => badge.id !== id),
				};
			});
		},
	});

	const useBadgesUserAdd = () => useMutation({
		mutationFn: async (opts: { badgeId: string, userId: string }) => {
			const req = await adminApiClient.badgeAddUser({
				badgeId: opts.badgeId,
				userId: opts.userId,
			});
			return req.response;
		},
		onSuccess(_, opts) {
			queryClient.setQueriesData<GetBadgesResponse>([BADGES_QUERY_KEY], (data) => {
				if (!data) return data;
				return {
					badges: data.badges.map((badge) => {
						if (badge.id === opts.badgeId) {
							return {
								...badge,
								users: [...badge.users, opts.userId],
							};
						}
						return badge;
					}),
				};
			});
		},
	});

	const useBadgesUserRemove = () => useMutation({
		mutationFn: async (opts: { badgeId: string, userId: string }) => {
			await adminApiClient.badgeDeleteUser({
				badgeId: opts.badgeId,
				userId: opts.userId,
			});
		},
		onSuccess(_, opts) {
			queryClient.setQueriesData<GetBadgesResponse>([BADGES_QUERY_KEY], (data) => {
				if (!data) return data;
				return {
					badges: data.badges.map((badge) => ({
						...badge,
						users: badge.id === opts.badgeId
							? badge.users.filter((userId) => userId !== opts.userId)
							: badge.users,
					})),
				};
			});
		},
	});

	return {
		useBadgesUpload,
		useBadgesUpdate,
		useBadgesDelete,
		useBadgesUserAdd,
		useBadgesUserRemove,
	};
};

export const useBadges = () => useQuery({
	queryKey: [BADGES_QUERY_KEY],
	queryFn: async () => {
		const request = await unprotectedApiClient.getBadgesWithUsers({});
		return request.response;
	},
});
