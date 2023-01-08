import { AppShell, ColorScheme, ColorSchemeProvider, MantineProvider, useMantineTheme } from '@mantine/core';
import { useColorScheme } from '@mantine/hooks';
import { QueryClientProvider } from '@tanstack/react-query';
import { getCookie, setCookie } from 'cookies-next';
import { GetServerSidePropsContext } from 'next';
import { AppProps } from 'next/app';
import Head from 'next/head';
import { useEffect, useState } from 'react';

import { SideBar } from '@/components/layout/sidebar';
import { queryClient } from '@/services/queryClient';

function App(props: AppProps & { colorScheme: ColorScheme }) {
  const { Component } = props;
  const theme = useMantineTheme();
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

  return (
    <>
      <Head>
        <title>Tsuwari</title>
        <meta name="viewport" content="minimum-scale=1, initial-scale=1, width=device-width"/>
      </Head>

      <QueryClientProvider client={queryClient}>
        <ColorSchemeProvider colorScheme={colorScheme} toggleColorScheme={toggleColorScheme}>
          <MantineProvider theme={{ colorScheme }} withGlobalStyles withNormalizeCSS>
            <AppShell
              styles={{
                main: {
                  background: colorScheme === 'dark' ? 'dark.8' : 'gray.0',
                  padding: 0,
                  width: '100%',
                },
              }}
              navbar={<SideBar opened={sidebarOpened} setOpened={setSidebarOpened} />}
              // header={<NavBar setOpened={setSidebarOpened} opened={sidebarOpened} />}
            >
              <AppShell styles={{
                main: {
                  background: props.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[0],
                },
              }}
              >
                <Component
                  styles={{
                    main: {
                      background: colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[0],
                    },
                  }}
                />
              </AppShell>
            </AppShell>
          </MantineProvider>
        </ColorSchemeProvider>
      </QueryClientProvider>
    </>
  );
}

App.getInitialProps = ({ ctx }: { ctx: GetServerSidePropsContext }) => ({
  colorScheme: getCookie('color-scheme', ctx),
});

export default App;
