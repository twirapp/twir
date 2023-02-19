import { AppShell, ColorScheme, useMantineTheme } from '@mantine/core';
import { useAtom } from 'jotai/index';
import { useContext, useEffect, useState } from 'react';

import { externalObsWsAtom, internalObsWsAtom, MyOBSWebsocket } from '../stores/obs';

import { NavBar } from '@/components/layout/navbar';
import { SideBar } from '@/components/layout/sidebar';
import { FetcherError, useProfile } from '@/services/api';
import { useObsModule } from '@/services/api/modules';
import { useInternalObsWs, useObs } from '@/services/obs/hook';
import { SelectedDashboardContext } from '@/services/selectedDashboardProvider';

type Props = React.PropsWithChildren<{
  colorScheme: ColorScheme;
}>;

export const AppProvider: React.FC<Props> = (props) => {
  const obs = useObs();
  const internalObs = useInternalObsWs();
  const obsModule = useObsModule();
  const { data: obsSettings } = obsModule.useSettings();

  const dashboardContext = useContext(SelectedDashboardContext);
  const { error: profileError, data: profileData } = useProfile();

  useEffect(() => {
    if (!obsSettings || !profileData) return;
    obs.connect().then(() => {
      return internalObs.connect();
    });

    return () => {
      obs.disconnect().then(() => {
        return internalObs.disconnect();
      });
    };
  }, [obsSettings, profileData]);

  useEffect(() => {
    if (!dashboardContext.id && profileData) {
      dashboardContext.setId(profileData.id);
    }
  }, [profileData]);

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
