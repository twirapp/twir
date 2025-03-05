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

interface Feature {
	title: string
	description: string
	icon: any
	fullWidth?: boolean
}

export const featuresData: Feature[] = [
	{
		title: 'Music Recognition',
		description: 'Shazam for your twitch stream! Free, without the need to connect any music services.',
		icon: MusicRecognizeSvg,
		fullWidth: true,
	},
	{
		title: 'Commands',
		description:
			'A powerful command system with aliases, counters, and all sorts of variables for all occasions',
		icon: CommandsSvg,
	},
	{
		title: 'Timers',
		description:
			'A simple system, but with verve, has become a popular announcement system from Twitch',
		icon: TimersSvg,
	},
	{
		title: 'Greetings',
		description: 'Do you want to somehow highlight your favorite viewers? Add a greeting!',
		icon: GreetingsSvg,
	},
	{
		title: 'Song requests',
		description:
			'Viewers request songs via chat commands. Bot manages queue, plays songs, and offers controls. Enhances stream with viewer engagement',
		icon: SongRequestsSvg,
	},
	{
		title: 'Keywords',
		description:
			'Identifies specified keywords in chat, triggers automated messages for engagement or information. Enhances interaction and delivers targeted content during live stream',
		icon: KeywordsSvg,
	},
	{
		title: 'Events',
		description:
			'With this powerful tool, you can easily set up customized listeners to keep track of specific events happening in the chat, or even trigger actions based on system events',
		icon: EventsSvg,
	},
	{
		title: 'Moderation',
		description: 'Create and manage chat filters to keep safe and kind communication',
		icon: ModerationSvg,
	},
	{
		title: 'OBS Websockets',
		description:
			'Highly integrate with obs studio via websockets. Change scenes, mute audio, toggle source visibility via bot',
		icon: ConnectionSvg,
	},
	{
		title: 'Stats tracking',
		description: 'Track users watch time, messages, used channel points',
		icon: StatsSvg,
	},
	{
		title: 'Overlays',
		description: 'An assortment of pre-designed overlays, including now playing, chat, emote wall, pixel dudes, and AFK displays',
		icon: OverlaysSvg,
	},
	// {
	// 	title: 'Chat alerts',
	// 	description: `If you seek streamlined chat notifications without the complexity of the entire event system, you're in the right place! Our simplified system is here to meet your needs`,
	// 	icon: ChatAlertsSvg,
	// },
	{
		title: 'Alerts',
		description: 'Want to sound alerts on rewards? We got you covered! Create custom alerts for your channel points, commands, keywords triggers',
		icon: AlertsSvg,
	},
	{
		title: 'Games',
		description: 'Looking to add a touch of fun and relaxation to the chat? No problem! We offer Russian roulette, duels, seppuku, voteban, and the magic 8-ball for your entertainment',
		icon: GamesSvg,
	},
]
