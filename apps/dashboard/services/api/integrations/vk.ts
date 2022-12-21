import { useSWRConfig } from 'swr';
import useSWR from 'swr';

import { mutationOptions, swrAuthFetcher } from '../fetchWrappers';

import { useSelectedDashboard } from '@/services/dashboard';

export type VKProfile = {
  id: number;
  first_name: string;
  last_name: string;
  photo_max_orig: string;
};

export const useVkIntegration = () => {
  const [selectedDashboard] = useSelectedDashboard();
  const { mutate } = useSWRConfig();

  return {
    getProfile() {
      return useSWR<VKProfile>(
        selectedDashboard
          ? `/api/v1/channels/${selectedDashboard.channelId}/integrations/vk`
          : null,
        swrAuthFetcher,
      );
    },
    async getAuthLink(): Promise<string | undefined> {
      if (!selectedDashboard) {
        throw new Error('Cannot get link because dashboard not selected');
      }

      try {
        const data = await swrAuthFetcher(`/api/v1/channels/${selectedDashboard.channelId}/integrations/vk/auth`);
        return data
      } catch {
        return
      }
    },
    async postCode(code: string) {
      if (!selectedDashboard) {
        throw new Error('Cannot post code because dashboard not slected');
      }

      await swrAuthFetcher(`/api/v1/channels/${selectedDashboard.channelId}/integrations/vk`, {
        body: JSON.stringify({ code }),
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      mutate(
        selectedDashboard
          ? `/api/v1/channels/${selectedDashboard.channelId}/integrations/vk`
          : null, // which cache keys are updated
        undefined, // update cache data to `undefined`
        { ...mutationOptions, revalidate: true }, // do not revalidate
      );
    },
    async logout() {
      if (!selectedDashboard) {
        throw new Error('Cannot post code because dashboard not slected');
      }

      await swrAuthFetcher(
        `/api/v1/channels/${selectedDashboard.channelId}/integrations/vk/logout`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
        },
      );

      mutate(
        selectedDashboard
          ? `/api/v1/channels/${selectedDashboard.channelId}/integrations/vk`
          : null, // which cache keys are updated
        undefined, // update cache data to `undefined`
        { ...mutationOptions, revalidate: true }, // do not revalidate
      );
    },
  };
};
