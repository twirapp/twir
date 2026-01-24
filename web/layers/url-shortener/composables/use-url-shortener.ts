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

export const useUrlShortener = defineStore('url-shortener', () => {
	const api = useOapi()
	const latestShortenedUrls = ref<LinkOutputDto[]>([])
	const customDomain = ref<CustomDomainOutputDto | null>(null)
	const isCustomDomainLoading = ref(false)

	async function shortUrl(opts: { url: string; alias?: string }) {
		try {
			const response = await api.v1.shortUrlCreate({
				url: opts.url,
				alias: opts.alias,
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
			sortBy: 'views' as ShortUrlProfileParamsSortByEnum,
		}
	) {
		const response = await api.v1.shortUrlProfile(opts)

		latestShortenedUrls.value = response.data.data.items

		return {
			data: response.data?.data,
			error: response.error,
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
