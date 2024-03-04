import { DudesSprite } from '@twir/types/overlays';
import { DudesLayers } from '@twirapp/dudes';
import type { SoundAsset, AssetsLoaderOptions, DudesTypes } from '@twirapp/dudes/types';

export const dudesTwir = 'Twir';

export type DudeSprite = keyof typeof DudesSprite

export function getSprite(sprite?: DudeSprite): DudesTypes.SpriteData {
	if (!sprite || sprite === 'random') {
		const sprites = Object.values(dudesSprites);
		return sprites[Math.floor(Math.random() * sprites.length)];
	}

	return dudesSprites[sprite];
}

export const assetsLoaderOptions: AssetsLoaderOptions = {
  basePath: location.origin + '/overlays/dudes/sprites/',
  defaultSearchParams: {
    ts: Date.now(),
  },
};

export const dudesSprites: Record<
	Exclude<DudeSprite, 'random'>,
	DudesTypes.SpriteData
> = {
	dude: {
		name: 'dude',
		layers: [
			{
				layer: DudesLayers.Body,
				src: 'body/dude.png',
			},
			{
				layer: DudesLayers.Eyes,
				src: 'eyes/dude.png',
			},
		],
	},
	agent: {
		name: 'agent',
		layers: [
			{
				layer: DudesLayers.Body,
				src: 'body/dude.png',
			},
			{
				layer: DudesLayers.Eyes,
				src: 'eyes/toned-glasses.png',
			},
			{
				layer: DudesLayers.Cosmetics,
				src: 'cosmetics/gun.png',
			},
		],
	},
	cat: {
		name: 'cat',
		layers: [
			{
				layer: DudesLayers.Body,
				src: 'body/cat.png',
			},
			{
				layer: DudesLayers.Eyes,
				src: 'eyes/dude.png',
			},
			{
				layer: DudesLayers.Mouth,
				src: 'mouth/cat.png',
			},
		],
	},
	girl: {
		name: 'girl',
		layers: [
			{
				layer: DudesLayers.Body,
				src: 'body/dude.png',
			},
			{
				layer: DudesLayers.Eyes,
				src: 'eyes/dude.png',
			},
			{
				layer: DudesLayers.Hat,
				src: 'hat/girl-ribbon.png',
			},
		],
	},
	santa: {
		name: 'santa',
		layers: [
			{
				layer: DudesLayers.Body,
				src: 'body/dude.png',
			},
			{
				layer: DudesLayers.Eyes,
				src: 'eyes/dude.png',
			},
			{
				layer: DudesLayers.Hat,
				src: 'hat/santa.png',
			},
		],
	},
	sith: {
		name: 'sith',
		layers: [
			{
				layer: DudesLayers.Body,
				src: 'body/devil.png',
			},
			{
				layer: DudesLayers.Eyes,
				src: 'eyes/dude.png',
			},
			{
				layer: DudesLayers.Cosmetics,
				src: 'cosmetics/lightsaber.png',
			},
		],
	},
};

export const dudesSounds: SoundAsset[] = [
	{
		alias: 'Jump',
		src: location.origin + '/overlays/dudes/sounds/jump.mp3',
	},
];
