import { useMantineTheme, useMantineColorScheme } from '@mantine/core';
import { useHotkeys } from '@mantine/hooks';

export const useTheme = () => {
  const theme = useMantineTheme();
  const { colorScheme, toggleColorScheme } = useMantineColorScheme();

  useHotkeys([['mod+J', () => toggleColorScheme()]]);

  return {
    theme,
    colorScheme,
    toggleColorScheme,
  };
};
