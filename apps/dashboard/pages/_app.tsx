import { AppShell, ColorScheme, ColorSchemeProvider, MantineProvider } from '@mantine/core';
import { useColorScheme } from '@mantine/hooks';
import { ModalsProvider } from '@mantine/modals';
import { NotificationsProvider } from '@mantine/notifications';
import { SpotlightProvider } from '@mantine/spotlight';
import { IconSearch } from '@tabler/icons';
import { QueryClientProvider } from '@tanstack/react-query';
import { getCookie, setCookie } from 'cookies-next';
import { GetServerSidePropsContext } from 'next';
import { appWithTranslation } from 'next-i18next';
import { AppProps } from 'next/app';
import Head from 'next/head';
import { useEffect, useState } from 'react';

import i18nconfig from '../next-i18next.config.js';

import { AppProvider } from '@/components/appProvider';
import { NavBar } from '@/components/layout/navbar';
import { SideBar } from '@/components/layout/sidebar';
import { queryClient } from '@/services/api';
import { SelectedDashboardContext } from '@/services/selectedDashboardProvider';

function App(props: AppProps & { colorScheme: ColorScheme }) {
  const { Component } = props;
  const cookieSelectedDashboard = getCookie('selectedDashboard') as string | null | undefined;
  const [selectedDashboard, setSelectedDashboard] = useState<string>(cookieSelectedDashboard || '');
  const [sidebarOpened, setSidebarOpened] = useState(false);
  const preferenceColorScheme = useColorScheme(undefined, {
    getInitialValueInEffect: true,
  });

  const [colorScheme, setColorScheme] = useState<ColorScheme>(
    props.colorScheme ?? preferenceColorScheme,
  );

  const toggleColorScheme = (value?: ColorScheme) => {
    const newColorScheme = value || (colorScheme === 'dark' ? 'light' : 'dark');
    setColorScheme(newColorScheme);
    setCookie('color-scheme', newColorScheme, { maxAge: 86400000 * 365 });
  };

  useEffect(() => {
    if (!props.colorScheme) {
      toggleColorScheme(preferenceColorScheme);
    }
  }, [preferenceColorScheme]);

  useEffect(() => {
    if (selectedDashboard) {
      setCookie('selectedDashboard', selectedDashboard, {
        // 1 month
        expires: new Date(Date.now() + 2_629_700_000),
      });
    }
  }, [selectedDashboard]);

  return (
    <>
      <Head>
        <title>Tsuwari</title>
        <meta name="viewport" content="minimum-scale=1, initial-scale=1, width=device-width" />
      </Head>
      <SelectedDashboardContext.Provider
        value={{ id: selectedDashboard, setId: setSelectedDashboard }}
      >
        <QueryClientProvider client={queryClient}>
          <ColorSchemeProvider colorScheme={colorScheme} toggleColorScheme={toggleColorScheme}>
            <MantineProvider theme={{ colorScheme }} withGlobalStyles withNormalizeCSS>
              <NotificationsProvider position={'top-center'} limit={5}>
                <SpotlightProvider
                  actions={[]}
                  searchIcon={<IconSearch size={18} />}
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
                          background: colorScheme === 'dark' ? 'dark.8' : 'gray.0',
                          padding: 0,
                          width: '100%',
                        },
                      }}
                      navbar={<SideBar opened={sidebarOpened} setOpened={setSidebarOpened} />}
                      header={<NavBar setOpened={setSidebarOpened} opened={sidebarOpened} />}
                    >
                      <AppProvider colorScheme={colorScheme}>
                        <Component
                          styles={{
                            main: {
                              background: colorScheme === 'dark' ? 'dark.8' : 'gray.0',
                            },
                          }}
                        />
                      </AppProvider>
                    </AppShell>
                  </ModalsProvider>
                </SpotlightProvider>
              </NotificationsProvider>
            </MantineProvider>
          </ColorSchemeProvider>
        </QueryClientProvider>
      </SelectedDashboardContext.Provider>
    </>
  );
}

App.getInitialProps = ({ ctx }: { ctx: GetServerSidePropsContext }) => ({
  colorScheme: getCookie('color-scheme', ctx),
});

export default appWithTranslation(App, i18nconfig);
