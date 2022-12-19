import { ChannelGreeting } from '@tsuwari/typeorm/entities/ChannelGreeting';
import useSWR from 'swr';

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
