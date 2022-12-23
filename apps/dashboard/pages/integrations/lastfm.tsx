import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import { useRouter } from 'next/router';
import { useEffect } from 'react';

import { useLastfm } from '@/services/api/integrations';
import { useSelectedDashboard } from '@/services/dashboard';

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
export const getServerSideProps = async ({ locale }) => ({
  props: {
    ...(await serverSideTranslations(locale, ['integrations', 'layout'])),
  },
});

export default function LastfmLogin() {
  const router = useRouter();
  const manager = useLastfm();
  const [dashboard] = useSelectedDashboard();
  useEffect(() => {
    if (!dashboard) {
      return;
    }
    
    const token = router.query.token;

    if (typeof token !== 'string') {
      router.push('/integrations');
    } else {
      manager.postCode.mutateAsync({ code: token }).finally(() => {
        router.push('/integrations');
      });
    }
  }, [dashboard]);
}
