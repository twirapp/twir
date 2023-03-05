import { Flex, Grid } from '@mantine/core';
import { NextPage } from 'next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';

import { OBSOverlay } from '@/components/overlays/obs';
import { TTSOverlay } from '@/components/overlays/tts';

const Overlays: NextPage = () => {
  return (
   <Flex direction={'row'} gap={'md'}>
     <TTSOverlay />
     <OBSOverlay />
   </Flex>
  );
};

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
export const getServerSideProps = async ({ locale }) => ({
  props: {
    ...(await serverSideTranslations(locale, ['events', 'layout', 'commands', 'application'])),
  },
});

export default Overlays;