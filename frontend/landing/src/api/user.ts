import { gqlUrl } from '@/api/gql.ts'

export async function getAuthenticatedUser(session: string) {
	const request = await fetch(gqlUrl, {
		headers: {
			'Cookie': `session=${session}`,
			'Content-Type': 'application/json',
		},
		credentials: 'include',
		body: '{"query":"query {\\n  authenticatedUser {\\n    twitchProfile {\\n      displayName\\n      profileImageUrl\\n    }\\n  }\\n}"}',
		method: 'POST',
	})

	const response = await request.json()

	if (!request.ok || response.errors) {
		console.log(response)
		throw new Error(response.errors.toString())
	}

	const profile = response.data.authenticatedUser.twitchProfile

	return profile
}
