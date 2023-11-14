import type { Settings_AnimationSettings } from '@twir/grpc/generated/api/api/overlays_kappagen';

export const animations: Settings_AnimationSettings[] = [
	{
		style: 'TheCube',
		prefs: {
			size: 0.2,
			center: false,
			faces: false,
			speed: 6,
			// message: [],
		},
		enabled: true,
	} as Settings_AnimationSettings,
	{
		style: 'Text',
		prefs: {
			message: ['Twir'],
			time: 3,
		},
		enabled: true,
	},
	{
		style: 'Confetti',
		count: 150,
		enabled: true,
	},
	{
		style: 'Spiral',
		count: 150,
		enabled: true,
	},
	{
		style: 'Stampede',
		count: 150,
		enabled: true,
	},
	{
		style: 'Burst',
		count: 50,
		enabled: true,
	},
	{
		style: 'Fountain',
		count: 50,
		enabled: true,
	},
	{
		style: 'SmallPyramid',
		enabled: true,
	},
	{
		style: 'Pyramid',
		enabled: true,
	},
	{
		style: 'Fireworks',
		count: 150,
		enabled: true,
	},
	{
		style: 'Conga',
		enabled: true,
	},
];
