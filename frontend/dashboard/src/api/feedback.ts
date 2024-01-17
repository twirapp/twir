import { useMutation } from '@tanstack/vue-query';

import { protectedApiClient } from '@/api/twirp';

export const useLeaveFeedback = () => useMutation({
	mutationKey: ['leaveFeedback'],
	mutationFn: async (message: string) => {
		await protectedApiClient.leaveFeedback({
			message,
		});
	},
});
