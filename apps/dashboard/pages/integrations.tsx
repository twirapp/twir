import { Grid } from '@mantine/core';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';

import { DonationAlertsIntegration } from '../components/integrations/donationalerts';
import { LastfmIntegration } from '../components/integrations/lastfm';
import { SpotifyIntegration } from '../components/integrations/spotify';
import { StreamlabsIntegration } from '../components/integrations/streamlabs';
import { VKIntegration } from '../components/integrations/vk';

// @ts-ignore
export const getServerSideProps = async ({ locale }) => ({
    props: {
        ...(await serverSideTranslations(locale, ['integrations', 'layout'])),
    },
});

export default function Integrations() {
  return (
    <Grid grow>
      <Grid.Col span={4}>
        <SpotifyIntegration />
      </Grid.Col>
      <Grid.Col span={4}>
        <LastfmIntegration />
      </Grid.Col>
      <Grid.Col span={4}>
        <VKIntegration />
      </Grid.Col>
      <Grid.Col span={4}>
        <DonationAlertsIntegration />
      </Grid.Col>
      <Grid.Col span={4}>
        <StreamlabsIntegration />
      </Grid.Col>
      <Grid.Col span={4}>
        <VKIntegration />
      </Grid.Col>
    </Grid>
  );
}
