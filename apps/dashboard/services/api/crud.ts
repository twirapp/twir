
import { useMutation, UseMutationResult, useQuery, UseQueryResult } from '@tanstack/react-query';
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

interface Crud<T> {
  getAll: UseQueryResult<T[], unknown>
  delete: UseMutationResult<any, unknown, string, unknown>
  patch: UseMutationResult<any, unknown, {id: string, data: Partial<T>}, unknown>
  createOrUpdate: UseMutationResult<any, unknown, {id?: string | undefined, data: T}, unknown>
}

// const cachedCruds: Map<string, Crud<any>> = new Map();

const createCrudManager = <T extends { id: string }>(system: string): Crud<T> => {
  // console.log(system, cachedCruds.has(system));
  // if (cachedCruds.has(system)) {
  //   return cachedCruds.get(system) as Crud<T>;
  // }
  //
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
      onSuccess: (result, data) => {
        queryClient.setQueryData<T[]>([getUrl(system)], old => {
          if (!old) {
            return [result];
          }
          const index = old?.findIndex(o => o.id === data.id);
          old[index!] = result;
          return old;
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
      onSuccess: (result, data) => {
        queryClient.setQueryData<T[]>([getUrl(system)], old => {
          if (!old) {
            return [result];
          }
          const index = old?.findIndex(o => o.id === data.id);
          if (index && index != -1) {
            old[index!] = result;
          } else {
            old.push(result);
          }
          console.log(old);
          return old;
        });
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