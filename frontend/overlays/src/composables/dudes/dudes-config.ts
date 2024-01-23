import type { DudeAsset } from '@twirapp/dudes/types';

export const dudesSprites = [
  'dude',
  'sith',
  'agent',
  'girl',
  'cat',
  'santa',
] as const;

export type DudeSprite = typeof dudesSprites[number];

const dudesPath = window.location.origin + '/overlays/';

export const dudesAssets: DudeAsset[] = [
  {
    alias: 'dude',
    src: dudesPath + 'dudes/dude/dude.json',
  },
  {
    alias: 'sith',
    src: dudesPath + './dudes/sith/sith.json',
  },
  {
    alias: 'agent',
    src: dudesPath + './dudes/agent/agent.json',
  },
  {
    alias: 'girl',
    src: dudesPath + './dudes/girl/girl.json',
  },
  {
    alias: 'cat',
    src: dudesPath + './dudes/cat/cat.json',
  },
  {
    alias: 'santa',
    src: dudesPath + './dudes/santa/santa.json',
  },
];