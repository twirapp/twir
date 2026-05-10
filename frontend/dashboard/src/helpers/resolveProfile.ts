export interface ResolvedProfile {
	avatar: string
	displayName: string
	login: string
	url: string
	platform: string
	notFound: boolean
}

export interface ProfileInput {
	profileImageUrl?: string | null
	displayName?: string | null
	login?: string | null
	platform?: string | null
	notFound?: boolean | null
}

function fallback<T>(value: T | null | undefined, defaultValue: T): T {
	if (value === null || value === undefined || value === '') {
		return defaultValue
	}
	return value
}

export function resolveProfile(profile: ProfileInput): ResolvedProfile {
	const platform = profile.platform ?? 'twitch'
	const login = fallback(profile.login, '')
	const displayName = fallback(profile.displayName, login)
	const avatar = fallback(profile.profileImageUrl, '')
	const notFound = profile.notFound ?? false

	let url = ''
	if (login && !notFound && !login.startsWith('[')) {
		url = platform === 'kick'
			? `https://kick.com/${login}`
			: `https://twitch.tv/${login}`
	}

	return { avatar, displayName, login, url, platform, notFound }
}
