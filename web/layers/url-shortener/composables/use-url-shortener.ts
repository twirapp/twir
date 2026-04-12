import { ContentType } from '@twir/api/openapi'
import type { ErrorModel, LinkOutputDto, ShortUrlProfileParamsSortByEnum } from '@twir/api/openapi'

import { useOapi } from '~/composables/use-oapi'

type CustomDomainOutputDto = {
	id: string
	domain: string
	verified: boolean
	verification_token: string
	verification_target: string
	created_at: string
}

type CustomDomainResponseDto = {
	$schema?: string
	data: CustomDomainOutputDto
}

export type BannedUserAgentDto = {
	id: string
	pattern: string
	description: string | null
	created_at: string
}

type BannedUserAgentsListResponseDto = {
	$schema?: string
	data: BannedUserAgentDto[]
}

type BannedUserAgentResponseDto = {
	$schema?: string
	data: BannedUserAgentDto
}

export type LinkBannedUserAgentDto = {
	id: string
	link_id: string
	pattern: string
	description: string | null
	created_at: string
}

type LinkBannedUserAgentsListResponseDto = {
	$schema?: string
	data: LinkBannedUserAgentDto[]
}

type LinkBannedUserAgentResponseDto = {
	$schema?: string
	data: LinkBannedUserAgentDto
}

