import { ChannelCommand } from '@tsuwari/typeorm/entities/ChannelCommand';
import useSWR, { useSWRConfig } from 'swr';

import { swrAuthFetcher } from '@/services/api';
import { useSelectedDashboard } from '@/services/dashboard';

export const useCommandManager = () => {
  const { mutate } = useSWRConfig();
  const [selectedDashboard] = useSelectedDashboard();

  return {
    getAll() {
      return useSWR<ChannelCommand[]>(
        selectedDashboard ? `/api/v1/channels/${selectedDashboard.channelId}/commands` : null,
        swrAuthFetcher,
      );
    },
    async delete(commandId: string) {
      if (selectedDashboard === null) {
        throw new Error('Selected dashboard is null, unable to delete command.');
      }

      return mutate<ChannelCommand[]>(
        `/api/v1/channels/${selectedDashboard.channelId}/commands`,
        async (commands) => {
          await swrAuthFetcher(
            `/api/v1/channels/${selectedDashboard.channelId}/commands/${commandId}`,
            {
              method: 'DELETE',
              headers: {
                'Content-Type': 'application/json',
              },
            },
          );

          return commands?.filter((c) => c.id != commandId);
        },
        {
          revalidate: false,
        },
      );
    },
    async patch(commandId: string, commandData: Partial<ChannelCommand>) {
      if (selectedDashboard === null) {
        throw new Error('Selected dashboard is null, unable to delete command.');
      }

      return mutate<ChannelCommand[]>(
        `/api/v1/channels/${selectedDashboard.channelId}/commands`,
        async (commands) => {
          const data = await swrAuthFetcher(
            `/api/v1/channels/${selectedDashboard.channelId}/commands/${commandId}`,
            {
              body: JSON.stringify(commandData),
              method: 'PATCH',
              headers: {
                'Content-Type': 'application/json',
              },
            },
          );

          const index = commands!.findIndex((c) => c.id === commandId);
          commands![index] = data;

          return commands;
        },
        {
          revalidate: false,
        },
      );
    },
    async createOrUpdate(command: ChannelCommand) {
      if (selectedDashboard === null) {
        throw new Error('Selected dashboard is null, unable to delete command.');
      }

      return mutate<ChannelCommand[]>(
        `/api/v1/channels/${selectedDashboard.channelId}/commands`,
        async (commands) => {
          const data = await swrAuthFetcher(
            `/api/v1/channels/${selectedDashboard.channelId}/commands/${command.id ?? ''}`,
            {
              body: JSON.stringify(command),
              method: command.id ? 'PUT' : 'POST',
              headers: {
                'Content-Type': 'application/json',
              },
            },
          );

          if (command.id) {
            const index = commands!.findIndex((c) => c.id === data.id);
            commands![index] = data;
          } else {
            commands?.push(data);
          }

          return commands;
        },
        {
          revalidate: false,
        },
      );
    },
  };
};
