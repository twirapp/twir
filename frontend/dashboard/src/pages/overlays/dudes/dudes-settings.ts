import { DudesSprite } from '@twir/types/overlays';

import type { DudesOverlaySettings } from '@/gql/graphql';

export type DudesSettingsWithOptionalId = Omit<DudesOverlaySettings, 'id'> & { id?: string };

interface User {
	id: string;
	displayName: string;
}

const twitchBots: User[] = [
	{
		id: '870280719',
		displayName: 'TwirApp',
	},
	{
		id: '19264788',
		displayName: 'Nightbot',
	},
	{
		id: '1564983',
		displayName: 'Moobot',
	},
	{
		id: '105166207',
		displayName: 'Streamlabs',
	},
	{
		id: '100135110',
		displayName: 'StreamElements',
	},
	{
		id: '52268235',
		displayName: 'WizeBot',
	},
	{
		id: '496585194',
		displayName: 'DonationAlerts_',
	},
	{
		id: '237719657',
		displayName: 'Fossabot',
	},
];

export const defaultDudesSettings: DudesSettingsWithOptionalId = {
	id: '',
	ignoreSettings: {
		ignoreCommands: true,
		ignoreUsers: true,
		users: twitchBots.map(bot => bot.id),
	},
	dudeSettings: {
		maxOnScreen: 0,
		defaultSprite: DudesSprite.random,
		visibleName: true,
		color: '#969696',
		eyesColor: '#FFFFFF',
		cosmeticsColor: '#FFFFFF',
		maxLifeTime: 1000 * 60 * 30, // 30 minutes
		growTime: 1000 * 60 * 5, // 5 minutes
    growMaxScale: 10,
		gravity: 400,
		scale: 4,
		soundsEnabled: true,
		soundsVolume: 0.01,
	},
	messageBoxSettings: {
		enabled: true,
		borderRadius: 10,
		boxColor: '#EEEEEE',
		fontFamily: 'roboto',
		fontSize: 20,
		padding: 10,
		showTime: 5 * 1000,
		fill: '#333333',
	},
	nameBoxSettings: {
		fontFamily: 'roboto',
		fontSize: 18,
		fill: ['#FFFFFF'],
		lineJoin: 'round',
		strokeThickness: 4,
		stroke: '#000000',
		fillGradientStops: [0],
		fillGradientType: 0,
		fontStyle: 'normal',
		fontVariant: 'normal',
		fontWeight: 400,
		dropShadow: false,
		dropShadowAlpha: 1,
		dropShadowAngle: 0,
		dropShadowBlur: 1,
		dropShadowDistance: 1,
		dropShadowColor: '#3AC7D9',
	},
	spitterEmoteSettings: {
		enabled: true,
	},
};
