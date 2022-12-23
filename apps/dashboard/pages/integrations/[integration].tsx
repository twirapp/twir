import { LoadingOverlay } from '@mantine/core';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import { useRouter } from 'next/router';
import { useEffect } from 'react';

import { useDonationAlerts, useFaceit, useLastfm, useSpotify, useStreamlabs, useVK, type Integration as IntegrationType } from '@/services/api/integrations';
import { useSelectedDashboard } from '@/services/dashboard';

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
export const getServerSideProps = async ({ locale }) => ({
  props: {
    ...(await serverSideTranslations(locale, ['integrations', 'layout'])),
  },
});

export default function Integration() {
  const router = useRouter();

  const managers: Record<string, IntegrationType<any>> = {
    'donationalerts': useDonationAlerts(),
    'faceit': useFaceit(),
    'lastfm': useLastfm(),
    'spotify': useSpotify(),
    'streamlabs': useStreamlabs(),
    'vk': useVK(),
  };
  const [dashboard] = useSelectedDashboard();

  useEffect(() => {
    if (!dashboard) {
      return;
    }

    const { integration, code } = router.query;

    if (typeof code !== 'string' || !(code in managers)) {
      router.push('/integrations');
      return;
    }

    const manager = managers[code];
    manager.postCode.mutateAsync({ code }).finally(() => {
      router.push('/integrations');
    });
  }, [dashboard]);

  return <LoadingOverlay visible={true} overlayBlur={2} />;
}
