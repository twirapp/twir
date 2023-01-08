import { AppShell, ColorScheme, useMantineTheme } from '@mantine/core';
import { useRouter } from 'next/router';
import { useContext, useEffect, useState } from 'react';

import { NavBar } from '@/components/layout/navbar';
import { SideBar } from '@/components/layout/sidebar';
import { FetcherError, useProfile } from '@/services/api';
import { useLocale } from '@/services/dashboard';
import { SelectedDashboardContext } from '@/services/selectedDashboardProvider';

type Props = React.PropsWithChildren<{
  colorScheme: ColorScheme;
}>;

export const AppProvider: React.FC<Props> = (props) => {
  const dashboardContext = useContext(SelectedDashboardContext);

  const router = useRouter();
  const { error: profileError, data: profileData } = useProfile();
  const { locale, toggleLocale, isSupportedLocale } = useLocale();

  useEffect(() => {
    if (!dashboardContext.id && profileData) {
      dashboardContext.setId(profileData.id);
    }
  }, [profileData]);

  useEffect(() => {
    // redirect to route with setted locale
    if (isSupportedLocale()) {
      const { pathname, asPath, query } = router;
      if (query.code || query.token) {
        return;
      }
      router.push({ pathname, query }, asPath, { locale });
    } else {
      toggleLocale();
    }
  }, [locale]);

  useEffect(() => {
    if (profileError) {
      if (profileError instanceof FetcherError && profileError.status === 403) {
        window.location.replace(`/api/auth?state=${window.btoa(window.location.origin)}`);
      } else {
        window.location.replace(`${window.location.origin}`);
      }
    }
  }, [profileError]);

  const theme = useMantineTheme();
  const [sidebarOpened, setSidebarOpened] = useState(false);

  return (
    <AppShell
      styles={{
        main: {
          background: props.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[0],
          width: '100%',
        },
      }}
      navbar={<SideBar opened={sidebarOpened} setOpened={setSidebarOpened} />}
      header={<NavBar setOpened={setSidebarOpened} opened={sidebarOpened} />}
    >
      {props.children}
    </AppShell>
  );
};
