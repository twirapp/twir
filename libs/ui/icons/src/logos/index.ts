/// <reference types="vite-svg-loader" />

import { FunctionalComponent, SVGAttributes } from 'vue';

import Dota2 from '@/logos/Dota2Logo.svg?component';
import Faceit from '@/logos/FaceitLogo.svg?component';
import LastFm from '@/logos/LastFmLogo.svg?component';
import Spotify from '@/logos/SpotifyLogo.svg?component';
import Tsuwari from '@/logos/TsuwariLogo.svg?component';
import VK from '@/logos/VKLogo.svg?component';

const logos = {
  Tsuwari,
  Dota2,
  Faceit,
  LastFm,
  Spotify,
  VK,
};

export type LogoName = keyof typeof logos;

export default logos as {
  readonly [K in LogoName]: FunctionalComponent<SVGAttributes>;
};
