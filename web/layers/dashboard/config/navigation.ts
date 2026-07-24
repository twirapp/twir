import { DISCORD_INVITE_URL, GITHUB_REPOSITORY_URL } from '@twir/brand'

export interface NavigationItem {
	name?: string
	icon: string
	path: string
	disabled?: boolean
	isNew?: boolean
	openStateKey?: string
	child?: NavigationItem[]
	translationKey?: string // i18n key for translation
}

export interface NavigationConfig {
	items: NavigationItem[]
}

// Base navigation items without translations
// Name will be populated from i18n in components
export const baseNavigationItems: Array<Partial<NavigationItem>> = [
	{
		translationKey: 'sidebar.dashboard',
		icon: 'lucide:layout-dashboard',
		path: '/dashboard',
	},
	{
		name: 'Platforms',
		icon: 'lucide:layers',
		path: '/dashboard/platforms',
	},
	{
		name: 'Modules',
		icon: 'lucide:component',
		path: '/dashboard/modules',
		isNew: true,
	},
	{
		translationKey: 'sidebar.integrations',
		icon: 'lucide:box',
		path: '/dashboard/integrations',
	},
	{
		translationKey: 'sidebar.alerts',
		icon: 'lucide:bell',
		path: '/dashboard/alerts',
	},
	{
		translationKey: 'sidebar.chatAlerts',
		icon: 'lucide:message-circle-warning',
		path: '/dashboard/events/chat-alerts',
	},
	{
		translationKey: 'sidebar.events',
		icon: 'tabler:calendar-cog',
		path: '/dashboard/events',
	},
	{
		translationKey: 'sidebar.overlays',
		icon: 'lucide:blend',
		path: '/dashboard/overlays',
		isNew: true,
	},
	{
		translationKey: 'sidebar.songRequests',
		icon: 'lucide:audio-lines',
		path: '/dashboard/song-requests',
	},
	{
		translationKey: 'sidebar.games',
		icon: 'lucide:dices',
		path: '/dashboard/games',
	},
	{
		translationKey: 'sidebar.commands.label',
		icon: 'lucide:package',
		path: '/dashboard/commands',
		openStateKey: 'commands',
		child: [
			{
				translationKey: 'sidebar.commands.custom',
				icon: 'lucide:package-plus',
				path: '/dashboard/commands/custom',
			},
			{
				translationKey: 'sidebar.commands.builtin',
				icon: 'lucide:package-check',
				path: '/dashboard/commands/builtin',
			},
		],
	},
	{
		translationKey: 'sidebar.community',
		icon: 'lucide:users',
		path: '/dashboard/community',
		openStateKey: 'community',
		child: [
			{
				name: 'Chat Logs',
				icon: 'lucide:scroll-text',
				path: '/dashboard/community?tab=chat-logs',
			},
			{
				translationKey: 'community.users.title',
				icon: 'lucide:users',
				path: '/dashboard/community?tab=users',
			},
			{
				translationKey: 'sidebar.roles',
				icon: 'lucide:shield-user',
				path: '/dashboard/community?tab=permissions',
			},
			{
				translationKey: 'community.emotesStatistic.title',
				icon: 'lucide:smile',
				path: '/dashboard/community?tab=emotes-stats',
			},
			{
				name: 'Rewards history',
				icon: 'lucide:sparkles',
				path: '/dashboard/community?tab=rewards-history',
			},
		],
	},
	{
		translationKey: 'sidebar.moderation',
		icon: 'lucide:shield',
		path: '/dashboard/moderation',
	},
	{
		translationKey: 'sidebar.timers',
		icon: 'lucide:timer',
		path: '/dashboard/timers',
	},
	{
		translationKey: 'sidebar.giveaways',
		icon: 'lucide:gift',
		path: '/dashboard/giveaways',
		isNew: true,
	},
	{
		translationKey: 'sidebar.keywords',
		icon: 'lucide:whole-word',
		path: '/dashboard/keywords',
	},
	{
		translationKey: 'sidebar.variables',
		icon: 'lucide:variable',
		path: '/dashboard/variables',
	},
	{
		translationKey: 'sidebar.greetings',
		icon: 'lucide:message-circle-heart',
		path: '/dashboard/greetings',
	},
	{
		name: 'Expiring Vips',
		icon: 'lucide:gem',
		path: '/dashboard/expiring-vips',
		isNew: true,
	},
	{
		translationKey: 'sidebar.import',
		icon: 'lucide:import',
		path: '/dashboard/import',
	},
]

export interface FooterNavigationItem {
	name?: string
	icon: Component | string // Component or SVG asset path
	href: string
	translationKey?: string
	isExternal?: boolean
	showNotificationsBadge?: boolean
	isPublicPageDependent?: boolean // Shows only if public page is available
	computedHref?: () => string // For dynamic hrefs like hastebin
}

export const footerNavigationItems: FooterNavigationItem[] = [
	{
		name: 'Discord',
		icon: 'discord', // Special case: will use DiscordLogo in sidebar
		href: DISCORD_INVITE_URL,
		isExternal: true,
	},
	{
		name: 'GitHub',
		icon: 'github', // Special case: will use GithubLogo in sidebar
		href: GITHUB_REPOSITORY_URL,
		isExternal: true,
	},
	{
		translationKey: 'sidebar.notifications',
		icon: 'lucide:bell',
		href: '/dashboard/notifications',
		showNotificationsBadge: true,
	},
	{
		translationKey: 'sidebar.publicPage',
		icon: 'lucide:globe',
		href: '', // Will be computed dynamically
		isExternal: true,
		isPublicPageDependent: true,
	},
	{
		name: 'URL Shortener',
		icon: 'lucide:link',
		href: '/url-shortener',
		isExternal: true,
		isPublicPageDependent: true,
	},
	{
		name: 'Hastebin',
		icon: 'lucide:clipboard-pen-line',
		href: '/h',
		isExternal: true,
		isPublicPageDependent: true,
	},
]

// Helper to get flattened routes for command menu (without sub-items)
export function getFlatNavigationItems() {
	const result: Array<{
		name?: string
		translationKey?: string | string[]
		icon: Component
		path: string
		disabled?: boolean
		isNew?: boolean
	}> = []

	for (const item of baseNavigationItems) {
		if (!item.icon || !item.path) continue

		if (item.child) {
			// Add parent
			// result.push({
			// 	name: item.name,
			// 	translationKey: item.translationKey,
			// 	icon: item.icon,
			// 	path: item.path,
			// 	disabled: item.disabled,
			// 	isNew: item.isNew,
			// })
			// Add children
			for (const child of item.child) {
				if (!child.icon || !child.path) continue
				result.push({
					name: child.name,
					translationKey: [item.translationKey ?? '', child.translationKey ?? ''].filter(Boolean),
					icon: child.icon,
					path: child.path,
					disabled: child.disabled,
					isNew: child.isNew,
				})
			}
		} else {
			result.push({
				name: item.name,
				translationKey: item.translationKey,
				icon: item.icon,
				path: item.path,
				disabled: item.disabled,
				isNew: item.isNew,
			})
		}
	}

	return result
}
