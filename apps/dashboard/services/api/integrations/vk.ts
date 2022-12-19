import { useSWRConfig } from 'swr';
import useSWR from 'swr';

import { swrAuthFetcher } from '../fetchWrappers';

import { useSelectedDashboard } from '@/services/dashboard';

export type vkprofile = {
  name: string;
  image: string;
  playCount: string;
};

export const useVkIntegration = () => {
  const [selectedDashboard] = useSelectedDashboard();
  const { mutate } = useSWRConfig();

  return {
    getProfile() {
      return useSWR<vkprofile>(
        selectedDashboard
          ? `/api/v1/channels/${selectedDashboard.channelId}/integrations/vk`
          : null,
        swrAuthFetcher,
      );
    },
    getAuthLink(): Promise<string> {
      if (!selectedDashboard) {
        throw new Error('Cannot get link because dashboard not selected');
      }

      return swrAuthFetcher(`/api/v1/channels/${selectedDashboard.channelId}/integrations/vk/auth`);
    },
    async postToken(token: string) {
      if (!selectedDashboard) {
        throw new Error('Cannot post code because dashboard not slected');
      }

      await swrAuthFetcher(`/api/v1/channels/${selectedDashboard.channelId}/integrations/vk`, {
        body: JSON.stringify({ token }),
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
        { revalidate: true }, // do not revalidate
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
        { revalidate: true }, // do not revalidate
      );
    },
  };
};
