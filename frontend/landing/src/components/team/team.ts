type Social = 'telegram' | 'github' | 'twitch' | 'instagram';
type Developer = {
	username: string;
	avatarUrl: string;
	description: string;
	social: Partial<Record<Social, string>>;
};

export const socialLinkLabels: Record<Social, string> = {
	github: 'GitHub Profile',
	instagram: 'Instagram Profile',
	telegram: 'Telegram Profile',
	twitch: 'Twitch Channel',
};

export const developers: Developer[] = [
	{
		username: 'Satont',
		avatarUrl: 'https://avatars.githubusercontent.com/u/42675886?v=4',
		description: 'Founder and Backend, Frontend developer',
		social: {
			telegram: 'https://t.me/satont',
			github: 'https://github.com/satont',
		},
	},
	{
		username: 'MelKam',
		description: 'UI-UX Designer, Frontend developer',
		avatarUrl: 'https://avatars.githubusercontent.com/u/51422045?v=4',
		social: {
			twitch: 'https://www.twitch.tv/mellkam',
			telegram: 'https://t.me/mellkam',
			github: 'https://github.com/MellKam/',
			instagram: 'https://www.instagram.com/mel._.kam/',
		},
	},
	{
		username: 'crashmax',
		description: 'Frontend developer',
		avatarUrl: 'https://avatars.githubusercontent.com/u/15673111?v=4',
		social: {
			twitch: 'https://www.twitch.tv/vs_code',
			telegram: 'https://t.me/crashmax',
			github: 'https://github.com/crashmax-dev',
		},
	},
];
