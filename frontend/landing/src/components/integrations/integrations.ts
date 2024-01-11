import DonationAlerts from '@/assets/social/donationalerts.svg?use';
import Faceit from '@/assets/social/faceit.svg?use';
import LastFm from '@/assets/social/lastfm.svg?use';
import OBS from '@/assets/social/obs.svg?use';
import Spotify from '@/assets/social/spotify.svg?use';
// import StreamElements from '@/assets/social/streamelements.svg?use';
import Streamlabs from '@/assets/social/streamlabs.svg?use';
import Twitch from '@/assets/social/twitch.svg?use';
import Vk from '@/assets/social/vk.svg?use';

interface Integration {
	icon: any;
	label: string;
	href: string;
}

export const integrations: Integration[] = [
  {
		icon: DonationAlerts,
		label: 'DonationAlerts',
		href: 'https://donationalerts.com',
	},
  {
		icon: Faceit,
		label: 'FaceIt',
		href: 'https://faceit.com',
	},
  {
		icon: OBS,
		label: 'OBS',
		href: 'https://obsproject.com',
	},
  {
		icon: Spotify,
		label: 'Spotify',
		href: 'https://spotify.com',
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
		icon: Twitch,
		label: 'Twitch',
		href: 'https://twitch.tv',
	},
  {
		icon: Vk,
		label: 'VK',
		href: 'https://vk.com',
	},
  {
		icon: LastFm,
		label: 'LastFM',
		href: 'https://last.fm',
	},
];
