import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query';
import type { Settings } from '@twir/grpc/generated/api/api/overlays_kappagen';

import { protectedApiClient } from '../twirp';

export const useKappaGenOverlayManager = () => {
  const queryClient = useQueryClient();
  const queryKey = ['kappagenSettings'];

  return {
    getSettings: () => useQuery({
      queryKey,
      queryFn: async (): Promise<Settings> => {
        const call = await protectedApiClient.overlayKappaGenGet({});
        return call.response;
      },
    }),
    updateSettings: () => useMutation({
      mutationKey: ['ttsUpdate'],
      mutationFn: async (opts: Settings) => {
        await protectedApiClient.overlayKappaGenUpdate(opts);
      },
      onSuccess: async () => {
        await queryClient.invalidateQueries(queryKey);
      },
    }),
  };
};
