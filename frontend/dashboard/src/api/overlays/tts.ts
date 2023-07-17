import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import type { PostRequest } from '@twir/grpc/generated/api/api/modules_tts';

import { protectedApiClient } from '@/api/twirp.js';


export const useTtsOverlayManager = () => {
  const queryClient = useQueryClient();
  const queryKey = ['ttsSettings'];

  return {
    getSettings: () => useQuery({
      queryKey,
      queryFn: async () => {
        const call = await protectedApiClient.modulesTTSGet({});
        return call.response;
      },
    }),
    updateSettings: () => useMutation({
      mutationKey: ['ttsUpdate'],
      mutationFn: async (opts: PostRequest) => {
        await protectedApiClient.modulesTTSUpdate(opts);
      },
      onSuccess: async () => {
        await queryClient.invalidateQueries(queryKey);
      },
    }),
  };
};
