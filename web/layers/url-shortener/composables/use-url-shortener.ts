import type { ErrorModel } from '@twir/api/openapi'

export function useUrlShortener() {
	const api = useOapi()

	async function shortUrl(opts: { url: string; alias?: string }) {
		try {
			const response = await api.v1.shortUrlCreate({
				url: opts.url,
				alias: opts.alias,
			})

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

	return {
		shortUrl,
	}
}
