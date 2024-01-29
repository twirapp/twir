import type { Settings } from '@twir/api/messages/overlays_dudes/overlays_dudes';

type DeepRequired<T> = Required<{
	[K in keyof T]: T[K] extends Required<T[K]> ? T[K] : DeepRequired<T[K]>;
}>;

export type DudesSettingsWithOptionalId = DeepRequired<Omit<Settings, 'id'> & { id?: string }>;

export const defaultDudesSettings: DudesSettingsWithOptionalId = {
	id: '',
	dudeSettings: {
		color: '#969696',
		maxLifeTime: 1000 * 60 * 30,
		gravity: 400,
		scale: 4,
		soundsEnabled: true,
		soundsVolume: 0.01,
	},
	messageBoxSettings: {
		borderRadius: 10,
		boxColor: '#eeeeee',
		fontFamily: 'Arial',
		fontSize: 20,
		padding: 10,
		showTime: 5 * 1000,
		fill: '#333333',
	},
	nameBoxSettings: {
		fontFamily: 'Arial',
		fontSize: 18,
		fill: ['#ffffff'],
		lineJoin: 'round',
		strokeThickness: 4,
		stroke: '#000000',
		fillGradientStops: [0],
		fillGradientType: 0,
		fontStyle: 'normal',
		fontVariant: 'normal',
		fontWeight: 'normal',
		dropShadow: false,
		dropShadowAlpha: 1,
		dropShadowAngle: 0,
		dropShadowBlur: 1,
		dropShadowDistance: 1,
		dropShadowColor: '#3ac7d9',
	},
};
