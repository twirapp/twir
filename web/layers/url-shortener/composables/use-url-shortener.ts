import type { ErrorModel, LinkOutputDto, ShortUrlProfileParamsSortByEnum } from '@twir/api/openapi'

import { useOapi } from '~/composables/use-oapi'

export const useUrlShortener = defineStore('url-shortener', () => {
	const api = useOapi()
	const latestShortenedUrls = ref<LinkOutputDto[]>([])

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
			if (e instanceof Response) {
				const errorData = (await e.json()) as ErrorModel
				return {
					data: null,
					error:
						errorData.detail ??
						errorData.errors?.map((e) => e.message)?.join(', ') ??
						errorData.title,
				}
			}

			return {
				data: null,
				error: 'Unknown error',
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
	}
})

if (import.meta.hot) {
	import.meta.hot.accept(acceptHMRUpdate(useUrlShortener, import.meta.hot))
}
