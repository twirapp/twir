import { ChannelIntegration } from '@tsuwari/typeorm/entities/ChannelIntegration';
import { useSWRConfig } from 'swr';
import useSWR from 'swr';

import { mutationOptions, authFetcher } from '../fetchWrappers';

import { useSelectedDashboard } from '@/services/dashboard';

export type SpotifyProfile = {
  id: string;
  display_name: string;
  images?: Array<{ url: string }>;
};

export const useSpotifyIntegration = () => {
  const [selectedDashboard] = useSelectedDashboard();
  const { mutate } = useSWRConfig();

  return {
    getProfile() {
      return useSWR<SpotifyProfile>(
        selectedDashboard
          ? `/api/v1/channels/${selectedDashboard.channelId}/integrations/spotify/profile`
          : null,
        authFetcher,
      );
    },
    async getAuthLink(): Promise<string | undefined> {
      if (!selectedDashboard) {
        throw new Error('Cannot get link because dashboard not selected');
      }

      try {
        const data = await authFetcher(
          `/api/v1/channels/${selectedDashboard.channelId}/integrations/spotify/auth`,
        );
        return data;
      } catch {
        return;
      }
    },
    async postCode(code: string) {
      if (!selectedDashboard) {
        throw new Error('Cannot post code because dashboard not selected');
      }

      await authFetcher(
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
        { ...mutationOptions, revalidate: true }, // do not revalidate
      );
    },
    async logout() {
      if (!selectedDashboard) {
        throw new Error('Cannot post code because dashboard not selected');
      }

      await authFetcher(
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
        { ...mutationOptions, revalidate: true }, // do not revalidate
      );
    },
  };
};
