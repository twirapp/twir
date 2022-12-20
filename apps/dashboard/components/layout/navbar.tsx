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
  Menu,
} from '@mantine/core';
import { useLocalStorage, useHotkeys, useMediaQuery } from '@mantine/hooks';
import { IconSun, IconMoonStars } from '@tabler/icons';
import { US, RU } from 'country-flag-icons/react/3x2';
import { useRouter } from 'next/router';
import { Dispatch, SetStateAction } from 'react';

import { Profile } from './profile';

import { useProfile } from '@/services/api';
import { useLocale } from '@/services/dashboard';

const flags = {
  en: <US style={{ height: 14 }} />,
  ru: <RU style={{ height: 14 }} />,
};

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
  const [locale, setLocale] = useLocale();

  const toggleLanguage = (newLocale: 'en' | 'ru') => {
    setLocale(newLocale);
  };

  const toggleColorScheme = (value?: ColorScheme) =>
    setColorScheme(value || (colorScheme === 'dark' ? 'light' : 'dark'));

  useHotkeys([['mod+J', () => toggleColorScheme()]]);

  const largeScreen = useMediaQuery('(min-width: 250px)');

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
                {flags[locale]}
              </ActionIcon>
            </Menu.Target>
            <Menu.Dropdown>
              <Menu.Label>Change language</Menu.Label>
              <Menu.Item onClick={() => toggleLanguage('en')}>{flags['en']} English</Menu.Item>
              <Menu.Item onClick={() => toggleLanguage('ru')}>{flags['ru']} Russian</Menu.Item>
            </Menu.Dropdown>
          </Menu>

          {isLoadingProfile && <Loader />}
          {!isLoadingProfile && userData && <Profile user={userData} />}
        </Group>
      </Grid>
    </Header>
  );
}
