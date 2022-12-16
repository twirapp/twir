import { Grid } from '@mantine/core';

import { DonationAlertsIntegration } from '../components/integrations/donationalerts';
import { LastfmIntegration } from '../components/integrations/lastfm';
import { SpotifyIntegration } from '../components/integrations/spotify';
import { StreamlabsIntegration } from '../components/integrations/streamlabs';
import { VKIntegration } from '../components/integrations/vk';

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
