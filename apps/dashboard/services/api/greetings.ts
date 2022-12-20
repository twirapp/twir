import { ChannelGreeting } from '@tsuwari/typeorm/entities/ChannelGreeting';
import useSWR, { useSWRConfig } from 'swr';

import { useSelectedDashboard } from '../dashboard/useSelectedDashboard';

import { mutationOptions, swrAuthFetcher } from '@/services/api';

export type Greeting = ChannelGreeting & { userName: string };

export const useGreetingsManager = () => {
  const [selectedDashboard] = useSelectedDashboard();
  const { mutate } = useSWRConfig();

  return {
    getAll() {
      return useSWR<Greeting[]>(
        selectedDashboard ? `/api/v1/channels/${selectedDashboard.channelId}/greetings` : null,
        swrAuthFetcher,
      );
    },
    createOrUpdate(greeting: ChannelGreeting) {
      if (selectedDashboard === null) {
        throw new Error('Selected dashboard is null, unable to delete greeting.');
      }

      return mutate<ChannelGreeting[]>(
        `/api/v1/channels/${selectedDashboard.channelId}/greetings`,
        async (greetings) => {
          const data = await swrAuthFetcher(
            `/api/v1/channels/${selectedDashboard.channelId}/greetings/${greeting.id ?? ''}`,
            {
              body: JSON.stringify(greeting),
              method: greeting.id ? 'PUT' : 'POST',
              headers: {
                'Content-Type': 'application/json',
              },
            },
          );

          if (greeting.id) {
            const index = greetings!.findIndex((c) => c.id === data.id);
            greetings![index] = data;
          } else {
            greetings?.push(data);
          }

          return greetings;
        },
        mutationOptions,
      );
    },
    async patch(greetingId: string, greetingData: Partial<ChannelGreeting>) {
      if (selectedDashboard === null) {
        throw new Error('Selected dashboard is null, unable to delete keyword.');
      }

      return mutate<ChannelGreeting[]>(
        `/api/v1/channels/${selectedDashboard.channelId}/greetings`,
        async (greetings) => {
          const data = await swrAuthFetcher(
            `/api/v1/channels/${selectedDashboard.channelId}/greetings/${greetingId}`,
            {
              body: JSON.stringify(greetingData),
              method: 'PATCH',
              headers: {
                'Content-Type': 'application/json',
              },
            },
          );

          const index = greetings!.findIndex((c) => c.id === greetingId);
          greetings![index] = data;

          return greetings;
        },
        mutationOptions,
      );
    },
    async delete(greetingId: string) {
      if (selectedDashboard === null) {
        throw new Error('Selected dashboard is null, unable to delete greeting.');
      }

      return mutate<ChannelGreeting[]>(
        `/api/v1/channels/${selectedDashboard.channelId}/greetings`,
        async (greetings) => {
          await swrAuthFetcher(
            `/api/v1/channels/${selectedDashboard.channelId}/greetings/${greetingId}`,
            {
              method: 'DELETE',
              headers: {
                'Content-Type': 'application/json',
              },
            },
          );

          return greetings?.filter((c) => c.id != greetingId);
        },
        mutationOptions,
      );
    },
  };
};
