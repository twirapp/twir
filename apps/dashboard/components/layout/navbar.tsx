import {
  ActionIcon,
  Burger,
  ColorScheme,
  Container,
  createStyles,
  Flex,
  Group,
  Header,
  Loader,
  Menu,
  Text,
  useMantineTheme,
} from '@mantine/core';
import { useHotkeys, useLocalStorage } from '@mantine/hooks';
import { IconMoonStars, IconSun, IconLanguage } from '@tabler/icons';
import { Dispatch, SetStateAction } from 'react';

import { Profile } from './profile';

import { useProfile } from '@/services/api';
import { useLocale, LOCALES } from '@/services/dashboard';

const useStyles = createStyles((theme) => ({
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    height: '100%',
  },

  hiddenMobile: {
    [theme.fn.smallerThan('sm')]: {
      display: 'none',
    },
  },

  hiddenDesktop: {
    [theme.fn.largerThan('sm')]: {
      display: 'none',
    },
  },
}));

export function NavBar({
  opened,
  setOpened,
}: {
  setOpened: Dispatch<SetStateAction<boolean>>;
  opened: boolean;
}) {
  const { classes } = useStyles();
  const theme = useMantineTheme();
  const [colorScheme, setColorScheme] = useLocalStorage<ColorScheme>({
    key: 'theme',
    getInitialValueInEffect: true,
  });
  const { locale, toggleLocale } = useLocale();

  const toggleColorScheme = (value?: ColorScheme) =>
    setColorScheme(value || (colorScheme === 'dark' ? 'light' : 'dark'));

  useHotkeys([['mod+J', () => toggleColorScheme()]]);

  const { data: userData, isLoading: isLoadingProfile } = useProfile();

  return (
    <Header height={60}>
      <Container maw="unset" className={classes.header}>
        <Flex gap="sm" justify="flex-start" align="center" direction="row">
          <Burger
            className={classes.hiddenDesktop}
            opened={opened}
            onClick={() => setOpened(!opened)}
            size="sm"
            color={theme.colors.gray[6]}
            mr="xl"
          />
          <Text fz="lg" className={classes.hiddenMobile}>
            Tsuwari
          </Text>
        </Flex>
        <Group position="center">
          <ActionIcon
            size="lg"
            variant="default"
            color={theme.colorScheme === 'dark' ? 'yellow' : 'blue'}
            onClick={() => toggleColorScheme()}
            title="Toggle color scheme"
          >
            {theme.colorScheme === 'dark' ? <IconSun size={18} /> : <IconMoonStars size={18} />}
          </ActionIcon>
          <Menu transition="pop" shadow="md" withArrow width={200}>
            <Menu.Target>
              <ActionIcon size="lg" title="Toggle language" variant="default">
                <IconLanguage size={18} />
              </ActionIcon>
            </Menu.Target>
            <Menu.Dropdown>
              <Menu.Label>Change language</Menu.Label>
              <Menu.Divider />
              {Array.from(LOCALES.entries()).map(([lang, { icon, name }]) => (
                <Menu.Item
                  style={{ fontWeight: lang === locale ? 'bold' : 'initial' }}
                  icon={icon}
                  key={lang}
                  onClick={() => toggleLocale(lang)}
                >
                  {name}
                </Menu.Item>
              ))}
            </Menu.Dropdown>
          </Menu>
          {isLoadingProfile && <Loader />}
          {!isLoadingProfile && userData && <Profile user={userData} />}
        </Group>
      </Container>
    </Header>
  );
}
