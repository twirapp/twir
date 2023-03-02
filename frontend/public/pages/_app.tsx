import {
  AppShell,
  ColorScheme,
  ColorSchemeProvider, Container,
  MantineProvider,
  useMantineTheme,
} from '@mantine/core';
import { useColorScheme } from '@mantine/hooks';
import { QueryClientProvider } from '@tanstack/react-query';
import { getCookie, setCookie } from 'cookies-next';
import { GetServerSidePropsContext } from 'next';
import { AppProps } from 'next/app';
import Head from 'next/head';
import { useEffect, useState } from 'react';

import { Header } from '@/components/layout/header';
import { queryClient } from '@/services/queryClient';

function App(props: AppProps & { colorScheme: ColorScheme }) {
  const { Component } = props;
  const theme = useMantineTheme();

  const preferenceColorScheme = useColorScheme(undefined, {
    getInitialValueInEffect: true,
  });

  const [colorScheme, setColorScheme] = useState<ColorScheme>(
    props.colorScheme ?? preferenceColorScheme,
  );

  const toggleColorScheme = (value?: ColorScheme) => {
    const newColorScheme = value || (colorScheme === 'dark' ? 'light' : 'dark');
    setColorScheme(newColorScheme);
    setCookie('color_scheme', newColorScheme, { maxAge: 86400000 * 365 });
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
        <meta name="viewport" content="minimum-scale=1, initial-scale=1, width=device-width" />
      </Head>

      <QueryClientProvider client={queryClient}>
        <ColorSchemeProvider colorScheme={colorScheme} toggleColorScheme={toggleColorScheme}>
          <MantineProvider theme={{ colorScheme }} withGlobalStyles withNormalizeCSS>
            <Container py="xl">
              <AppShell>
                <Header />
                <Component {...props.pageProps} />
              </AppShell>
              </Container>
          </MantineProvider>
        </ColorSchemeProvider>
      </QueryClientProvider>
    </>
  );
}

App.getInitialProps = ({ ctx }: { ctx: GetServerSidePropsContext }) => ({
  colorScheme: getCookie('color_scheme', ctx),
});

export default App;
