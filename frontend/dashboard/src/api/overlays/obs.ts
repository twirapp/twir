import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query';
import type {
	GetResponse,
	PostRequest,
} from '@twir/api/messages/modules_obs_websocket/modules_obs_websocket';
import { Ref, unref } from 'vue';

import { protectedApiClient } from '@/api/twirp.js';

export const useObsOverlayManager = () => {
	const queryClient = useQueryClient();
	const queryKey = ['obsSettings'];

	return {
		getSettings: (refetchEnabled = true) => useQuery({
			queryKey,
			queryFn: async (): Promise<GetResponse | null> => {
				try {
					const call = await protectedApiClient.modulesOBSWebsocketGet({});
					return call.response;
				} catch {
					return null;
				}
			},
			refetchInterval: refetchEnabled ? 1000 : 0,
		}),
		updateSettings: () => useMutation({
			mutationKey: ['obsSettingsUpdate'],
			mutationFn: async (opts: PostRequest | Ref<PostRequest>) => {
				const data = unref(opts);
				await protectedApiClient.modulesOBSWebsocketUpdate(data);
			},
			onSuccess: async () => {
				await queryClient.invalidateQueries(queryKey);
			},
		}),
	};
};
