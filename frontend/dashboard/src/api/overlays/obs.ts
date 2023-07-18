import { useQueryClient, useQuery, useMutation } from '@tanstack/vue-query';
import type {
  GetResponse,
  PostRequest,
} from '@twir/grpc/generated/api/api/modules_obs_websocket';
import { Ref, toRaw, unref } from 'vue';

import { protectedApiClient } from '@/api/twirp.js';

export const useObsOverlayManager = () => {
  const queryClient = useQueryClient();
  const queryKey = ['obsSettings'];

  return {
    getSettings: () => useQuery({
      queryKey,
      queryFn: async (): Promise<GetResponse> => {
        const call = await protectedApiClient.modulesOBSWebsocketGet({});
        return call.response;
      },
    }),
    updateSettings: () => useMutation({
      mutationKey: ['obsSettingsUpdate'],
      mutationFn: async (opts: PostRequest | Ref<PostRequest>) => {
        const data = unref(opts);
        await protectedApiClient.modulesOBSWebsocketUpdate(data);
      },
    }),
  };
};
