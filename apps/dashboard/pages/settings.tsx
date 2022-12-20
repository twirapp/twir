import { Grid } from '@mantine/core';

import { DashboardAccess } from '../components/settings.tsx/dashboardAccess';

export default function Settings() {
  return (
    <Grid grow>
      <Grid.Col span={4}>
        <DashboardAccess />
      </Grid.Col>
      <Grid.Col span={4}></Grid.Col>
      <Grid.Col span={4}></Grid.Col>
    </Grid>
  );
}
