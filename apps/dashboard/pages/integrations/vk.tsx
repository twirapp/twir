import { useRouter } from 'next/router';
import { useEffect } from 'react';

import { useVkIntegration } from '@/services/api/integrations';

export default function LastfmLogin() {
  const router = useRouter();
  const manager = useVkIntegration();

  useEffect(() => {
    const code = router.query.code;

    if (typeof code !== 'string') {
      router.push('/integrations');
    } else {
      manager.postCode(code).finally(() => {
        router.push('/integrations');
      });
    }
  }, [router.query]);
}
