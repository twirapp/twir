import { defineMiddleware } from 'astro/middleware';

import { protectedClient, unProtectedClient } from '@/api/twirp.js';

export const onRequest = defineMiddleware(async (context, next) => {
	const session = context.cookies.get('session');
	const location = context.url.origin;

	await Promise.all([
		(async () => {
			if (session && session.value) {
				try {
					const request = await protectedClient.authUserProfile({}, {
						meta: { Cookie: `session=${session.value}` },
					});
					context.locals.profile = request.response;
				} catch (err) {
					console.log('User profile error:', err);
				}
			}
		})(),
		(async () => {
			const state = Buffer.from(location, 'base64').toString('hex');
			const request = await unProtectedClient.authGetLink({ state });
			context.locals.authLink = request.response.link;
		})(),
	]);

	next();
});
