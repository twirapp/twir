import type { APIContext } from 'astro';
import { defineMiddleware } from 'astro/middleware';

import { getAuthLink } from '@/api/auth-link.ts';
import { getAuthenticatedUser } from '@/api/user.ts';

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
			context.locals.profile = await getAuthenticatedUser(session.value);
		} catch (err) {
			console.log('User profile error:', err);
		}
	}
};

const assignLoginLink = async (context: APIContext) => {
	const redirectTo = `${context.url.origin}/dashboard`;

	console.dir({ redirectTo });
	try {
		context.locals.authLink = await getAuthLink(redirectTo);
	} catch { /* empty */
	}
};
