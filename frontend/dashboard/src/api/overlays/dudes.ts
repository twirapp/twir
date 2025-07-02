import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query';
import type {
	GetAllResponse,
	Settings,
	UpdateRequest,
} from '@twir/api/messages/overlays_dudes/overlays_dudes';
import { unref } from 'vue';
import type { MaybeRef } from 'vue';

import { protectedApiClient } from '@/api/twirp.js';

export const useDudesOverlayManager = () => {
	const queryClient = useQueryClient();
	const queryKey = 'dudesOverlay';

	return {
		useGet: (id: MaybeRef<string>) => useQuery({
			queryKey: [queryKey, id],
			queryFn: async (): Promise<Settings | null> => {
				try {
					const call = await protectedApiClient.overlayDudesGet({
						id: unref(id),
					});
					return call.response;
				} catch {
					return null;
				}
			},
		}),
		useGetAll: () => useQuery({
			queryKey: [queryKey],
			queryFn: async (): Promise<GetAllResponse> => {
				const call = await protectedApiClient.overlayDudesGetAll({});
				return call.response;
			},
		}),
		useCreate: () => useMutation({
			mutationKey: ['dudesOverlayCreate'],
			mutationFn: async (opts: MaybeRef<Settings>) => {
				const data = unref(opts);
				const call = await protectedApiClient.overlayDudesCreate(data);
				return call.response;
			},
			onSuccess: async () => {
				await queryClient.invalidateQueries([queryKey]);
			},
		}),
		useUpdate: () => useMutation({
			mutationKey: ['dudesOverlayUpdate'],
			mutationFn: async (opts: MaybeRef<UpdateRequest>) => {
				const data = unref(opts);
				await protectedApiClient.overlayDudesUpdate(data);
			},
			onSuccess: async (_, opts) => {
				const data = unref(opts);
				await queryClient.invalidateQueries([queryKey, data.id]);
			},
		}),
		useDelete: () => useMutation({
			mutationKey: ['dudesOverlayDelete'],
			mutationFn: async (id: MaybeRef<string>) => {
				await protectedApiClient.overlayDudesDelete({
					id: unref(id),
				});
			},
			onSuccess: async () => {
				await queryClient.invalidateQueries([queryKey]);
			},
		}),
	};
};
