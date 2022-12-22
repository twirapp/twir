import { AppShell, ColorScheme, ColorSchemeProvider, MantineProvider, useMantineTheme } from '@mantine/core';
import { useColorScheme, useHotkeys, useLocalStorage } from '@mantine/hooks';
import { ModalsProvider } from '@mantine/modals';
import { NotificationsProvider } from '@mantine/notifications';
import { SpotlightProvider } from '@mantine/spotlight';
import { IconSearch } from '@tabler/icons';
import { QueryClientProvider } from '@tanstack/react-query';
import { setCookie } from 'cookies-next';
import { useRouter } from 'next/router';
import { useEffect, useState } from 'react';

import { NavBar } from '@/components/layout/navbar';
import { SideBar } from '@/components/layout/sidebar';
import { queryClient, useProfile } from '@/services/api';
import { SELECTED_DASHBOARD_KEY, useLocale, useSelectedDashboard } from '@/services/dashboard';

export const AppProvider: React.FC<React.PropsWithChildren> = (props) => {
  const [selectedDashboard] = useSelectedDashboard();
  const router = useRouter();
  const { error: profileError, data } = useProfile();
  const [locale] = useLocale();

  useEffect(() => {
    if (selectedDashboard) {
      setCookie(SELECTED_DASHBOARD_KEY, selectedDashboard.channelId, {
        // 1 month
        expires: new Date(Date.now() + 2_629_700_000),
      });
    }
  }, [selectedDashboard]);

  useEffect(() => {
    if (locale) {
      const { pathname, asPath, query } = router;
      if (query.code || query.token) {
        return;
      }
      router.push({ pathname, query }, asPath, { locale });
    }
  }, [locale]);

  useEffect(() => {
    if (profileError) {
      window.location.replace(`${window.location.origin}`);
    }
  }, [profileError]);

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
                navbar={<SideBar opened={sidebarOpened} setOpened={setSidebarOpened} />}
                header={<NavBar setOpened={setSidebarOpened} opened={sidebarOpened} />}
              >
                {props.children}
              </AppShell>
            </ModalsProvider>
          </SpotlightProvider>
        </NotificationsProvider>
      </MantineProvider>
    </ColorSchemeProvider>
  );
};