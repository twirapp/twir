import { useRouter } from 'next/router';
import { useEffect } from 'react';

import { useLastfmIntegration } from '@/services/api/integrations';
import { useSelectedDashboard } from '@/services/dashboard';

export default function LastfmLogin() {
  const router = useRouter();
  const manager = useLastfmIntegration();
  const [dashboard] = useSelectedDashboard();
  useEffect(() => {
    if (!dashboard) {
      return;
    }
    
    const token = router.query.token;

    if (typeof token !== 'string') {
      router.push('/integrations');
    } else {
      manager.postToken(token).finally(() => {
        router.push('/integrations');
      });
    }
  }, [dashboard]);
}
