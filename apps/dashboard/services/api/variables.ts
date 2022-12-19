import { ChannelCustomvar } from '@tsuwari/typeorm/entities/ChannelCustomvar';
import useSWR, { useSWRConfig } from 'swr';

import { swrAuthFetcher } from '@/services/api';
import { useSelectedDashboard } from '@/services/dashboard';

export const useVariablesManager = () => {
  const { mutate } = useSWRConfig();
  const [selectedDashboard] = useSelectedDashboard();

  return {
    getBuiltin() {
      return useSWR<ChannelCustomvar[]>(
        selectedDashboard
          ? `/api/v1/channels/${selectedDashboard.channelId}/variables/builtin`
          : null,
        swrAuthFetcher,
      );
    },
    getCreated() {
      return useSWR<ChannelCustomvar[]>(
        selectedDashboard ? `/api/v1/channels/${selectedDashboard.channelId}/variables` : null,
        swrAuthFetcher,
      );
    },
    async delete(variableId: string) {
      if (selectedDashboard === null) {
        throw new Error('Selected dashboard is null, unable to delete variable.');
      }

      return mutate<ChannelCustomvar[]>(
        `/api/v1/channels/${selectedDashboard.channelId}/variables`,
        async (variables) => {
          await swrAuthFetcher(
            `/api/v1/channels/${selectedDashboard.channelId}/variables/${variableId}`,
            {
              method: 'DELETE',
              headers: {
                'Content-Type': 'application/json',
              },
            },
          );

          return variables?.filter((c) => c.id != variableId);
        },
        {
          revalidate: false,
        },
      );
    },
    async patch(variableId: string, variableData: Partial<ChannelCustomvar>) {
      if (selectedDashboard === null) {
        throw new Error('Selected dashboard is null, unable to delete keyword.');
      }

      return mutate<ChannelCustomvar[]>(
        `/api/v1/channels/${selectedDashboard.channelId}/variables`,
        async (variables) => {
          const data = await swrAuthFetcher(
            `/api/v1/channels/${selectedDashboard.channelId}/variables/${variableId}`,
            {
              body: JSON.stringify(variableData),
              method: 'PATCH',
              headers: {
                'Content-Type': 'application/json',
              },
            },
          );

          const index = variables!.findIndex((c) => c.id === variableId);
          variables![index] = data;

          return variables;
        },
        {
          revalidate: false,
        },
      );
    },
    async createOrUpdate(variable: ChannelCustomvar) {
      if (selectedDashboard === null) {
        throw new Error('Selected dashboard is null, unable to delete variable.');
      }

      return mutate<ChannelCustomvar[]>(
        `/api/v1/channels/${selectedDashboard.channelId}/variables`,
        async (variables) => {
          const data = await swrAuthFetcher(
            `/api/v1/channels/${selectedDashboard.channelId}/variables/${variable.id ?? ''}`,
            {
              body: JSON.stringify(variable),
              method: variable.id ? 'PUT' : 'POST',
              headers: {
                'Content-Type': 'application/json',
              },
            },
          );

          if (variable.id) {
            const index = variables!.findIndex((c) => c.id === data.id);
            variables![index] = data;
          } else {
            variables?.push(data);
          }

          return variables;
        },
        {
          revalidate: false,
        },
      );
    },
  };
};
