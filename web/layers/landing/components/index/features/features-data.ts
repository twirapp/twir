import AlertsSvg from 'assets/icons/features/alerts.svg?component'
import CommandsSvg from 'assets/icons/features/commands.svg?component'
import ConnectionSvg from 'assets/icons/features/connection.svg?component'
import EventsSvg from 'assets/icons/features/events.svg?component'
import GamesSvg from 'assets/icons/features/games.svg?component'
import GreetingsSvg from 'assets/icons/features/greetings.svg?component'
import KeywordsSvg from 'assets/icons/features/keywords.svg?component'
import ModerationSvg from 'assets/icons/features/moderation.svg?component'
import MusicRecognizeSvg from 'assets/icons/features/music-recognize.svg?component'
import OverlaysSvg from 'assets/icons/features/overlays.svg?component'
import SongRequestsSvg from 'assets/icons/features/song-requests.svg?component'
import StatsSvg from 'assets/icons/features/stats.svg?component'
import TimersSvg from 'assets/icons/features/timers.svg?component'

import { Icon } from '#components'

interface Feature {
	id: string
	title: string
	description: string
	icon: any
	fullWidth?: boolean
}

export const featuresData: Feature[] = [
	{
		id: 'music-recognition',
		title: 'musicRecognition',
		description: 'musicRecognition',
		icon: MusicRecognizeSvg,
	},
	{
		id: 'vips',
		title: 'vips',
		description: 'vips',
		icon: h(Icon, { name: 'lucide:gem' }),
	},
	{
		id: 'giveaways',
		title: 'giveaways',
		description: 'giveaways',
		icon: h(Icon, { name: 'lucide:gift' }),
	},
	{
		id: 'commands',
		title: 'commands',
		description: 'commands',
		icon: CommandsSvg,
	},
	{
		id: 'variables',
		title: 'variables',
		description: 'variables',
		icon: h(Icon, { name: 'lucide:code' }),
	},
	{
		id: 'mcp',
		title: 'mcp',
		description: 'mcp',
		icon: h(Icon, { name: 'lucide:bot' }),
	},
	{
		id: 'timers',
		title: 'timers',
		description: 'timers',
		icon: TimersSvg,
	},
	{
		id: 'greetings',
		title: 'greetings',
		description: 'greetings',
		icon: GreetingsSvg,
	},
	{
		id: 'song-requests',
		title: 'songRequests',
		description: 'songRequests',
		icon: SongRequestsSvg,
	},
	{
		id: 'keywords',
		title: 'keywords',
		description: 'keywords',
		icon: KeywordsSvg,
	},
	{
		id: 'events',
		title: 'events',
		description: 'events',
		icon: EventsSvg,
	},
	{
		id: 'moderation',
		title: 'moderation',
		description: 'moderation',
		icon: ModerationSvg,
	},
	{
		id: 'obs',
		title: 'obs',
		description: 'obs',
		icon: ConnectionSvg,
	},
	{
		id: 'stats',
		title: 'stats',
		description: 'stats',
		icon: StatsSvg,
	},
	{
		id: 'overlays',
		title: 'overlays',
		description: 'overlays',
		icon: OverlaysSvg,
	},
	// {
	// 	title: 'Chat alerts',
	// 	description: `If you seek streamlined chat notifications without the complexity of the entire event system, you're in the right place! Our simplified system is here to meet your needs`,
	// 	icon: ChatAlertsSvg,
	// },
	{
		id: 'alerts',
		title: 'alerts',
		description: 'alerts',
		icon: AlertsSvg,
	},
	{
		id: 'games',
		title: 'games',
		description: 'games',
		icon: GamesSvg,
	},
	{
		id: 'short-urls',
		title: 'shortUrls',
		description: 'shortUrls',
		icon: h(Icon, { name: 'lucide:link' }),
	},
	{
		id: 'hastebins',
		title: 'hastebins',
		description: 'hastebins',
		icon: h(Icon, { name: 'lucide:file-scan' }),
	},
]
