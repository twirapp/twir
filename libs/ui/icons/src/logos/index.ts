/// <reference types="vite-svg-loader" />

import { FunctionalComponent, SVGAttributes } from 'vue';

import Dota2 from '@/logos/Dota2Logo.svg?component';
import Faceit from '@/logos/FaceitLogo.svg?component';
import LastFm from '@/logos/LastFmLogo.svg?component';
import Spotify from '@/logos/SpotifyLogo.svg?component';
import Tsuwari from '@/logos/TsuwariLogo.svg?component';
import VK from '@/logos/VKLogo.svg?component';

export type LogoName = 'Tsuwari' | 'VK' | 'Spotify' | 'Faceit' | 'Dota2' | 'LastFm';

type Logos = {
  [K in LogoName]: FunctionalComponent<SVGAttributes>;
};

const logos: Logos = {
  Tsuwari,
  Dota2,
  Faceit,
  LastFm,
  Spotify,
  VK,
};

export default logos;
