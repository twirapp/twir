import Discord from 'assets/icons/social/discord.svg?component'
import Donatello from 'assets/icons/social/donatello.svg?component'
import Donatepay from 'assets/icons/social/donatepay.svg?component'
import Donatestream from 'assets/icons/social/donatestream.svg?component'
import DonationAlerts from 'assets/icons/social/donationalerts.svg?component'
import Faceit from 'assets/icons/social/faceit.svg?component'
import LastFm from 'assets/icons/social/lastfm.svg?component'
import OBS from 'assets/icons/social/obs.svg?component'
import SevenTv from 'assets/icons/social/seventv.svg?component'
import Spotify from 'assets/icons/social/spotify.svg?component'
import Streamlabs from 'assets/icons/social/streamlabs.svg?component'
import Twitch from 'assets/icons/social/twitch.svg?component'
import Valorant from 'assets/icons/social/valorant.svg?component'
import Vk from 'assets/icons/social/vk.svg?component'

interface Integration {
	icon: any
	label: string
	href: string
	width?: string
}

export const integrationsData: Integration[] = [

	{
		icon: Twitch,
		label: 'Twitch',
		href: 'https://twitch.tv',
	},
	{
		icon: SevenTv,
		label: '7TV',
		href: 'https://7tv.app',
	},
	{
		icon: OBS,
		label: 'OBS',
		href: 'https://obsproject.com',
	},
	{
		icon: DonationAlerts,
		label: 'DonationAlerts',
		href: 'https://donationalerts.com',
	},
	{
		icon: Donatello,
		label: 'Donatello',
		href: 'https://donatello.to',
	},
	{
		icon: Donatepay,
		label: 'DonatePay',
		href: 'https://donatepay.ru',
		width: '120px',
	},
	{
		icon: Donatestream,
		label: 'Donate.stream',
		href: 'https://donate.stream',
		width: '120px',
	},
	{
		icon: Discord,
		label: 'Discord',
		href: 'https://discord.com',
	},
	// {
	// 	icon: StreamElements,
	// 	label: 'Stream Elements',
	// 	href: 'https://streamelements.com',
	// },
	{
		icon: Streamlabs,
		label: 'Streamlabs',
		href: 'https://streamlabs.com',
	},
	{
		icon: Faceit,
		label: 'FaceIt',
		href: 'https://faceit.com',
	},
	{
		icon: Valorant,
		label: 'Valorant',
		href: 'https://playvalorant.com/en-us/',
	},
	{
		icon: Vk,
		label: 'VK',
		href: 'https://vk.com',
	},
	{
		icon: Spotify,
		label: 'Spotify',
		href: 'https://spotify.com',
	},
	{
		icon: LastFm,
		label: 'LastFM',
		href: 'https://last.fm',
	},
]
