import {
  ActionIcon,
  Burger,
  ColorScheme,
  Flex,
  Grid,
  Group,
  Header,
  Loader,
  MediaQuery,
  Menu,
  Text,
  useMantineTheme,
} from '@mantine/core';
import { useHotkeys, useLocalStorage, useMediaQuery } from '@mantine/hooks';
import { IconMoonStars, IconSun } from '@tabler/icons';
import { Dispatch, SetStateAction } from 'react';

import { Profile } from './profile';

import { useProfile } from '@/services/api';
import { useLocale, LOCALES } from '@/services/dashboard';

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
  const { locale, toggleLocale } = useLocale();

  const toggleColorScheme = (value?: ColorScheme) =>
    setColorScheme(value || (colorScheme === 'dark' ? 'light' : 'dark'));

  useHotkeys([['mod+J', () => toggleColorScheme()]]);

  const largeScreen = useMediaQuery('(min-width: 300px)');

  const { data: userData, isLoading: isLoadingProfile } = useProfile();

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
          <Menu shadow="md" width={200}>
            <Menu.Target>
              <ActionIcon title="Toggle language" variant="subtle">
                {LOCALES.get(locale)?.icon}
              </ActionIcon>
            </Menu.Target>
            <Menu.Dropdown>
              <Menu.Label>Change language</Menu.Label>
              {Array.from(LOCALES.entries()).map(([locale, { icon, name }]) => (
                <Menu.Item key={locale} onClick={() => toggleLocale(locale)}>
                  {icon} {name}
                </Menu.Item>
              ))}
            </Menu.Dropdown>
          </Menu>
          {isLoadingProfile && <Loader />}
          {!isLoadingProfile && userData && <Profile user={userData} />}
        </Group>
      </Grid>
    </Header>
  );
}
