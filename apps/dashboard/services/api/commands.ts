import { ChannelCommand } from '@tsuwari/typeorm/entities/ChannelCommand';
import useSWR, { useSWRConfig } from 'swr';

import { swrAuthFetcher } from '@/services/api';
import { useSelectedDashboard } from '@/services/dashboard';

export const useCommands = () => {
  const [selectedDashboard] = useSelectedDashboard();

  return useSWR<ChannelCommand[]>(
    selectedDashboard ? `/api/v1/channels/${selectedDashboard.channelId}/commands` : null,
    swrAuthFetcher,
  );
};

export const useUpdateCommand = () => {
  const { mutate } = useSWRConfig();
  const [selectedDashboard] = useSelectedDashboard();

  if (selectedDashboard === null) {
    throw new Error('Selected dashboard is null, unable to post command.');
  }

  return (command: ChannelCommand) =>
    mutate<ChannelCommand[]>(
      `/api/v1/channels/${selectedDashboard.channelId}/commands`,
      async (commands) => {
        const data = await swrAuthFetcher(
          `/api/v1/channels/${selectedDashboard.channelId}/commands/${command.id}`,
          {
            body: JSON.stringify(command),
            method: 'PUT',
            headers: {
              'Content-Type': 'application/json',
            },
          },
        );

        const index = commands!.findIndex((c) => c.id === data.id);
        commands![index] = data;
        return commands;
      },
      {
        revalidate: false,
      },
    );
};
