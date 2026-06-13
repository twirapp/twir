import type { Component } from 'vue'

import { IconCalendarCog } from '@tabler/icons-vue'
import { DISCORD_INVITE_URL, GITHUB_REPOSITORY_URL } from '@twir/brand'
import {
	AudioLines,
	Bell,
	Blend,
	Box,
	ClipboardPenLine,
	ComponentIcon,
	Dices,
	GemIcon,
	GiftIcon,
	Globe,
	Import,
	LayoutDashboard,
	LinkIcon,
	MessageCircleHeart,
	MessageCircleWarning,
	Package,
	PackageCheck,
	PackagePlus,
	ScrollTextIcon,
	SettingsIcon,
	Shield,
	ShieldUser,
	Smile,
	SparklesIcon,
	Timer,
	Users,
	Variable,
	WholeWord,
} from 'lucide-vue-next'

export interface NavigationItem {
	name?: string
	icon: Component
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
		icon: LayoutDashboard,
		path: '/dashboard',
	},
	{
		name: 'Bot Settings',
		icon: SettingsIcon,
		path: '/dashboard/bot-settings',
	},
	{
		name: 'Modules',
		icon: ComponentIcon,
		path: '/dashboard/modules',
		isNew: true,
	},
	{
		translationKey: 'sidebar.integrations',
		icon: Box,
		path: '/dashboard/integrations',
	},
	{
		translationKey: 'sidebar.alerts',
		icon: Bell,
		path: '/dashboard/alerts',
	},
	{
		translationKey: 'sidebar.chatAlerts',
		icon: MessageCircleWarning,
		path: '/dashboard/events/chat-alerts',
	},
	{
		translationKey: 'sidebar.events',
		icon: IconCalendarCog,
		path: '/dashboard/events',
	},
	{
		translationKey: 'sidebar.overlays',
		icon: Blend,
		path: '/dashboard/overlays',
		isNew: true,
	},
	{
		translationKey: 'sidebar.songRequests',
		icon: AudioLines,
		path: '/dashboard/song-requests',
	},
	{
		translationKey: 'sidebar.games',
		icon: Dices,
		path: '/dashboard/games',
	},
	{
		translationKey: 'sidebar.commands.label',
		icon: Package,
		path: '/dashboard/commands',
		openStateKey: 'commands',
		child: [
			{
				translationKey: 'sidebar.commands.custom',
				icon: PackagePlus,
				path: '/dashboard/commands/custom',
			},
			{
				translationKey: 'sidebar.commands.builtin',
				icon: PackageCheck,
				path: '/dashboard/commands/builtin',
			},
		],
	},
	{
		translationKey: 'sidebar.community',
		icon: Users,
		path: '/dashboard/community',
		openStateKey: 'community',
		child: [
			{
				name: 'Chat Logs',
				icon: ScrollTextIcon,
				path: '/dashboard/community?tab=chat-logs',
			},
			{
				translationKey: 'community.users.title',
				icon: Users,
				path: '/dashboard/community?tab=users',
			},
			{
				translationKey: 'sidebar.roles',
				icon: ShieldUser,
				path: '/dashboard/community?tab=permissions',
			},
			{
				translationKey: 'community.emotesStatistic.title',
				icon: Smile,
				path: '/dashboard/community?tab=emotes-stats',
			},
			{
				name: 'Rewards history',
				icon: SparklesIcon,
				path: '/dashboard/community?tab=rewards-history',
			},
		],
	},
	{
		translationKey: 'sidebar.moderation',
		icon: Shield,
		path: '/dashboard/moderation',
	},
	{
		translationKey: 'sidebar.timers',
		icon: Timer,
		path: '/dashboard/timers',
	},
	{
		translationKey: 'sidebar.giveaways',
		icon: GiftIcon,
		path: '/dashboard/giveaways',
		isNew: true,
	},
	{
		translationKey: 'sidebar.keywords',
		icon: WholeWord,
		path: '/dashboard/keywords',
	},
	{
		translationKey: 'sidebar.variables',
		icon: Variable,
		path: '/dashboard/variables',
	},
	{
		translationKey: 'sidebar.greetings',
		icon: MessageCircleHeart,
		path: '/dashboard/greetings',
	},
	{
		name: 'Expiring Vips',
		icon: GemIcon,
		path: '/dashboard/expiring-vips',
		isNew: true,
	},
	{
		translationKey: 'sidebar.import',
		icon: Import,
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
		icon: Bell,
		href: '/dashboard/notifications',
		showNotificationsBadge: true,
	},
	{
		translationKey: 'sidebar.publicPage',
		icon: Globe,
		href: '', // Will be computed dynamically
		isExternal: true,
		isPublicPageDependent: true,
	},
	{
		name: 'URL Shortener',
		icon: LinkIcon,
		href: '/url-shortener',
		isExternal: true,
		isPublicPageDependent: true,
	},
	{
		name: 'Hastebin',
		icon: ClipboardPenLine,
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
