import { ChannelIntegration } from '@tsuwari/typeorm/entities/ChannelIntegration';
import { useSWRConfig } from 'swr';
import useSWR from 'swr';

import { swrAuthFetcher } from '../fetchWrappers';

import { useSelectedDashboard } from '@/services/dashboard';

export type Profile = {
  id: string;
  display_name: string;
  images?: Array<{ url: string }>;
};

export const useSpotifyIntegration = () => {
  const [selectedDashboard] = useSelectedDashboard();
  const { mutate } = useSWRConfig();

  return {
    getIntegration() {
      return useSWR<ChannelIntegration>(
        selectedDashboard
          ? `/api/v1/channels/${selectedDashboard.channelId}/integrations/spotify`
          : null,
        swrAuthFetcher,
      );
    },
    getProfile() {
      return useSWR<Profile>(
        selectedDashboard
          ? `/api/v1/channels/${selectedDashboard.channelId}/integrations/spotify/profile`
          : null,
        swrAuthFetcher,
      );
    },
    getAuthLink(): Promise<string> {
      if (!selectedDashboard) {
        throw new Error('Cannot get link because dashboard not selected');
      }

      return swrAuthFetcher(
        `/api/v1/channels/${selectedDashboard.channelId}/integrations/spotify/auth`,
      );
    },
    async postCode(code: string) {
      if (!selectedDashboard) {
        throw new Error('Cannot post code because dashboard not slected');
      }

      await swrAuthFetcher(
        `/api/v1/channels/${selectedDashboard.channelId}/integrations/spotify/token`,
        {
          body: JSON.stringify({ code }),
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
        },
      );

      mutate(
        selectedDashboard
          ? `/api/v1/channels/${selectedDashboard.channelId}/integrations/spotify/profile`
          : null, // which cache keys are updated
        undefined, // update cache data to `undefined`
        { revalidate: false }, // do not revalidate
      );
    },
    async logout() {
      if (!selectedDashboard) {
        throw new Error('Cannot post code because dashboard not slected');
      }

      await swrAuthFetcher(
        `/api/v1/channels/${selectedDashboard.channelId}/integrations/spotify/logout`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
        },
      );

      mutate(
        selectedDashboard
          ? `/api/v1/channels/${selectedDashboard.channelId}/integrations/spotify/profile`
          : null, // which cache keys are updated
        undefined, // update cache data to `undefined`
        { revalidate: false }, // do not revalidate
      );
    },
  };
};
