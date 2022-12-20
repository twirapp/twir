import {
  AppShell,
  ColorScheme,
  ColorSchemeProvider,
  MantineProvider,
  useMantineTheme,
} from '@mantine/core';
import { useColorScheme, useHotkeys, useLocalStorage } from '@mantine/hooks';
import { ModalsProvider } from '@mantine/modals';
import { NotificationsProvider } from '@mantine/notifications';
import { SpotlightProvider } from '@mantine/spotlight';
import { IconSearch } from '@tabler/icons';
import { AppProps } from 'next/app';
import Head from 'next/head';
import { useEffect, useState } from 'react';
import { SWRConfig } from 'swr';

import { NavBar } from '../components/layout/navbar';
import { SideBar } from '../components/layout/sidebar';

import { swrAuthFetcher, useProfile } from '@/services/api';

export default function App(props: AppProps) {
  const { Component } = props;

  const { error } = useProfile();

  useEffect(() => {
    if (error) {
      window.location.replace(`${window.location.origin}`);
    }
  }, [error]);

  const preferredColorScheme = useColorScheme();
  const [colorScheme, setColorScheme] = useLocalStorage<ColorScheme>({
    key: 'theme',
    defaultValue: preferredColorScheme,
    getInitialValueInEffect: true,
  });

  const toggleColorScheme = (value?: ColorScheme) =>
    setColorScheme(value || (colorScheme === 'dark' ? 'light' : 'dark'));

  useHotkeys([['mod+J', () => toggleColorScheme()]]);

  const theme = useMantineTheme();
  const [opened, setOpened] = useState(false);

  return (
    <>
      <Head>
        <title>Tsuwari</title>
        <meta name="viewport" content="minimum-scale=1, initial-scale=1, width=device-width" />
      </Head>
      <ColorSchemeProvider colorScheme={colorScheme} toggleColorScheme={toggleColorScheme}>
        <MantineProvider theme={{ colorScheme }} withGlobalStyles withNormalizeCSS>
          <NotificationsProvider>
            <SpotlightProvider
              actions={[]}
              searchIcon={<IconSearch size={18} />}
              searchPlaceholder="Search..."
              shortcut={['mod+k']}
              nothingFoundMessage="Nothing found..."
            >
              <SWRConfig
                value={{
                  fetcher: swrAuthFetcher,
                }}
              >
                <ModalsProvider>
                  <AppShell
                    styles={{
                      main: {
                        background:
                          colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[0],
                      },
                    }}
                    navbarOffsetBreakpoint="sm"
                    asideOffsetBreakpoint="sm"
                    navbar={<SideBar opened={opened} />}
                    header={<NavBar setOpened={setOpened} opened={opened} />}
                  >
                    <Component
                      styles={{
                        main: {
                          background:
                            colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[0],
                        },
                      }}
                    />
                  </AppShell>
                </ModalsProvider>
              </SWRConfig>
            </SpotlightProvider>
          </NotificationsProvider>
        </MantineProvider>
      </ColorSchemeProvider>
    </>
  );
}
