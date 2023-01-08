import { useQuery } from '@tanstack/react-query';
import { type HelixUserData } from '@twurple/api';

export const useUsersByNames = (names: string[]) => useQuery({
  queryKey: ['users', ...names],
  queryFn: async ({ queryKey }): Promise<HelixUserData[]> => {
    const params = new URLSearchParams();
    params.set('names', queryKey.slice(1)!.join(','));

    const data = await fetch(`/api/v1/twitch/users?${params}`);

    return data.json();
  },
});