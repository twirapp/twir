import { useQueryClient, useQuery, useMutation } from '@tanstack/vue-query';
import type {
	Settings,
	UpdateRequest,
	GetAllResponse,
} from '@twir/api/messages/overlays_now_playing/overlays_now_playing';
import { unref } from 'vue';
import type { MaybeRef } from 'vue';

import { protectedApiClient } from '@/api/twirp.js';

export const useNowPlayingOverlayManager = () => {
	const queryClient = useQueryClient();
	const queryKey = 'nowPlayingOverlay';

	return {
		useGetAll: () => useQuery({
			queryKey: [queryKey],
			queryFn: async (): Promise<GetAllResponse> => {
				const call = await protectedApiClient.overlaysNowPlayingGetAll({});
				return call.response;
			},
		}),
		useCreate: () => useMutation({
			mutationKey: ['nowPlayingOverlayCreate'],
			mutationFn: async (opts: MaybeRef<Settings>) => {
				const data = unref(opts);
				const call = await protectedApiClient.overlaysNowPlayingCreate(data);
				return call.response;
			},
			onSuccess: async () => {
				await queryClient.invalidateQueries([queryKey]);
			},
		}),
		useUpdate: () => useMutation({
			mutationKey: ['nowPlayingOverlayUpdate'],
			mutationFn: async (opts: MaybeRef<UpdateRequest>) => {
				const data = unref(opts);
				await protectedApiClient.overlaysNowPlayingUpdate(data);
			},
			onSuccess: async (_, opts) => {
				const data = unref(opts);
				await queryClient.invalidateQueries([queryKey, data.id]);
			},
		}),
		useDelete: () => useMutation({
			mutationKey: ['nowPlayingOverlayDelete'],
			mutationFn: async (id: MaybeRef<string>) => {
				await protectedApiClient.overlaysNowPlayingDelete({
					id: unref(id),
				});
			},
			onSuccess: async () => {
				await queryClient.invalidateQueries([queryKey]);
			},
		}),
	};
};
