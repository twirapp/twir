import { Grid } from '@mantine/core';
import { NextPage } from 'next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';

import { TTSOverlay } from '@/components/overlays/tts';

const Overlays: NextPage = () => {
  return (
   <Grid>
     <Grid.Col span={2}><TTSOverlay /></Grid.Col>
   </Grid>
  );
};

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
export const getServerSideProps = async ({ locale }) => ({
  props: {
    ...(await serverSideTranslations(locale, ['events', 'layout'])),
  },
});

export default Overlays;