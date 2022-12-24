import { AppShell, ColorScheme, useMantineTheme } from '@mantine/core';
import { setCookie } from 'cookies-next';
import { useRouter } from 'next/router';
import { useEffect, useState } from 'react';

import { NavBar } from '@/components/layout/navbar';
import { SideBar } from '@/components/layout/sidebar';
import { useProfile } from '@/services/api';
import { SELECTED_DASHBOARD_KEY, useLocale, useSelectedDashboard } from '@/services/dashboard';
import i18nconfig from '../next-i18next.config'

type Props = React.PropsWithChildren<{
  colorScheme: ColorScheme
}>

const supportedLocales = i18nconfig.i18n.locales

export const AppProvider: React.FC<Props> = (props) => {
  const [selectedDashboard] = useSelectedDashboard();
  const router = useRouter();
  const { error: profileError } = useProfile();
  const [locale, setLocale] = useLocale();

  useEffect(() => {
    if (selectedDashboard) {
      console.log(JSON.stringify(selectedDashboard));
      setCookie(SELECTED_DASHBOARD_KEY, selectedDashboard.channelId, {
        // 1 month
        expires: new Date(Date.now() + 2_629_700_000),
      });
    }
  }, [selectedDashboard]);

  useEffect(() => {
    // redirect to route with setted locale
    if (locale && supportedLocales.includes(locale)) {
      const { pathname, asPath, query } = router;
      if (query.code || query.token) {
        return;
      }
      router.push({ pathname, query }, asPath, { locale });
    }

    // set locale if unsupported is setted
    if (!supportedLocales.includes(locale)) {
      setLocale('en');
    }
  }, [locale]);

  useEffect(() => {
    if (profileError) {
      window.location.replace(`${window.location.origin}`);
    }
  }, [profileError]);

  const theme = useMantineTheme();
  const [sidebarOpened, setSidebarOpened] = useState(false);

  return (
    <AppShell
      styles={{
        main: {
          background:
            props.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[0],
        },
      }}
      navbarOffsetBreakpoint="sm"
      asideOffsetBreakpoint="sm"
      navbar={<SideBar opened={sidebarOpened} setOpened={setSidebarOpened} />}
      header={<NavBar setOpened={setSidebarOpened} opened={sidebarOpened} />}
    >
      {props.children}
    </AppShell>
  );
};