
import { useMutation, useQuery } from '@tanstack/react-query';
import { Dashboard } from '@tsuwari/shared';
import { ChannelCommand } from '@tsuwari/typeorm/entities/ChannelCommand';
import { ChannelCustomvar } from '@tsuwari/typeorm/entities/ChannelCustomvar';
import { ChannelGreeting } from '@tsuwari/typeorm/entities/ChannelGreeting';
import { ChannelKeyword } from '@tsuwari/typeorm/entities/ChannelKeyword';
import { ChannelTimer } from '@tsuwari/typeorm/entities/ChannelTimer';
import { getCookie } from 'cookies-next';

import { authFetcher } from '@/services/api';
import { queryClient } from '@/services/api';
import { SELECTED_DASHBOARD_KEY } from '@/services/dashboard';

export type Greeting = ChannelGreeting & { userName: string };

const getUrl = (system: string) => `/api/v1/channels/${getCookie(SELECTED_DASHBOARD_KEY)}/${system}`;

const createCrudManager = <T>(system: string) => {
  return {
    getAll: useQuery<T[]>({
      queryKey: [getUrl(system)],
      queryFn: () => authFetcher(getUrl(system)),
    }),
    delete: useMutation({
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
        queryClient.setQueryData<ChannelCommand[]>([getUrl(system)], old => {
          return old?.filter(c => c.id !== id);
        });
      },
      mutationKey: [getUrl(system)],
    }),
    patch: useMutation({
      mutationFn: ({ id, data }: { id: string, data:  Partial<T> }) => {
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
      onSuccess: (result) => {
        queryClient.setQueryData<T[]>([getUrl(system)], old => {
          return [...old ?? [], result];
        });
      },
      mutationKey: [getUrl(system)],
    }),
    createOrUpdate: useMutation({
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
      onSuccess: (result) => {
        queryClient.setQueryData<T[]>([getUrl(system)], old => {
          return [...old ?? [], result];
        });
      },
      mutationKey: [getUrl(system)],
    }),
  };
};

export const commandsManager = createCrudManager<ChannelCommand>('commands');
export const greetingsManager = createCrudManager<Greeting>('greetings');
export const keywordsManager =  createCrudManager<ChannelKeyword>('keywords');
export const timersManager =  createCrudManager<ChannelTimer>('timers');
export const variablesManager = createCrudManager<ChannelCustomvar>('variables');
export const dashboardAccessManager = createCrudManager<Dashboard>('dashboard-access');