export const useUrlShortener = defineStore('url-shortener', () => {
	const api = useOapi()
	const latestShortenedUrls = ref<LinkOutputDto[]>([])
	const customDomain = ref<CustomDomainOutputDto | null>(null)
	const isCustomDomainLoading = ref(false)

	const globalBannedUserAgents = ref<BannedUserAgentDto[]>([])
	const isGlobalBannedUserAgentsLoading = ref(false)

	const perLinkBannedUserAgents = ref<Map<string, LinkBannedUserAgentDto[]>>(new Map())
	const isPerLinkBannedUserAgentsLoading = ref<Map<string, boolean>>(new Map())

	async function shortUrl(opts: { url: string; alias?: string; useCustomDomain?: boolean }) {
		try {
			const response = await api.v1.shortUrlCreate({
				url: opts.url,
				alias: opts.alias,
				use_custom_domain: opts.useCustomDomain,
			})

			latestShortenedUrls.value = [response.data.data, ...latestShortenedUrls.value.slice(0, 2)]

			return {
				data: response.data,
				error: response.error,
			}
		} catch (e) {
			return {
				data: null,
				error: await parseApiError(e),
			}
		}
	}

	async function fetchCustomDomain() {
		isCustomDomainLoading.value = true

		try {
			const response = await api.http.request<CustomDomainResponseDto, ErrorModel>({
				path: '/v1/short-links/custom-domain',
				method: 'GET',
				format: 'json',
			})

			customDomain.value = response.data.data

			return {
				data: response.data,
				error: response.error,
			}
		} catch (e) {
			if (e instanceof Response && e.status === 404) {
				customDomain.value = null
				return {
					data: null,
					error: null,
				}
			}

			return {
				data: null,
				error: await parseApiError(e),
			}
		} finally {
			isCustomDomainLoading.value = false
		}
	}

	async function createCustomDomain(domain: string) {
		try {
			const response = await api.http.request<CustomDomainResponseDto, ErrorModel>({
				path: '/v1/short-links/custom-domain',
				method: 'POST',
				type: ContentType.Json,
				format: 'json',
				body: {
					domain,
				},
			})

			customDomain.value = response.data.data

			return {
				data: response.data,
				error: response.error,
			}
		} catch (e) {
			return {
				data: null,
				error: await parseApiError(e),
			}
		}
	}

	async function verifyCustomDomain() {
		try {
			const response = await api.http.request<CustomDomainResponseDto, ErrorModel>({
				path: '/v1/short-links/custom-domain/verify',
				method: 'POST',
				format: 'json',
			})

			customDomain.value = response.data.data

			return {
				data: response.data,
				error: response.error,
			}
		} catch (e) {
			return {
				data: null,
				error: await parseApiError(e),
			}
		}
	}

	async function deleteCustomDomain() {
		try {
			const response = await api.http.request<Record<string, never>, ErrorModel>({
				path: '/v1/short-links/custom-domain',
				method: 'DELETE',
				format: 'json',
			})

			customDomain.value = null

			return {
				data: response.data,
				error: response.error,
			}
		} catch (e) {
			return {
				data: null,
				error: await parseApiError(e),
			}
		}
	}

	async function refetchLatestShortenedUrls(
		opts: { page?: number; perPage?: number; sortBy?: ShortUrlProfileParamsSortByEnum } = {
			page: 0,
			perPage: 3,
			sortBy: 'created_at' as ShortUrlProfileParamsSortByEnum,
		}
	) {
		try {
			const response = await api.v1.shortUrlProfile(opts)

			latestShortenedUrls.value = response.data.data.items

			return {
				data: response.data?.data,
				error: response.error,
			}
		} catch (e) {
			return {
				data: null,
				error: await parseApiError(e),
			}
		}
	}

	async function fetchGlobalBannedUserAgents() {
		isGlobalBannedUserAgentsLoading.value = true
		try {
			const response = await api.http.request<BannedUserAgentsListResponseDto, ErrorModel>({
				path: '/v1/short-links/banned-user-agents',
				method: 'GET',
				format: 'json',
			})
			globalBannedUserAgents.value = response.data.data
			return { data: response.data, error: response.error }
		} catch (e) {
			return { data: null, error: await parseApiError(e) }
		} finally {
			isGlobalBannedUserAgentsLoading.value = false
		}
	}

	async function createGlobalBannedUserAgent(opts: { pattern: string; description?: string | null }) {
		try {
			const response = await api.http.request<BannedUserAgentResponseDto, ErrorModel>({
				path: '/v1/short-links/banned-user-agents',
				method: 'POST',
				type: ContentType.Json,
				format: 'json',
				body: {
					pattern: opts.pattern,
					description: opts.description,
				},
			})
			globalBannedUserAgents.value = [...globalBannedUserAgents.value, response.data.data]
			return { data: response.data, error: response.error }
		} catch (e) {
			return { data: null, error: await parseApiError(e) }
		}
	}

	async function deleteGlobalBannedUserAgent(id: string) {
		try {
			const response = await api.http.request<Record<string, never>, ErrorModel>({
				path: `/v1/short-links/banned-user-agents/${id}`,
				method: 'DELETE',
				format: 'json',
			})
			globalBannedUserAgents.value = globalBannedUserAgents.value.filter((item) => item.id !== id)
			return { data: response.data, error: response.error }
		} catch (e) {
			return { data: null, error: await parseApiError(e) }
		}
	}

	async function fetchPerLinkBannedUserAgents(linkId: string) {
		isPerLinkBannedUserAgentsLoading.value.set(linkId, true)
		try {
			const response = await api.http.request<LinkBannedUserAgentsListResponseDto, ErrorModel>({
				path: `/v1/short-links/${linkId}/banned-user-agents`,
				method: 'GET',
				format: 'json',
			})
			perLinkBannedUserAgents.value.set(linkId, response.data.data)
			return { data: response.data, error: response.error }
		} catch (e) {
			return { data: null, error: await parseApiError(e) }
		} finally {
			isPerLinkBannedUserAgentsLoading.value.set(linkId, false)
		}
	}

	async function createPerLinkBannedUserAgent(linkId: string, opts: { pattern: string; description?: string | null }) {
		try {
			const response = await api.http.request<LinkBannedUserAgentResponseDto, ErrorModel>({
				path: `/v1/short-links/${linkId}/banned-user-agents`,
				method: 'POST',
				type: ContentType.Json,
				format: 'json',
				body: {
					pattern: opts.pattern,
					description: opts.description,
				},
			})
			const current = perLinkBannedUserAgents.value.get(linkId) || []
			perLinkBannedUserAgents.value.set(linkId, [...current, response.data.data])
			return { data: response.data, error: response.error }
		} catch (e) {
			return { data: null, error: await parseApiError(e) }
		}
	}

	async function deletePerLinkBannedUserAgent(linkId: string, id: string) {
		try {
			const response = await api.http.request<Record<string, never>, ErrorModel>({
				path: `/v1/short-links/${linkId}/banned-user-agents/${id}`,
				method: 'DELETE',
				format: 'json',
			})
			const current = perLinkBannedUserAgents.value.get(linkId) || []
			perLinkBannedUserAgents.value.set(linkId, current.filter((item) => item.id !== id))
			return { data: response.data, error: response.error }
		} catch (e) {
			return { data: null, error: await parseApiError(e) }
		}
	}

	return {
		shortUrl,
		refetchLatestShortenedUrls,
		latestShortenedUrls,
		fetchCustomDomain,
		createCustomDomain,
		verifyCustomDomain,
		deleteCustomDomain,
		customDomain,
		isCustomDomainLoading,
		globalBannedUserAgents,
		isGlobalBannedUserAgentsLoading,
		fetchGlobalBannedUserAgents,
		createGlobalBannedUserAgent,
		deleteGlobalBannedUserAgent,
		perLinkBannedUserAgents,
		isPerLinkBannedUserAgentsLoading,
		fetchPerLinkBannedUserAgents,
		createPerLinkBannedUserAgent,
		deletePerLinkBannedUserAgent,
	}
})

async function parseApiError(error: unknown) {
	if (error instanceof Response) {
		const errorData = (await error.json()) as ErrorModel
		return (
			errorData.detail ??
			errorData.errors?.map((err) => err.message)?.join(', ') ??
			errorData.title ??
			'Unknown error'
		)
	}

	return 'Unknown error'
}

if (import.meta.hot) {
	import.meta.hot.accept(acceptHMRUpdate(useUrlShortener, import.meta.hot))
}
