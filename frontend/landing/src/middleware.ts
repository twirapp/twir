import { defineMiddleware } from 'astro/middleware'

import type { APIContext } from 'astro'

import { getAuthLink } from '@/api/auth-link.ts'
import { getAuthenticatedUser } from '@/api/user.ts'

export const onRequest = defineMiddleware(async (context, next) => {
	await Promise.all([
		assignProfile(context),
		assignLoginLink(context),
	])

	await next()
})

async function assignProfile(context: APIContext) {
	const session = context.cookies.get('session')

	if (session && session.value) {
		try {
			context.locals.profile = await getAuthenticatedUser(session.value)
		} catch (err) {
			console.log('User profile error:', err)
		}
	}
}

async function assignLoginLink(context: APIContext) {
	let redirectTo = `${context.url.origin}/dashboard`
	redirectTo = redirectTo
		.replace(':80', '')
		.replace(':443', '')

	try {
		context.locals.authLink = await getAuthLink(redirectTo)
	} catch { /* empty */
	}
}
