import { Grid } from '@mantine/core';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';

import { DashboardAccess } from '../components/settings/dashboardAccess';

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
export const getServerSideProps = async ({ locale }) => ({
    props: {
        ...(await serverSideTranslations(locale, ['settings', 'layout'])),
    },
});

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
