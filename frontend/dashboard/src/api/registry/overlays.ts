import { useMutation } from '@tanstack/vue-query';

import { protectedApiClient } from '@/api/twirp';

function fromBinary(binary: string) {
	const bytes = new Uint8Array(binary.length);
	for (let i = 0; i < bytes.length; i++) {
		bytes[i] = binary.charCodeAt(i);
	}
	return String.fromCharCode(...new Uint16Array(bytes.buffer));
}

export const useOverlaysParseHtml = () => useMutation({
	mutationFn: async (htmlString: string) => {
		if (!htmlString) {
			return '';
		}
		const req = await protectedApiClient.overlaysParseHtml({
			html: btoa(htmlString),
		});

		return fromBinary(req.response.html);
	},
});
