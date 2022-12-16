import {
  Header,
  Grid,
  MediaQuery,
  Burger,
  ActionIcon,
  Text,
  ColorScheme,
  useMantineTheme,
  Flex,
  Group,
} from '@mantine/core';
import { useLocalStorage, useHotkeys, useMediaQuery } from '@mantine/hooks';
import { IconSun, IconMoonStars } from '@tabler/icons';
import { Dispatch, SetStateAction } from 'react';

import { Profile } from './profile';

export function NavBar({
  opened,
  setOpened,
}: {
  setOpened: Dispatch<SetStateAction<boolean>>;
  opened: boolean;
}) {
  const theme = useMantineTheme();
  const [colorScheme, setColorScheme] = useLocalStorage<ColorScheme>({
    key: 'theme',
    getInitialValueInEffect: true,
  });

  const toggleColorScheme = (value?: ColorScheme) =>
    setColorScheme(value || (colorScheme === 'dark' ? 'light' : 'dark'));

  useHotkeys([['mod+J', () => toggleColorScheme()]]);

  const largeScreen = useMediaQuery('(min-width: 250px)');

  return (
    <Header height={{ base: 50, md: 50 }} p="md">
      <Grid justify="space-between" align="center">
        <Flex gap="sm" justify="flex-start" align="center" direction="row">
          <MediaQuery largerThan="sm" styles={{ display: 'none' }} aria-label="Open navigation">
            <Burger
              opened={opened}
              onClick={() => setOpened((o) => !o)}
              size="sm"
              color={theme.colors.gray[6]}
              mr="xl"
            />
          </MediaQuery>

          <Text hidden={!largeScreen}>Tsuwari</Text>
        </Flex>
        <Group position="center">
          <ActionIcon
            variant="outline"
            color={theme.colorScheme === 'dark' ? 'yellow' : 'blue'}
            onClick={() => toggleColorScheme()}
            title="Toggle color scheme"
          >
            {theme.colorScheme === 'dark' ? <IconSun size={18} /> : <IconMoonStars size={18} />}
          </ActionIcon>

          <Profile />
        </Group>
      </Grid>
    </Header>
  );
}
