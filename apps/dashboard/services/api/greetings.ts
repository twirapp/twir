import { ChannelGreeting } from '@tsuwari/typeorm/entities/ChannelGreeting';
import useSWR, { useSWRConfig } from 'swr';

import { useSelectedDashboard } from '../dashboard/useSelectedDashboard';

import { swrAuthFetcher } from '@/services/api';

export type Greeting = ChannelGreeting & { userName: string };

export const useGreetings = () => {
  const [selectedDashboard] = useSelectedDashboard();

  return useSWR<Greeting[]>(
    selectedDashboard ? `/api/v1/channels/${selectedDashboard.channelId}/greetings` : null,
    swrAuthFetcher,
  );
};

export const useManageGreeting = () => {
  const { mutate } = useSWRConfig();
  const [selectedDashboard] = useSelectedDashboard();

  if (selectedDashboard === null) {
    throw new Error('Selected dashboard is null, unable to post command.');
  }

  return (greeting: ChannelGreeting) =>
    mutate<ChannelGreeting[]>(
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
      {
        revalidate: false,
      },
    );
};
