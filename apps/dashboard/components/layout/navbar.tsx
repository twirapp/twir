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
  Loader,
} from '@mantine/core';
import { useLocalStorage, useHotkeys, useMediaQuery } from '@mantine/hooks';
import { IconSun, IconMoonStars } from '@tabler/icons';
import { AuthUser } from '@tsuwari/shared';
import { Dispatch, SetStateAction } from 'react';
import useSWR from 'swr';

import { swrFetcher } from '../../services/swrFetcher';
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

  const { data: userData, isLoading: isLoadingProfile } = useSWR<AuthUser>(
    '/api/auth/profile',
    swrFetcher,
  );

  return (
    <Header height={{ base: 50, md: 50 }} p="md">
      <Grid justify="space-between" align="center">
        <Flex gap="sm" justify="flex-start" align="center" direction="row">
          <MediaQuery largerThan="sm" styles={{ display: 'none' }} aria-label="Open navigation">
            <Burger
              opened={opened}
              onClick={() => setOpened(!opened)}
              size="sm"
              color={theme.colors.gray[6]}
              mr="xl"
            />
          </MediaQuery>

          <Text hidden={!largeScreen}>Tsuwari</Text>
        </Flex>
        <Group position="center">
          <ActionIcon
            variant="subtle"
            color={theme.colorScheme === 'dark' ? 'yellow' : 'blue'}
            onClick={() => toggleColorScheme()}
            title="Toggle color scheme"
          >
            {theme.colorScheme === 'dark' ? <IconSun size={18} /> : <IconMoonStars size={18} />}
          </ActionIcon>

          {isLoadingProfile && <Loader />}
          {!isLoadingProfile && userData && <Profile user={userData} />}
        </Group>
      </Grid>
    </Header>
  );
}
