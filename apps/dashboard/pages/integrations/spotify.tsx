import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import { useRouter } from 'next/router';
import { useEffect } from 'react';

import { useSpotifyIntegration } from '@/services/api/integrations';
import { useSelectedDashboard } from '@/services/dashboard';

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
export const getServerSideProps = async ({ locale }) => ({
  props: {
    ...(await serverSideTranslations(locale, ['integrations', 'layout'])),
  },
});

export default function SpotifyLogin() {
  const router = useRouter();
  const manager = useSpotifyIntegration();
  const [dashboard] = useSelectedDashboard();

  useEffect(() => {
    if (!dashboard) {
      return;
    }
    const code = router.query.code;

    if (typeof code !== 'string') {
      router.push('/integrations');
    } else {
      manager.postCode(code).finally(() => {
        router.push('/integrations');
      });
    }
  }, [dashboard]);
}
