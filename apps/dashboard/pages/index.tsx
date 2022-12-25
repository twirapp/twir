import { Grid } from '@mantine/core';
import type { NextPage } from 'next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';

import { BotWidget } from '@/components/dashboard/bot';
import { YoutubePlayer } from '@/components/dashboard/youtoubePlayer';

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
export const getServerSideProps = async ({ locale }) => ({
    props: {
        ...(await serverSideTranslations(locale, ['dashboard', 'layout'])),
    },
});

const cols = {
  xs: 12,
  sm: 12,
  md: 6,
  lg: 4,
  xl: 4,
};

const Home: NextPage = () => {
  return (
    <Grid>
      <Grid.Col {...cols}><YoutubePlayer /></Grid.Col>
      <Grid.Col {...cols}><BotWidget /></Grid.Col>
    </Grid>
  );
};

export default Home;
