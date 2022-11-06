import { persistentAtom } from '@nanostores/persistent';

export const youtubeAutoPlay = persistentAtom<boolean>('youtube-autoplay', true, {
  encode: JSON.stringify,
  decode: JSON.parse,
});
