import { defineMiddleware } from 'astro/middleware';

import { protectedClient, unProtectedClient } from './api/twirp.js';

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
				// eslint-disable-next-line no-empty
				} catch {}
			}
		})(),
		(async () => {
			const request = await unProtectedClient.authGetLink({ state: Buffer.from(location, 'base64').toString('hex') });
			context.locals.authLink = request.response.link;
		})(),
	]);

	next();
});
