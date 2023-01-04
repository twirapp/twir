import { useMantineTheme, type ColorScheme } from '@mantine/core';
import { useHotkeys, useLocalStorage } from '@mantine/hooks';
import { useCallback } from 'react';

export const useTheme = () => {
  const theme = useMantineTheme();

  const [colorScheme, setColorScheme] = useLocalStorage<ColorScheme>({
    key: 'theme',
    getInitialValueInEffect: true,
  });

  const toggleTheme = useCallback((value?: ColorScheme) => {
    setColorScheme(value || (colorScheme === 'dark' ? 'light' : 'dark'));
  }, [colorScheme]);

  useHotkeys([['mod+J', () => toggleTheme()]]);

  return {
    theme,
    toggleTheme,
  };
};
