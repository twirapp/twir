import { useQueryClient, useQuery, useMutation } from '@tanstack/vue-query';
import type {
  Settings,
} from '@twir/grpc/generated/api/api/modules_chat_overlay';
import { Ref, unref } from 'vue';

import { protectedApiClient } from '@/api/twirp.js';

export const useChatOverlayManager = () => {
  const queryClient = useQueryClient();
  const queryKey = ['chatOverlay'];

  return {
    getSettings: () => useQuery({
      queryKey,
      queryFn: async (): Promise<Settings | null> => {
				try {
					const call = await protectedApiClient.modulesChatOverlayGet({});
					return call.response;
				} catch {
					return null;
				}
      },
    }),
    updateSettings: () => useMutation({
      mutationKey: ['chatOverlayUpdate'],
      mutationFn: async (opts: Settings | Ref<Settings>) => {
        const data = unref(opts);
        await protectedApiClient.modulesChatOverlayUpdate(data);
      },
      onSuccess: async () => {
        await queryClient.invalidateQueries(queryKey);
      },
    }),
  };
};
