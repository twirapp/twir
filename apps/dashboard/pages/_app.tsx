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
import { QueryClientProvider } from '@tanstack/react-query';
import { setCookie } from 'cookies-next';
import { appWithTranslation } from 'next-i18next';
import { AppProps } from 'next/app';
import Head from 'next/head';
import { useRouter } from 'next/router';
import { useEffect, useState } from 'react';

import i18nconfig from '../next-i18next.config.js';

import { AppProvider } from '@/components/appProvider';
import { NavBar } from '@/components/layout/navbar';
import { SideBar } from '@/components/layout/sidebar';
import { useProfile, queryClient } from '@/services/api';
import { SELECTED_DASHBOARD_KEY, useSelectedDashboard } from '@/services/dashboard';
import { useLocale } from '@/services/dashboard/useLocale';

const app = function App(props: AppProps) {
  const { Component, pageProps } = props;

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

  return (
    <>
      <Head>
        <title>Tsuwari</title>
        <meta name="viewport" content="minimum-scale=1, initial-scale=1, width=device-width" />
      </Head>
      <QueryClientProvider client={queryClient} contextSharing={true}>
      <AppProvider><Component
        styles={{
          main: {
            background:
              colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[0],
          },
        }}
        {...pageProps}
      /></AppProvider>
      </QueryClientProvider>
    </>
  );
};

export default appWithTranslation(app, i18nconfig);
