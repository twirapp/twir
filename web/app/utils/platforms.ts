import { Platform } from '~/gql/graphql.js'

/**
 * Single source of truth for Platform (gql enum) presentation:
 * labels, simple-icons names and brand colors.
 *
 * Always use these mappings instead of hardcoding platform labels,
 * icons or hex colors in components.
 */
export interface PlatformMeta {
	platform: Platform
	/** Human-readable platform name, e.g. "VK Video Live" */
	label: string
	/** simple-icons brand icon for the Nuxt <Icon /> component */
	icon: string
	/** Brand color as hex, e.g. "#9146FF" */
	color: string
	/** Tailwind class painting text/icon in the brand color */
	colorClass: string
	/** Classes for an outline Badge in brand colors */
	badgeClass: string
	/** Classes applied to an active chip in PlatformSelector */
	activeClass: string
}

export const PLATFORM_META = {
	[Platform.Twitch]: {
		platform: Platform.Twitch,
		label: 'Twitch',
		icon: 'simple-icons:twitch',
		color: '#9146FF',
		colorClass: 'text-[#9146FF]',
		badgeClass: 'border-[#9146FF]/30 bg-[#9146FF]/10 text-[#7C3AED]',
		activeClass:
			'data-[active=true]:border-[#9146FF] data-[active=true]:bg-[#9146FF]/10 data-[active=true]:text-[#9146FF]',
	},
	[Platform.Kick]: {
		platform: Platform.Kick,
		label: 'Kick',
		icon: 'simple-icons:kick',
		color: '#53FC18',
		colorClass: 'text-[#53FC18]',
		badgeClass: 'border-[#53FC18]/30 bg-[#53FC18]/10 text-[#3CB30F]',
		activeClass:
			'data-[active=true]:border-[#53FC18] data-[active=true]:bg-[#53FC18]/10 data-[active=true]:text-[#53FC18]',
	},
	[Platform.VkVideoLive]: {
		platform: Platform.VkVideoLive,
		label: 'VK Video Live',
		icon: 'simple-icons:vk',
		color: '#0077FF',
		colorClass: 'text-[#0077FF]',
		badgeClass: 'border-[#0077FF]/30 bg-[#0077FF]/10 text-[#0077FF]',
		activeClass:
			'data-[active=true]:border-[#0077FF] data-[active=true]:bg-[#0077FF]/10 data-[active=true]:text-[#0077FF]',
	},
	[Platform.Youtube]: {
		platform: Platform.Youtube,
		label: 'YouTube',
		icon: 'simple-icons:youtube',
		color: '#FF0000',
		colorClass: 'text-[#FF0000]',
		badgeClass: 'border-[#FF0000]/30 bg-[#FF0000]/10 text-[#FF0000]',
		activeClass:
			'data-[active=true]:border-[#FF0000] data-[active=true]:bg-[#FF0000]/10 data-[active=true]:text-[#FF0000]',
	},
} satisfies Record<Platform, PlatformMeta>

/** All platforms in stable display order */
export const PLATFORMS = [
	Platform.Twitch,
	Platform.Kick,
	Platform.VkVideoLive,
	Platform.Youtube,
] as const satisfies readonly Platform[]

/** Meta list in display order — handy for v-for option lists */
export const PLATFORM_OPTIONS: readonly PlatformMeta[] = PLATFORMS.map(
	(platform) => PLATFORM_META[platform],
)

/** Meta keyed by lowercase slug ("twitch", "kick", "vk_video_live", "youtube") for non-gql string platforms */
export const PLATFORM_META_BY_SLUG: Readonly<Record<string, PlatformMeta>> = Object.fromEntries(
	PLATFORM_OPTIONS.map((option) => [option.platform.toLowerCase(), option]),
)

export function getPlatformMeta(platform: Platform): PlatformMeta {
	return PLATFORM_META[platform]
}

/** Lookup by lowercase slug ("twitch", "vk_video_live", ...) as returned by non-gql APIs; null when unknown. */
export function getPlatformMetaBySlug(slug: string): PlatformMeta | null {
	return PLATFORM_META_BY_SLUG[slug.toLowerCase()] ?? null
}

/** Accepts raw strings too: some gql types still expose platforms as plain strings; unknown values pass through. */
export function getPlatformLabel(platform: Platform | string): string {
	const meta = PLATFORM_OPTIONS.find((option) => option.platform === platform)
	return meta?.label ?? platform
}
