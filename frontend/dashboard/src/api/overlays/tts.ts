import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import type {
	PostRequest,
	GetUsersSettingsResponse,
	GetInfoResponse,
	GetResponse,
} from '@twir/grpc/generated/api/api/modules_tts';
import { Ref, unref } from 'vue';

import { protectedApiClient } from '@/api/twirp.js';


export const useTtsOverlayManager = () => {
  const queryClient = useQueryClient();
  const queryKey = ['ttsSettings'];
  const usersQueryKey = ['ttsUsersSettings'];

  return {
    getSettings: () => useQuery({
      queryKey,
      queryFn: async (): Promise<GetResponse> => {
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
    getInfo: () => useQuery({
      queryKey: ['ttsInfo'],
      queryFn: async (): Promise<GetInfoResponse> => {
        const call = await protectedApiClient.modulesTTSGetInfo({});
        return call.response;
      },
    }),
    getUsersSettings: () => useQuery({
      queryKey: usersQueryKey,
      queryFn: async (): Promise<GetUsersSettingsResponse> => {
        const call = await protectedApiClient.modulesTTSGetUsersSettings({});
        return call.response;
      },
    }),
    deleteUsersSettings: () => useMutation({
      mutationKey: ['ttsUsersSettingsDelete'],
      mutationFn: async (ids: string[] | Ref<string[]>) => {
        const usersIds = unref(ids);

        await protectedApiClient.modulesTTSUsersDelete({ usersIds });
      },
      onSuccess: async () => {
        await queryClient.invalidateQueries(usersQueryKey);
      },
    }),
  };
};
