import Discord from '@/assets/social/discord.svg?use';
import Donatello from '@/assets/social/donatello.svg?use';
import Donatepay from '@/assets/social/donatepay.svg?use';
import Donatestream from '@/assets/social/donatestream.svg?use';
import DonationAlerts from '@/assets/social/donationalerts.svg?use';
import Faceit from '@/assets/social/faceit.svg?use';
import LastFm from '@/assets/social/lastfm.svg?use';
import OBS from '@/assets/social/obs.svg?use';
import SevenTv from '@/assets/social/seventv.svg?use';
import Spotify from '@/assets/social/spotify.svg?use';
// import StreamElements from '@/assets/social/streamelements.svg?use';
import Streamlabs from '@/assets/social/streamlabs.svg?use';
import Twitch from '@/assets/social/twitch.svg?use';
import Valorant from '@/assets/social/valorant.svg?use';
import Vk from '@/assets/social/vk.svg?use';

interface Integration {
	icon: any;
	label: string;
	href: string;
	width?: string;
}

export const integrations: Integration[] = [

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
];
