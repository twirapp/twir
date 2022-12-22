import { useSWRConfig } from 'swr';
import useSWR from 'swr';

import { mutationOptions, authFetcher } from '../fetchWrappers';

import { useSelectedDashboard } from '@/services/dashboard';

export type StreamLabsProfile = {
  avatar: string;
  name: string;
};

export const useStreamLabsIntegration = () => {
  const [selectedDashboard] = useSelectedDashboard();
  const { mutate } = useSWRConfig();

  return {
    getIntegration() {
      return useSWR<StreamLabsProfile>(
        selectedDashboard
          ? `/api/v1/channels/${selectedDashboard.channelId}/integrations/streamlabs`
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
          `/api/v1/channels/${selectedDashboard.channelId}/integrations/streamlabs/auth`,
        );
        return data;
      } catch {
        return;
      }
    },
    async postCode(code: string) {
      if (!selectedDashboard) {
        throw new Error('Cannot post code because dashboard not slected');
      }

      await authFetcher(
        `/api/v1/channels/${selectedDashboard.channelId}/integrations/streamlabs/token`,
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
          ? `/api/v1/channels/${selectedDashboard.channelId}/integrations/streamlabs`
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
        `/api/v1/channels/${selectedDashboard.channelId}/integrations/streamlabs/logout`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
        },
      );

      mutate(
        selectedDashboard
          ? `/api/v1/channels/${selectedDashboard.channelId}/integrations/streamlabs`
          : null, // which cache keys are updated
        undefined, // update cache data to `undefined`
        { ...mutationOptions, revalidate: true }, // do not revalidate
      );
    },
  };
};
