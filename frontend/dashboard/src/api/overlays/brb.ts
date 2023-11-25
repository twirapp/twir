import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query';
import { Settings } from '@twir/grpc/generated/api/api/overlays_be_right_back';
import { Ref, unref } from 'vue';

import { protectedApiClient } from '../twirp.js';

export const useBeRightBackOverlayManager = () => {
  const queryClient = useQueryClient();
  const queryKey = ['brbOverlay'];

  return {
    getSettings: () => useQuery({
      queryKey,
      queryFn: async (): Promise<Settings | null> => {
				try {
					const call = await protectedApiClient.overlayBeRightBackGet({});
					return call.response;
				} catch {
					return null;
				}
      },
    }),
    updateSettings: () => useMutation({
      mutationKey: ['brbOverlayUpdate'],
      mutationFn: async (opts: Settings | Ref<Settings>) => {
        const data = unref(opts);
        await protectedApiClient.overlayBeRightBackUpdate(data);
      },
      onSuccess: async () => {
        await queryClient.invalidateQueries(queryKey);
      },
    }),
  };
};
