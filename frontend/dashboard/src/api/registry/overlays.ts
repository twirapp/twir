import { useMutation } from '@tanstack/vue-query';

import { protectedApiClient } from '@/api/twirp';

export const useOverlaysParseHtml = () => useMutation({
	mutationFn: async (htmlString: string) => {
		if (!htmlString) {
			return '';
		}
		const req = await protectedApiClient.overlaysParseHtml({
			html: btoa(htmlString),
		});

		return atob(req.response.html);
	},
});
