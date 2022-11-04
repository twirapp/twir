import { persistentAtom } from '@nanostores/persistent';

export const youtubeAutoPlay = persistentAtom<boolean>('youtube-autoplay', false, {
  encode: JSON.stringify,
  decode: JSON.parse,
});
