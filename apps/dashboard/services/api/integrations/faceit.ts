import { useSWRConfig } from 'swr';
import useSWR from 'swr';

import { mutationOptions, swrAuthFetcher } from '../fetchWrappers';

import { useSelectedDashboard } from '@/services/dashboard';


export const useFaceitIntegration = () => {
  const [selectedDashboard] = useSelectedDashboard();
  const { mutate } = useSWRConfig();

  return {
    getProfile() {
      return useSWR(
        selectedDashboard
          ? `/api/v1/channels/${selectedDashboard.channelId}/integrations/faceit`
          : null,
        swrAuthFetcher,
      );
    },
    async getAuthLink(): Promise<string | undefined> {
      if (!selectedDashboard) {
        throw new Error('Cannot get link because dashboard not selected');
      }

      try {
        const data = await swrAuthFetcher(`/api/v1/channels/${selectedDashboard.channelId}/integrations/faceit/auth`);
        return data;
      } catch {
        return;
      }
    },
    async postCode(code: string) {
      if (!selectedDashboard) {
        throw new Error('Cannot post code because dashboard not selected');
      }

      await swrAuthFetcher(`/api/v1/channels/${selectedDashboard.channelId}/integrations/faceit`, {
        body: JSON.stringify({ code }),
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      // mutate(
      //   selectedDashboard
      //     ? `/api/v1/channels/${selectedDashboard.channelId}/integrations/faceit`
      //     : null, // which cache keys are updated
      //   undefined, // update cache data to `undefined`
      //   { ...mutationOptions, revalidate: true }, // do not revalidate
      // );
    },
    async logout() {
      if (!selectedDashboard) {
        throw new Error('Cannot post code because dashboard not selected');
      }

      await swrAuthFetcher(
        `/api/v1/channels/${selectedDashboard.channelId}/integrations/faceit/logout`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
        },
      );

      mutate(
        selectedDashboard
          ? `/api/v1/channels/${selectedDashboard.channelId}/integrations/faceit`
          : null, // which cache keys are updated
        undefined, // update cache data to `undefined`
        { ...mutationOptions, revalidate: true }, // do not revalidate
      );
    },
  };
};
