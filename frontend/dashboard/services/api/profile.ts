import { useMutation, useQuery } from '@tanstack/react-query';
import { deleteCookie } from 'cookies-next';
import { UnwrapPromise } from 'next/dist/lib/coalesced-function';

import { protectedApiClient } from '@/services/api/twirp';

export const useProfile = () =>
  useQuery<UnwrapPromise<ReturnType<typeof protectedApiClient.authUserProfile>['response']>>({
    queryKey: [`/api/auth/profile`],
    queryFn: async () => {
			const call = await protectedApiClient.authUserProfile({});
			return call.response;
		},
    retry: false,
  });

export const useLogoutMutation = () =>
  useMutation({
    mutationFn: async () => {
			await protectedApiClient.authLogout({});
    },
    onSuccess() {
      localStorage.removeItem('access_token');
      deleteCookie('dashboard_id');
      window.location.replace(window.location.origin);
    },
  });

type DashboardsResponse = UnwrapPromise<ReturnType<typeof protectedApiClient.authGetDashboards>['response']>

export const useDashboards = () => useQuery<DashboardsResponse>({
  queryKey: [`/api/auth/profile/dashboards`],
  queryFn: async () => {
		const call = await protectedApiClient.authGetDashboards({});
		return call.response;
	},
  retry: false,
});
