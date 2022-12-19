import { useRouter } from 'next/router';
import { useEffect } from 'react';

import { useLastfmIntegration } from '@/services/api/integrations';

export default function LastfmLogin() {
  const router = useRouter();
  const manager = useLastfmIntegration();

  useEffect(() => {
    const token = router.query.token;

    if (typeof token !== 'string') {
      router.push('/integrations');
    } else {
      manager.postToken(token).finally(() => {
        router.push('/integrations');
      });
    }
  }, [router.query]);
}
