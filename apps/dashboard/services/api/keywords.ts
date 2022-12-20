import { ChannelKeyword } from '@tsuwari/typeorm/entities/ChannelKeyword';
import useSWR, { useSWRConfig } from 'swr';

import { mutationOptions, swrAuthFetcher } from '@/services/api';
import { useSelectedDashboard } from '@/services/dashboard';

export const useKeywordsManager = () => {
  const { mutate } = useSWRConfig();
  const [selectedDashboard] = useSelectedDashboard();

  return {
    getAll() {
      return useSWR<ChannelKeyword[]>(
        selectedDashboard ? `/api/v1/channels/${selectedDashboard.channelId}/keywords` : null,
        swrAuthFetcher,
      );
    },
    async delete(keywordId: string) {
      if (selectedDashboard === null) {
        throw new Error('Selected dashboard is null, unable to delete keyword.');
      }

      return mutate<ChannelKeyword[]>(
        `/api/v1/channels/${selectedDashboard.channelId}/keywords`,
        async (keywords) => {
          await swrAuthFetcher(
            `/api/v1/channels/${selectedDashboard.channelId}/keywords/${keywordId}`,
            {
              method: 'DELETE',
              headers: {
                'Content-Type': 'application/json',
              },
            },
          );

          return keywords?.filter((c) => c.id != keywordId);
        },
        mutationOptions,
      );
    },
    async patch(keywordId: string, keywordData: Partial<ChannelKeyword>) {
      if (selectedDashboard === null) {
        throw new Error('Selected dashboard is null, unable to delete keyword.');
      }

      return mutate<ChannelKeyword[]>(
        `/api/v1/channels/${selectedDashboard.channelId}/keywords`,
        async (keywords) => {
          const data = await swrAuthFetcher(
            `/api/v1/channels/${selectedDashboard.channelId}/keywords/${keywordId}`,
            {
              body: JSON.stringify(keywordData),
              method: 'PATCH',
              headers: {
                'Content-Type': 'application/json',
              },
            },
          );

          const index = keywords!.findIndex((c) => c.id === keywordId);
          keywords![index] = data;

          return keywords;
        },
        mutationOptions,
      );
    },
    async createOrUpdate(keyword: ChannelKeyword) {
      if (selectedDashboard === null) {
        throw new Error('Selected dashboard is null, unable to delete keyword.');
      }

      return mutate<ChannelKeyword[]>(
        `/api/v1/channels/${selectedDashboard.channelId}/keywords`,
        async (keywords) => {
          const data = await swrAuthFetcher(
            `/api/v1/channels/${selectedDashboard.channelId}/keywords/${keyword.id ?? ''}`,
            {
              body: JSON.stringify(keyword),
              method: keyword.id ? 'PUT' : 'POST',
              headers: {
                'Content-Type': 'application/json',
              },
            },
          );

          if (keyword.id) {
            const index = keywords!.findIndex((c) => c.id === data.id);
            keywords![index] = data;
          } else {
            keywords?.push(data);
          }

          return keywords;
        },
        mutationOptions,
      );
    },
  };
};
