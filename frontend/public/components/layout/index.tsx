import { AppShell, ColorScheme, useMantineTheme } from '@mantine/core';
import { useState } from 'react';

import { NavBar } from './navbar';
import { SideBar } from './sidebar';

type Props = React.PropsWithChildren<{
  colorScheme: ColorScheme;
}>;

export const AppLayout: React.FC<Props> = (props) => {
  const theme = useMantineTheme();
  const [sidebarOpened, setSidebarOpened] = useState(false);

  return (
    <AppShell
      styles={{
        main: {
          background: theme.colorScheme === 'dark' ? 'dark.8' : 'gray.0',
          padding: 0,
          width: '100%',
        },
      }}
      navbar={<SideBar opened={sidebarOpened} setOpened={setSidebarOpened} />}
      header={<NavBar setOpened={setSidebarOpened} opened={sidebarOpened} />}
    >
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
    </AppShell>
  );
};
