
import { useMutation, UseMutationResult, useQuery, UseQueryResult } from '@tanstack/react-query';
import { Dashboard } from '@tsuwari/shared';
import { ChannelCommand } from '@tsuwari/typeorm/entities/ChannelCommand';
import { ChannelCommandGroup } from '@tsuwari/typeorm/entities/ChannelCommandGroup';
import { ChannelCustomvar } from '@tsuwari/typeorm/entities/ChannelCustomvar';
import { ChannelGreeting } from '@tsuwari/typeorm/entities/ChannelGreeting';
import { ChannelKeyword } from '@tsuwari/typeorm/entities/ChannelKeyword';
import { ChannelRole } from '@tsuwari/typeorm/entities/ChannelRole';
import { ChannelTimer } from '@tsuwari/typeorm/entities/ChannelTimer';
import { Event } from '@tsuwari/typeorm/entities/events/Event';
import { useContext } from 'react';

import { authFetcher } from '@/services/api';
import { queryClient } from '@/services/api';
import { SelectedDashboardContext } from '@/services/selectedDashboardProvider';

export type Greeting = ChannelGreeting & { userName: string, avatar?: string };

interface Crud<T> {
  useGetAll: () => UseQueryResult<T[], unknown>
  useDelete: () => UseMutationResult<any, unknown, string, unknown>
  usePatch: () => UseMutationResult<any, unknown, {id: string, data: Partial<T>}, unknown>
  useCreateOrUpdate: () => UseMutationResult<any, unknown, {id?: string | undefined, data: T}, unknown>,
}

const createCrudManager = <T extends { id: string }>(system: string, opts?: {
  additionalSystemsKeys?: string[],
}): Crud<T> => {
  const dashboard = useContext(SelectedDashboardContext);

  const getUrl = (system: string) => `/api/v1/channels/${dashboard.id}/${system}`;

  return {
    useGetAll: () => useQuery<T[]>({
      queryKey: [getUrl(system)],
      queryFn: () => authFetcher(getUrl(system)),
    }),
    useDelete: () => useMutation({
      mutationFn: (id: string) => {
        return authFetcher(
          `${getUrl(system)}/${id}`,
          {
            method: 'DELETE',
            headers: {
              'Content-Type': 'application/json',
            },
          },
        );
      },
      onSuccess: (result, id, context) => {
        queryClient.invalidateQueries([getUrl(system)]);

        for (const additionalSystem of opts?.additionalSystemsKeys ?? []) {
          queryClient.refetchQueries([getUrl(additionalSystem)]);
        }
      },
      mutationKey: [getUrl(system)],
    }),
    usePatch: () => useMutation({
      mutationFn: ({ id, data }: { id: string, data: Partial<T> }) => {
        return authFetcher(
          `${getUrl(system)}/${id}`,
          {
            body: JSON.stringify(data),
            method: 'PATCH',
            headers: {
              'Content-Type': 'application/json',
            },
          },
        );
      },
      onSuccess: (result, data) => {
        queryClient.invalidateQueries([getUrl(system)]);

        for (const additionalSystem of opts?.additionalSystemsKeys ?? []) {
          queryClient.refetchQueries([getUrl(additionalSystem)]);
        }
      },
      mutationKey: [getUrl(system)],
    }),
    useCreateOrUpdate: () => useMutation({
      mutationFn: ({ id, data }: { id?: string, data: T }) => {
        return authFetcher(
          `${getUrl(system)}/${id ?? ''}`,
          {
            body: JSON.stringify(data),
            method: id ? 'PUT' : 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
          },
        );
      },
      onSuccess: (result, data) => {
        queryClient.invalidateQueries([getUrl(system)]);
        for (const additionalSystem of opts?.additionalSystemsKeys ?? []) {
          queryClient.refetchQueries([getUrl(additionalSystem)]);
        }
      },
      mutationKey: [getUrl(system)],
    }),
  };
};

export const commandsManager = () => createCrudManager<ChannelCommand>('commands');
export const greetingsManager = () => createCrudManager<Greeting>('greetings');
export const keywordsManager =  () => createCrudManager<ChannelKeyword>('keywords');
export const timersManager =  () => createCrudManager<ChannelTimer>('timers');
export const variablesManager = () => createCrudManager<ChannelCustomvar>('variables');
export const dashboardAccessManager = () => createCrudManager<Dashboard>('settings/dashboard-access');
export const eventsManager = () => createCrudManager<Event>('events');
export const commandsGroupManager = () => createCrudManager<ChannelCommandGroup>('commands/groups', {
  additionalSystemsKeys: ['commands'],
});
export const useRolesApi = () => createCrudManager<ChannelRole>('roles');