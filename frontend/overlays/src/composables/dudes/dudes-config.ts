import type { SoundAsset, DudeAsset } from '@twirapp/dudes/types';

export const dudesSprites = [
  'dude',
  'sith',
  'agent',
  'girl',
  'cat',
];

const dudesEventSpites = [
	'santa',
];

const day = new Date().getDate();
const month = new Date().getMonth();
const isMaryChristmas =
	(month === 11 && day >= 25) ||
	(month === 0 && day <= 15);

if (isMaryChristmas) {
	dudesSprites.push(...dudesEventSpites);
}

const dudesAssetsPath = window.location.origin + '/overlays/dudes/';

export const dudesAssets: DudeAsset[] = [
  {
    alias: 'dude',
    src: dudesAssetsPath + 'sprites/dude/dude.json',
  },
  {
    alias: 'sith',
    src: dudesAssetsPath + 'sprites/sith/sith.json',
  },
  {
    alias: 'agent',
    src: dudesAssetsPath + 'sprites/agent/agent.json',
  },
  {
    alias: 'girl',
    src: dudesAssetsPath + 'sprites/girl/girl.json',
  },
  {
    alias: 'cat',
    src: dudesAssetsPath + 'sprites/cat/cat.json',
  },
  {
    alias: 'santa',
    src: dudesAssetsPath + 'sprites/santa/santa.json',
  },
];

export const dudesSounds: SoundAsset[] = [
	{
		alias: 'jump',
		src: dudesAssetsPath + 'sounds/jump.mp3',
	},
];
