import { ColorScheme, ColorSchemeProvider, MantineProvider } from '@mantine/core';
import { useColorScheme } from '@mantine/hooks';
import { QueryClientProvider } from '@tanstack/react-query';
import { getCookie, setCookie } from 'cookies-next';
import { GetServerSidePropsContext } from 'next';
import { AppProps } from 'next/app';
import Head from 'next/head';
import { useEffect, useState } from 'react';

import { AppLayout } from '../components/layout';
import { queryClient } from '../services/queryClient';

// put in constants.ts
const ONE_MONTH = 2_629_700_000;

interface Props {
  dashboardId: string | null | undefined;
  locale: string;
  colorScheme: ColorScheme;
}

function App(props: AppProps & Props) {
  const { Component } = props;

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

  return (
    <>
      <Head>
        <title>Tsuwari</title>
        <meta name="viewport" content="minimum-scale=1, initial-scale=1, width=device-width" />
      </Head>

      <QueryClientProvider client={queryClient}>
        <ColorSchemeProvider colorScheme={colorScheme} toggleColorScheme={toggleColorScheme}>
          <MantineProvider theme={{ colorScheme }} withGlobalStyles withNormalizeCSS>
            <AppLayout colorScheme={colorScheme}>
              <Component />
            </AppLayout>
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
