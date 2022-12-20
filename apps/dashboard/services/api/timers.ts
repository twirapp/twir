import { ChannelTimer } from '@tsuwari/typeorm/entities/ChannelTimer';
import useSWR, { useSWRConfig } from 'swr';

import { mutationOptions, swrAuthFetcher } from '@/services/api';
import { useSelectedDashboard } from '@/services/dashboard';

export const useTimersManager = () => {
  const { mutate } = useSWRConfig();
  const [selectedDashboard] = useSelectedDashboard();

  return {
    getAll() {
      return useSWR<ChannelTimer[]>(
        selectedDashboard ? `/api/v1/channels/${selectedDashboard.channelId}/timers` : null,
        swrAuthFetcher,
      );
    },
    async createOrUpdate(timer: ChannelTimer) {
      if (selectedDashboard === null) {
        throw new Error('Selected dashboard is null, unable to delete timer.');
      }

      return mutate<ChannelTimer[]>(
        `/api/v1/channels/${selectedDashboard.channelId}/timers`,
        async (timers) => {
          const data = await swrAuthFetcher(
            `/api/v1/channels/${selectedDashboard.channelId}/timers/${timer.id ?? ''}`,
            {
              body: JSON.stringify(timer),
              method: timer.id ? 'PUT' : 'POST',
              headers: {
                'Content-Type': 'application/json',
              },
            },
          );

          if (timer.id) {
            const index = timers!.findIndex((t) => t.id === data.id);
            timers![index] = data;
          } else {
            timers?.push(data);
          }

          return timers;
        },
        mutationOptions,
      );
    },
    async delete(timerID: string) {
      if (selectedDashboard === null) {
        throw new Error('Selected dashboard is null, unable to delete timer.');
      }

      return mutate<ChannelTimer[]>(
        `/api/v1/channels/${selectedDashboard.channelId}/timers`,
        async (timers) => {
          await swrAuthFetcher(
            `/api/v1/channels/${selectedDashboard.channelId}/timers/${timerID}`,
            {
              method: 'DELETE',
              headers: {
                'Content-Type': 'application/json',
              },
            },
          );

          return timers?.filter((c) => c.id != timerID);
        },
        mutationOptions,
      );
    },
    async patch(timerId: string, timerData: Partial<ChannelTimer>) {
      if (selectedDashboard === null) {
        throw new Error('Selected dashboard is null, unable to delete timer.');
      }

      return mutate<ChannelTimer[]>(
        `/api/v1/channels/${selectedDashboard.channelId}/timers`,
        async (timers) => {
          const data = await swrAuthFetcher(
            `/api/v1/channels/${selectedDashboard.channelId}/timers/${timerId}`,
            {
              body: JSON.stringify(timerData),
              method: 'PATCH',
              headers: {
                'Content-Type': 'application/json',
              },
            },
          );

          const index = timers!.findIndex((c) => c.id === timerId);
          timers![index] = data;

          return timers;
        },
        mutationOptions,
      );
    },
  };
};
