import { useMutation, useQuery } from '@tanstack/react-query';
import { ChannelRoleUser } from '@twir/typeorm/entities/ChannelRoleUser';
import { useContext } from 'react';

import { authFetcher } from '@/services/api/fetchWrappers';
import { queryClient } from '@/services/api/queryClient';
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
    useUpdate: () => useMutation({
      mutationFn: (data: { userNames: string[], roleId: string }) => {
        return authFetcher(
          `/api/v1/channels/${dashboard.id}/roles/${data.roleId}/users`,
          {
            method: 'PUT',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({ userNames: data.userNames }),
          },
        );
      },
      onSuccess: (d, data) => {
        queryClient.refetchQueries([`/api/v1/channels/${dashboard.id}/roles/${data.roleId}/users`]);
      },
    }),
  };
};
