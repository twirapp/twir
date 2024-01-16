import { useQueryClient, useQuery, useMutation } from '@tanstack/vue-query';
import type {
	Settings,
	UpdateRequest,
	GetAllResponse,
} from '@twir/api/messages/overlays_chat/overlays_chat';
import { unref } from 'vue';
import type { MaybeRef } from 'vue';

import { protectedApiClient } from '@/api/twirp.js';

export const useChatOverlayManager = () => {
	const queryClient = useQueryClient();
	const queryKey = 'chatOverlay';

	return {
		useGet: (id: MaybeRef<string>) => useQuery({
			queryKey: [queryKey, id],
			queryFn: async (): Promise<Settings | null> => {
				try {
					const call = await protectedApiClient.overlayChatGet({
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
				const call = await protectedApiClient.overlayChatGetAll({});
				return call.response;
			},
		}),
		useCreate: () => useMutation({
			mutationKey: ['chatOverlayCreate'],
			mutationFn: async (opts: MaybeRef<Settings>) => {
				const data = unref(opts);
				const call = await protectedApiClient.overlayChatCreate(data);
				return call.response;
			},
			onSuccess: async () => {
				await queryClient.invalidateQueries([queryKey]);
			},
		}),
		useUpdate: () => useMutation({
			mutationKey: ['chatOverlayUpdate'],
			mutationFn: async (opts: MaybeRef<UpdateRequest>) => {
				const data = unref(opts);
				await protectedApiClient.overlayChatUpdate(data);
			},
			onSuccess: async (_, opts) => {
				const data = unref(opts);
				await queryClient.invalidateQueries([queryKey, data.id]);
			},
		}),
		useDelete: () => useMutation({
			mutationKey: ['chatOverlayDelete'],
			mutationFn: async (id: MaybeRef<string>) => {
				await protectedApiClient.overlayChatDelete({
					id: unref(id),
				});
			},
			onSuccess: async () => {
				await queryClient.invalidateQueries([queryKey]);
			},
		}),
	};
};
