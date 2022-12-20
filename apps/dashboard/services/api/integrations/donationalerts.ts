import { ChannelIntegration } from '@tsuwari/typeorm/entities/ChannelIntegration';
import { useSWRConfig } from 'swr';
import useSWR from 'swr';

import { mutationOptions, swrAuthFetcher } from '../fetchWrappers';

import { useSelectedDashboard } from '@/services/dashboard';

export type DAPRofile = {
  avatar: string;
  name: string;
};

export const useDonationAlertsIntegration = () => {
  const [selectedDashboard] = useSelectedDashboard();
  const { mutate } = useSWRConfig();

  return {
    getIntegration() {
      return useSWR<DAPRofile>(
        selectedDashboard
          ? `/api/v1/channels/${selectedDashboard.channelId}/integrations/donationalerts`
          : null,
        swrAuthFetcher,
      );
    },
    getAuthLink(): Promise<string> {
      if (!selectedDashboard) {
        throw new Error('Cannot get link because dashboard not selected');
      }

      return swrAuthFetcher(
        `/api/v1/channels/${selectedDashboard.channelId}/integrations/donationalerts/auth`,
      );
    },
    async postCode(code: string) {
      if (!selectedDashboard) {
        throw new Error('Cannot post code because dashboard not slected');
      }

      await swrAuthFetcher(
        `/api/v1/channels/${selectedDashboard.channelId}/integrations/donationalerts/token`,
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
          ? `/api/v1/channels/${selectedDashboard.channelId}/integrations/donationalerts`
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
        `/api/v1/channels/${selectedDashboard.channelId}/integrations/donationalerts/logout`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
        },
      );

      mutate(
        selectedDashboard
          ? `/api/v1/channels/${selectedDashboard.channelId}/integrations/donationalerts`
          : null, // which cache keys are updated
        undefined, // update cache data to `undefined`
        { ...mutationOptions, revalidate: true }, // do not revalidate
      );
    },
  };
};
