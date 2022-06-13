import { persistentAtom } from '@nanostores/persistent';

export const localeStore = persistentAtom<string>('locale', 'en');