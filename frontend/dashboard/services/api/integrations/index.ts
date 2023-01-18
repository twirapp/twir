import { useMutation, UseMutationResult, useQuery, UseQueryResult } from '@tanstack/react-query';
import { useContext } from 'react';

import { authFetcher, queryClient } from '@/services/api';
import { SelectedDashboardContext } from '@/services/selectedDashboardProvider';

export * from './donatepay';

export interface Integration<T> {
  useData: () => UseQueryResult<T, unknown>,
  useGetAuthLink: () => UseQueryResult<string, unknown>,
  usePostCode: () => UseMutationResult<any, unknown, { code: string }, unknown>,
  useLogout: () => UseMutationResult<any, unknown, void, unknown>,
}

const createIntegrationManager = <T>(system: string): Integration<T> => {
  const dashboard = useContext(SelectedDashboardContext);
  const getUrl = (system: string) => `/api/v1/channels/${dashboard.id}/integrations/${system}`;

  return {
    useData: () => useQuery<T>({
      queryKey: [getUrl(system)],
      queryFn: () => authFetcher(getUrl(system)),
    }),
    useGetAuthLink: () => useQuery<string>({
      queryKey: [`${getUrl(system)}/auth`],
      queryFn: () => authFetcher(`${getUrl(system)}/auth`),
    }),
    usePostCode: () => useMutation({
      mutationFn: ({ code }) => {
        return authFetcher(
          getUrl(system),
          {
            body: JSON.stringify({ code }),
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
          },
        );
      },
      onSuccess: () => {
        queryClient.invalidateQueries([getUrl(system)]);
      },
      mutationKey: [getUrl(system)],
    }),
    useLogout: () => useMutation({
      mutationFn: () => {
        return authFetcher(
          `${getUrl(system)}/logout`,
          {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
          },
        );
      },
      onSuccess: () => {
        queryClient.setQueryData<T | null>([getUrl(system)], () => {
          return null;
        });
      },
      mutationKey: [getUrl(system)],
    }),
  };
};

export const useDonationAlerts = () => createIntegrationManager<{
  avatar: string;
  name: string;
}>('donationalerts');
export const useStreamlabs = () => createIntegrationManager<{
  avatar: string;
  name: string;
}>('streamlabs');
export const useFaceit = () => createIntegrationManager<{
  avatar: string;
  name: string;
}>('faceit');
export const useSpotify = () => createIntegrationManager<{
  id: string;
  display_name: string;
  images?: Array<{ url: string }>;
}>('spotify');
export const useLastfm = () => createIntegrationManager<{
  name: string;
  image: string;
  playCount: string;
}>('lastfm');
export const useVK = () => createIntegrationManager<{
  id: number;
  first_name: string;
  last_name: string;
  photo_max_orig: string;
}>('vk');
