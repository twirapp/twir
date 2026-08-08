import type {
	CustomDomainOutputDto,
	ErrorModel,
	LinkBannedUserAgentDto,
	LinkOutputDto,
	LinkPresetDto,
	PresetDto,
	PresetPatternDto,
	ShortUrlProfileParamsSortByEnum,
} from '@twir/api/openapi'

import { useOapi } from '~/composables/use-oapi'

export const useUrlShortener = defineStore('url-shortener', () => {
	const api = useOapi()
	const latestShortenedUrls = ref<LinkOutputDto[]>([])
	const customDomain = ref<CustomDomainOutputDto | null>(null)
	const isCustomDomainLoading = ref(false)

	const perLinkBannedUserAgents = ref<Map<string, LinkBannedUserAgentDto[]>>(new Map())
	const isPerLinkBannedUserAgentsLoading = ref<Map<string, boolean>>(new Map())

	const bannedUaPresets = ref<PresetDto[]>([])
	const isBannedUaPresetsLoading = ref(false)
	const presetPatterns = ref<Map<string, PresetPatternDto[]>>(new Map())
	const isPresetPatternsLoading = ref<Map<string, boolean>>(new Map())

	const linkPresets = ref<Map<string, LinkPresetDto[]>>(new Map())
	const isLinkPresetsLoading = ref<Map<string, boolean>>(new Map())

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
			const response = await api.v1.shortLinksGetCustomDomain()

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
			const response = await api.v1.shortLinksCreateCustomDomain({ domain })

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
			const response = await api.v1.shortLinksVerifyCustomDomain()

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
			const response = await api.v1.shortLinksDeleteCustomDomain()

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

	async function fetchPerLinkBannedUserAgents(linkId: string) {
		isPerLinkBannedUserAgentsLoading.value.set(linkId, true)
		try {
			const response = await api.v1.shortLinksListLinkBannedUserAgents(linkId)
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
			const response = await api.v1.shortLinksCreateLinkBannedUserAgent(linkId, {
				pattern: opts.pattern,
				description: opts.description ?? undefined,
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
			const response = await api.v1.shortLinksDeleteLinkBannedUserAgent(linkId, id)
			const current = perLinkBannedUserAgents.value.get(linkId) || []
			perLinkBannedUserAgents.value.set(linkId, current.filter((item) => item.id !== id))
			return { data: response.data, error: response.error }
		} catch (e) {
			return { data: null, error: await parseApiError(e) }
		}
	}

	async function fetchBannedUaPresets() {
		isBannedUaPresetsLoading.value = true
		try {
			const response = await api.v1.shortLinksListPresets()
			bannedUaPresets.value = response.data.data
			return { data: response.data, error: response.error }
		} catch (e) {
			return { data: null, error: await parseApiError(e) }
		} finally {
			isBannedUaPresetsLoading.value = false
		}
	}

	async function createBannedUaPreset(opts: { name: string; description?: string | null }) {
		try {
			const response = await api.v1.shortLinksCreatePreset({
				name: opts.name,
				description: opts.description ?? undefined,
			})
			bannedUaPresets.value = [...bannedUaPresets.value, response.data.data]
			return { data: response.data, error: response.error }
		} catch (e) {
			return { data: null, error: await parseApiError(e) }
		}
	}

	async function updateBannedUaPreset(
		presetId: string,
		opts: { name?: string; description?: string | null }
	) {
		try {
			const response = await api.v1.shortLinksUpdatePreset(presetId, {
				name: opts.name,
				description: opts.description ?? undefined,
			})
			bannedUaPresets.value = bannedUaPresets.value.map((preset) =>
				preset.id === presetId ? response.data.data : preset
			)
			return { data: response.data, error: response.error }
		} catch (e) {
			return { data: null, error: await parseApiError(e) }
		}
	}

	async function deleteBannedUaPreset(presetId: string) {
		try {
			const response = await api.v1.shortLinksDeletePreset(presetId)
			bannedUaPresets.value = bannedUaPresets.value.filter((preset) => preset.id !== presetId)
			presetPatterns.value.delete(presetId)
			for (const [linkId, presets] of linkPresets.value) {
				linkPresets.value.set(
					linkId,
					presets.filter((item) => item.preset_id !== presetId)
				)
			}
			return { data: response.data, error: response.error }
		} catch (e) {
			return { data: null, error: await parseApiError(e) }
		}
	}

	async function fetchPresetPatterns(presetId: string) {
		isPresetPatternsLoading.value.set(presetId, true)
		try {
			const response = await api.v1.shortLinksListPresetPatterns(presetId)
			presetPatterns.value.set(presetId, response.data.data)
			return { data: response.data, error: response.error }
		} catch (e) {
			return { data: null, error: await parseApiError(e) }
		} finally {
			isPresetPatternsLoading.value.set(presetId, false)
		}
	}

	async function createPresetPattern(
		presetId: string,
		opts: { pattern: string; description?: string | null }
	) {
		try {
			const response = await api.v1.shortLinksCreatePresetPattern(presetId, {
				pattern: opts.pattern,
				description: opts.description ?? undefined,
			})
			const current = presetPatterns.value.get(presetId) || []
			presetPatterns.value.set(presetId, [...current, response.data.data])
			return { data: response.data, error: response.error }
		} catch (e) {
			return { data: null, error: await parseApiError(e) }
		}
	}

	async function deletePresetPattern(presetId: string, patternId: string) {
		try {
			const response = await api.v1.shortLinksDeletePresetPattern(presetId, patternId)
			const current = presetPatterns.value.get(presetId) || []
			presetPatterns.value.set(
				presetId,
				current.filter((item) => item.id !== patternId)
			)
			return { data: response.data, error: response.error }
		} catch (e) {
			return { data: null, error: await parseApiError(e) }
		}
	}

	async function fetchLinkPresets(linkId: string) {
		isLinkPresetsLoading.value.set(linkId, true)
		try {
			const response = await api.v1.shortLinksListLinkPresets(linkId)
			linkPresets.value.set(linkId, response.data.data)
			return { data: response.data, error: response.error }
		} catch (e) {
			return { data: null, error: await parseApiError(e) }
		} finally {
			isLinkPresetsLoading.value.set(linkId, false)
		}
	}

	async function applyPresetToLink(linkId: string, presetId: string) {
		try {
			const response = await api.v1.shortLinksApplyPresetToLink(linkId, {
				preset_id: presetId,
			})
			const current = linkPresets.value.get(linkId) || []
			linkPresets.value.set(linkId, [...current, response.data.data])
			return { data: response.data, error: response.error }
		} catch (e) {
			return { data: null, error: await parseApiError(e) }
		}
	}

	async function removePresetFromLink(linkId: string, presetId: string) {
		try {
			const response = await api.v1.shortLinksRemovePresetFromLink(linkId, presetId)
			const current = linkPresets.value.get(linkId) || []
			linkPresets.value.set(
				linkId,
				current.filter((item) => item.preset_id !== presetId)
			)
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
		perLinkBannedUserAgents,
		isPerLinkBannedUserAgentsLoading,
		fetchPerLinkBannedUserAgents,
		createPerLinkBannedUserAgent,
		deletePerLinkBannedUserAgent,
		bannedUaPresets,
		isBannedUaPresetsLoading,
		presetPatterns,
		isPresetPatternsLoading,
		fetchBannedUaPresets,
		createBannedUaPreset,
		updateBannedUaPreset,
		deleteBannedUaPreset,
		fetchPresetPatterns,
		createPresetPattern,
		deletePresetPattern,
		linkPresets,
		isLinkPresetsLoading,
		fetchLinkPresets,
		applyPresetToLink,
		removePresetFromLink,
	}
})

async function parseApiError(error: unknown) {
	if (error instanceof Response) {
		try {
			const errorData = (await error.json()) as ErrorModel
			return (
				errorData.detail ??
				errorData.errors?.map((err) => err.message)?.join(', ') ??
				errorData.title ??
				'Unknown error'
			)
		} catch {
			return `HTTP ${error.status}: ${error.statusText || 'Unknown error'}`
		}
	}

	if (error instanceof Error) {
		return error.message
	}

	return 'Unknown error'
}

if (import.meta.hot) {
	import.meta.hot.accept(acceptHMRUpdate(useUrlShortener, import.meta.hot))
}
