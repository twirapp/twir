import { useQuery } from '@tanstack/react-query';

export const useChangelog = () => useQuery<Commit[]>({
  queryKey: ['github/changelog'],
  queryFn: async () => {
		const call = await fetch('https://api.github.com/repos/satont/tsuwari/commits?per_page=100');

		return await call.json();
	},
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
