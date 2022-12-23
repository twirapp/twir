import { useLocalStorage } from '@mantine/hooks';

// Local storage key
const LOCALE_KEY = 'locale';

export const useLocale = () =>
  useLocalStorage<'en' | 'ru'>({
    key: LOCALE_KEY
  });
