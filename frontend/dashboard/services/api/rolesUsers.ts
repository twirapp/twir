import { useMutation, useQuery } from '@tanstack/react-query';
import { ChannelRoleUser } from '@tsuwari/typeorm/entities/ChannelRoleUser';
import { useContext } from 'react';

import { authFetcher } from '@/services/api/fetchWrappers';
import { SelectedDashboardContext } from '@/services/selectedDashboardProvider';


export type RolesUser = ChannelRoleUser & { userName: string, userAvatar: string, userDisplayName: string };

export const useRolesUsers = () => {
  const dashboard = useContext(SelectedDashboardContext);

  return {
    useGetAll: (roleId: string) => useQuery<RolesUser[]>({
      queryKey: [`/api/v1/channels/${dashboard.id}/roles/${roleId}/users`],
      queryFn: () => {
        if (roleId === '')  {
          return [];
        }

        return authFetcher(`/api/v1/channels/${dashboard.id}/roles/${roleId}/users`);
      },
    }),
    useUpdate: (roleId: string) => useMutation({
      mutationFn: (data: { userNames: string }) => {
        return authFetcher(
          `/api/v1/channels/${dashboard.id}/roles/${roleId}/users`,
          {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
          },
        );
      },
    }),
  };
};