import { FunctionalComponent, SVGAttributes } from 'vue';

import Dota2Logo from '@/logos/Dota2Logo.svg?component';
import FaceitLogo from '@/logos/FaceitLogo.svg?component';
import LastFmLogo from '@/logos/LastFmLogo.svg?component';
import SpotifyLogo from '@/logos/SpotifyLogo.svg?component';
import TsuwariLogo from '@/logos/TsuwariLogo.svg?component';
import VKLogo from '@/logos/VKLogo.svg?component';

type SvgComponent = FunctionalComponent<SVGAttributes>;

export const VK = VKLogo as SvgComponent;
export const Tsuwari = TsuwariLogo as SvgComponent;
export const Spotify = SpotifyLogo as SvgComponent;
export const LastFm = LastFmLogo as SvgComponent;
export const Faceit = FaceitLogo as SvgComponent;
export const Dota2 = Dota2Logo as SvgComponent;
