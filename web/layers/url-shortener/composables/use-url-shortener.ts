export function useUrlShortener() {
	const api = useOapi()

	async function shortUrl(opts: { url: string; alias?: string }) {
		const response = await api.v1.shortUrlCreate({
			url: opts.url,
			alias: opts.alias
		})

		return {
			data: response.data,
			error: response.error,
		}
	}

	return {
		shortUrl,
	}
}
