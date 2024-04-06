import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query';
import type { GetBadgesResponse } from '@twir/api/messages/badges_unprotected/badges_unprotected';
import { defineStore } from 'pinia';
import { computed } from 'vue';

import { adminApiClient, unprotectedApiClient } from '@/api/twirp';

export const useBadges = defineStore('admin-panel/badges', () => {
	const { data } = useQuery({
		queryKey: ['admin/badges'],
		queryFn: async () => {
			const request = await unprotectedApiClient.getBadgesWithUsers({});
			return request.response;
		},
	});

	const queryClient = useQueryClient();
	const deleter = useMutation({
		mutationFn: async (id: string) => {
			return await adminApiClient.badgesDelete({ id });
		},
		onSuccess(_, id) {
			queryClient.setQueriesData<GetBadgesResponse>(['admin/badges'], (data) => {
				if (!data) return data;

				return {
					badges: badges.value.filter(badge => badge.id != id),
				};
			});
		},
	});

	async function deleteBadge(id: string) {
		await deleter.mutateAsync(id);
	}

	const badges = computed(() => data.value?.badges ?? []);

	return {
		badges,
		deleteBadge,
	};
});
