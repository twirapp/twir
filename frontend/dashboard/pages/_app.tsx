import { AppShell, ColorScheme, ColorSchemeProvider, MantineProvider } from '@mantine/core';
import { useColorScheme } from '@mantine/hooks';
import { ModalsProvider } from '@mantine/modals';
import { NotificationsProvider } from '@mantine/notifications';
import { SpotlightProvider } from '@mantine/spotlight';
import { IconSearch } from '@tabler/icons';
import { QueryClientProvider } from '@tanstack/react-query';
import { getCookie, setCookie } from 'cookies-next';
import { Provider as JotaiProvider } from 'jotai';
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
import { InternalObsWebsocketProvider, OBSWebsocketProvider } from '@/services/obs/provider';
import { SelectedDashboardContext } from '@/services/selectedDashboardProvider';

import '../styles/global.css';
import { obsStore } from '../stores/obs';


// put in constants.ts
const ONE_MONTH = 2_629_700_000;

interface Props {
  dashboardId: string | null | undefined;
  locale: string;
  colorScheme: ColorScheme;
}

function App(props: AppProps & Props) {
  const { Component } = props;
  const [selectedDashboard, setSelectedDashboard] = useState<string>(props.dashboardId || '');
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
    setCookie('color_scheme', newColorScheme, {
      expires: new Date(Date.now() + ONE_MONTH * 12),
    });
  };

  useEffect(() => {
    if (!props.colorScheme) {
      toggleColorScheme(preferenceColorScheme);
    }
  }, [preferenceColorScheme]);

  useEffect(() => {
    if (selectedDashboard) {
      setCookie('dashboard_id', selectedDashboard, {
        expires: new Date(Date.now() + ONE_MONTH),
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
                    <JotaiProvider store={obsStore}>
                      <AppProvider colorScheme={colorScheme}>
                        <Component
                          styles={{
                            main: {
                              background: colorScheme === 'dark' ? 'dark.8' : 'gray.0',
                            },
                          }}
                        />
                      </AppProvider>
                    </JotaiProvider>
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
  locale: getCookie('locale', ctx),
  colorScheme: getCookie('color_scheme', ctx),
  dashboardId: getCookie('dashboard_id', ctx),
});

export default appWithTranslation(App, i18nconfig);
