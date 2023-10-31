import { useQuery } from '@tanstack/vue-query';

import { unprotectedApiClient } from '@/api/twirp';

export const useGoogleFontsList = () => useQuery({
	queryKey: ['googleFontsList'],
	queryFn: async () => {
		const call = unprotectedApiClient.getGoogleFonts({});
		return call.response;
	},
});
