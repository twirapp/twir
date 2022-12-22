import { useSWRConfig } from 'swr';
import useSWR from 'swr';

import { mutationOptions, authFetcher } from '@/services/api';
import { useSelectedDashboard } from '@/services/dashboard';

export type LastfmProfile = {
  name: string;
  image: string;
  playCount: string;
};

export const useLastfmIntegration = () => {
  const [selectedDashboard] = useSelectedDashboard();
  const { mutate } = useSWRConfig();

  return {
    getProfile() {
      return useSWR<LastfmProfile>(
        selectedDashboard
          ? `/api/v1/channels/${selectedDashboard.channelId}/integrations/lastfm/profile`
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
          `/api/v1/channels/${selectedDashboard.channelId}/integrations/lastfm/auth`,
        );
        return data;
      } catch {
        return;
      }
    },
    async postToken(token: string) {
      if (!selectedDashboard) {
        throw new Error('Cannot post code because dashboard not slected');
      }

      await authFetcher(`/api/v1/channels/${selectedDashboard.channelId}/integrations/lastfm`, {
        body: JSON.stringify({ token }),
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      mutate(
        selectedDashboard
          ? `/api/v1/channels/${selectedDashboard.channelId}/integrations/lastfm/profile`
          : null, // which cache keys are updated
        undefined, // update cache data to `undefined`
        { ...mutationOptions, revalidate: true }, // do not revalidate
      );
    },
    async logout() {
      if (!selectedDashboard) {
        throw new Error('Cannot post code because dashboard not slected');
      }

      await authFetcher(
        `/api/v1/channels/${selectedDashboard.channelId}/integrations/lastfm/logout`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
        },
      );

      mutate(
        selectedDashboard
          ? `/api/v1/channels/${selectedDashboard.channelId}/integrations/lastfm/profile`
          : null, // which cache keys are updated
        undefined, // update cache data to `undefined`
        { ...mutationOptions, revalidate: true }, // do not revalidate
      );
    },
  };
};
