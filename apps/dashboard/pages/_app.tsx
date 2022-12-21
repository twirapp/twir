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
import { appWithTranslation } from 'next-i18next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import { AppProps } from 'next/app';
import Head from 'next/head';
import { useRouter } from 'next/router';
import { useEffect, useState } from 'react';
import { SWRConfig } from 'swr';

import { NavBar } from '@/components/layout/navbar';
import { SideBar } from '@/components/layout/sidebar';
import i18nconfig from '../next-i18next.config.js';

import { swrAuthFetcher, useProfile } from '@/services/api';
import { useLocale } from '@/services/dashboard/useLocale';

const app = function App(props: AppProps) {
  const { Component, pageProps } = props;

  const router = useRouter();
  const { error } = useProfile();
  const [locale] = useLocale();

  useEffect(() => {
    if (locale) {
      const { pathname, asPath, query } = router;
      router.push({ pathname, query }, asPath, { locale });
    }
  }, [locale]);

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
                      {...pageProps}
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
};

export default appWithTranslation(app, i18nconfig);
