import { AppShell, ColorScheme, ColorSchemeProvider, MantineProvider, useMantineTheme } from '@mantine/core';
import { useColorScheme, useHotkeys, useLocalStorage } from '@mantine/hooks';
import { ModalsProvider } from '@mantine/modals';
import { NotificationsProvider } from '@mantine/notifications';
import { SpotlightProvider } from '@mantine/spotlight';
import { IconSearch } from '@tabler/icons';
import { QueryClientProvider } from '@tanstack/react-query';
import { appWithTranslation } from 'next-i18next';
import { AppProps } from 'next/app';
import Head from 'next/head';
import { useState } from 'react';

import i18nconfig from '../next-i18next.config.js';

import { AppProvider } from '@/components/appProvider';
import { NavBar } from '@/components/layout/navbar';
import { SideBar } from '@/components/layout/sidebar';
import { queryClient } from '@/services/api';


const app = function App(props: AppProps) {
  const { Component } = props;

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
  const [sidebarOpened, setSidebarOpened] = useState(false);

  return (
    <>
      <Head>
        <title>Tsuwari</title>
        <meta name="viewport" content="minimum-scale=1, initial-scale=1, width=device-width"/>
      </Head>
      <QueryClientProvider client={queryClient}>
        <ColorSchemeProvider colorScheme={colorScheme} toggleColorScheme={toggleColorScheme}>
          <MantineProvider theme={{ colorScheme }} withGlobalStyles withNormalizeCSS>
            <NotificationsProvider position={'top-center'} limit={5}>
              <SpotlightProvider
                actions={[]}
                searchIcon={<IconSearch size={18}/>}
                searchPlaceholder="Search..."
                shortcut={['mod+k']}
                nothingFoundMessage="Nothing found..."
                limit={Number.MAX_SAFE_INTEGER}
                centered={true}
                styles={{
                  spotlight: {
                    marginBottom: 20,
                  },
                }}
              >
                <ModalsProvider>
                  <AppShell
                    styles={{
                      main: {
                        background:
                          colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[0],
                        padding: 0,
                        width: '100%',
                      },
                    }}
                    navbarOffsetBreakpoint="sm"
                    asideOffsetBreakpoint="sm"
                    navbar={<SideBar opened={sidebarOpened} setOpened={setSidebarOpened}/>}
                    header={<NavBar setOpened={setSidebarOpened} opened={sidebarOpened}/>}
                  >
                    <AppProvider colorScheme={colorScheme}><Component
                      styles={{
                        main: {
                          background:
                            colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[0],
                        },
                      }}
                    /></AppProvider>
                  </AppShell>
                </ModalsProvider>
              </SpotlightProvider>
            </NotificationsProvider>
          </MantineProvider>
        </ColorSchemeProvider>
      </QueryClientProvider>
    </>
  );
};

export default appWithTranslation(app, i18nconfig);
