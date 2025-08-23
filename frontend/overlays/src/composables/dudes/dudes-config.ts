import { DudesSprite } from '@twir/types'
import { DudesLayers } from '@twirapp/dudes-vue'
import type { AssetsLoaderOptions, DudesTypes, SoundAsset } from '@twirapp/dudes-vue/types'

export const dudeMock = {
	id: 'Twir',
	name: 'Twir',
	color: '#8a2be2',
}

export type DudeSprite = keyof typeof DudesSprite

export function getSprite(sprite?: DudeSprite): DudesTypes.SpriteData {
	if (!sprite || sprite === DudesSprite.random) {
		const sprites = Object.values(dudesSprites)
		const spriteData = sprites[Math.floor(Math.random() * sprites.length)]
		return { ...spriteData }
	}

	const spriteData = dudesSprites[sprite]
	return { ...spriteData }
}

export const assetsLoaderOptions: AssetsLoaderOptions = {
	basePath: location.origin + '/overlays/dudes/sprites/',
	defaultSearchParams: {
		ts: Date.now(),
	},
}

export const dudesSprites: Record<Exclude<DudeSprite, 'random'>, DudesTypes.SpriteData> = {
	dude: {
		name: 'dude',
		layers: [
			{
				layer: DudesLayers.body,
				src: 'body/dude.png',
			},
			{
				layer: DudesLayers.eyes,
				src: 'eyes/dude.png',
			},
		],
	},
	agent: {
		name: 'agent',
		layers: [
			{
				layer: DudesLayers.body,
				src: 'body/dude.png',
			},
			{
				layer: DudesLayers.eyes,
				src: 'eyes/toned-glasses.png',
			},
			{
				layer: DudesLayers.cosmetics,
				src: 'cosmetics/gun.png',
			},
		],
	},
	cat: {
		name: 'cat',
		layers: [
			{
				layer: DudesLayers.body,
				src: 'body/cat.png',
			},
			{
				layer: DudesLayers.eyes,
				src: 'eyes/dude.png',
			},
			{
				layer: DudesLayers.mouth,
				src: 'mouth/cat.png',
			},
		],
	},
	girl: {
		name: 'girl',
		layers: [
			{
				layer: DudesLayers.body,
				src: 'body/dude.png',
			},
			{
				layer: DudesLayers.eyes,
				src: 'eyes/dude.png',
			},
			{
				layer: DudesLayers.hat,
				src: 'hat/girl-ribbon.png',
			},
		],
	},
	santa: {
		name: 'santa',
		layers: [
			{
				layer: DudesLayers.body,
				src: 'body/dude.png',
			},
			{
				layer: DudesLayers.eyes,
				src: 'eyes/dude.png',
			},
			{
				layer: DudesLayers.hat,
				src: 'hat/santa.png',
			},
		],
	},
	sith: {
		name: 'sith',
		layers: [
			{
				layer: DudesLayers.body,
				src: 'body/devil.png',
			},
			{
				layer: DudesLayers.eyes,
				src: 'eyes/dude.png',
			},
			{
				layer: DudesLayers.cosmetics,
				src: 'cosmetics/lightsaber.png',
			},
		],
	},
}

export const dudesSounds: SoundAsset[] = [
	{
		alias: 'Jump',
		src: location.origin + '/overlays/dudes/sounds/jump.mp3',
	},
]
