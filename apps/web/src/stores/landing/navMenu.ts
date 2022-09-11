import { atom } from 'nanostores';

import type { NavMenuLocale } from '@/types/navMenu.js';

export const navMenuLocaleStore = atom<NavMenuLocale[]>([]);
