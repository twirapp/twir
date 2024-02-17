import type { SoundAsset, DudeAsset, AssetsLoadOptions } from '@twirapp/dudes/types';

export const dudesTwir = 'Twir';

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

export const assetsLoadOptions: AssetsLoadOptions = {
  basePath: location.origin + '/overlays/dudes/sprites/',
  defaultSearchParams: {
    ts: Date.now(),
  },
};

export const dudesAssets: DudeAsset[] = [
  {
    alias: 'dude',
    src: 'dude/dude.json',
  },
  {
    alias: 'sith',
    src: 'sith/sith.json',
  },
  {
    alias: 'agent',
    src: 'agent/agent.json',
  },
  {
    alias: 'girl',
    src: 'girl/girl.json',
  },
  {
    alias: 'cat',
    src: 'cat/cat.json',
  },
  {
    alias: 'santa',
    src: 'santa/santa.json',
  },
];

export const dudesSounds: SoundAsset[] = [
	{
		alias: 'jump',
		src: location.origin + '/overlays/dudes/sounds/jump.mp3',
	},
];
