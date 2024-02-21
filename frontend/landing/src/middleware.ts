import type { APIContext } from 'astro';
import { defineMiddleware } from 'astro/middleware';

import { protectedClient, unProtectedClient } from '@/api/twirp.js';

export const onRequest = defineMiddleware(async (context, next) => {
	await Promise.all([
		assignProfile(context),
		assignLoginLink(context),
	]);

	await next();
});

const assignProfile = async (context: APIContext) => {
	const session = context.cookies.get('session');

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
};

const assignLoginLink = async (context: APIContext) => {
	const redirectTo = `${context.url.origin}/dashboard`;

	try {
		const request = await unProtectedClient.authGetLink({ redirectTo });
		context.locals.authLink = request.response.link;
	} catch { /* empty */
	}
};
