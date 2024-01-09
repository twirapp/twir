import { unProtectedClient } from '@/api/twirp.js';

export async function getStats() {
	try {
		const { response } = await unProtectedClient.getStats({}, { timeout: 1000 });
		return response;
	} catch (err) {
		console.error(err);
	}
}
