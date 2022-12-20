import { Dashboard } from '@tsuwari/shared';
import useSWR, { useSWRConfig } from 'swr';

import { useSelectedDashboard } from '../dashboard/useSelectedDashboard';

import { swrAuthFetcher } from '@/services/api';

export const useDashboardAccess = () => {
  const [selectedDashboard] = useSelectedDashboard();
  const { mutate } = useSWRConfig();

  return {
    getAll() {
      return useSWR<Dashboard[]>(
        selectedDashboard
          ? `/api/v1/channels/${selectedDashboard.channelId}/settings/dashboard-access`
          : null,
        swrAuthFetcher,
      );
    },

    delete(dashboardId: string) {
      if (selectedDashboard === null) {
        throw new Error('Selected dashboard is null, unable to delete dashboard.');
      }

      return mutate<Dashboard[]>(
        `/api/v1/channels/${selectedDashboard.channelId}/settings/dashboard-access`,
        async (dashboards) => {
          await swrAuthFetcher(
            `/api/v1/channels/${selectedDashboard.channelId}/settings/dashboard-access/${dashboardId}`,
            {
              method: 'DELETE',
              headers: {
                'Content-Type': 'application/json',
              },
            },
          );

          return dashboards?.filter((c) => c.id != dashboardId);
        },
        {
          revalidate: false,
        },
      );
    },
    create(userName: string) {
      if (selectedDashboard === null) {
        throw new Error('Selected dashboard is null, unable to delete dashboard.');
      }

      return mutate<Dashboard[]>(
        `/api/v1/channels/${selectedDashboard.channelId}/settings/dashboard-access`,
        async (dashboards) => {
          const data = await swrAuthFetcher(
            `/api/v1/channels/${selectedDashboard.channelId}/settings/dashboard-access`,
            {
              body: JSON.stringify({ userName }),
              method: 'POST',
              headers: {
                'Content-Type': 'application/json',
              },
            },
          );

          dashboards?.push(data);

          return dashboards;
        },
        {
          revalidate: false,
        },
      );
    },
  };
};
