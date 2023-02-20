import { useQuery } from '@tanstack/react-query';

import { fetcher } from '@/services/api/fetchWrappers';

export const useChangelog = () => useQuery<Commit[]>({
  queryKey: ['github/changelog'],
  queryFn: () => fetcher('https://api.github.com/repos/satont/tsuwari/commits?per_page=100'),
});

export type Commit = {
  sha: string,
  commit: {
    author: {
      name: string,
      date: string,
    },
    message: string,
  }
}