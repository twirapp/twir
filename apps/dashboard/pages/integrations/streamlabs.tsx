import { useRouter } from 'next/router';
import { useEffect } from 'react';

import { useStreamLabsIntegration } from '@/services/api/integrations';
import { useSelectedDashboard } from '@/services/dashboard';

export default function DonationAlertsLogin() {
  const router = useRouter();
  const manager = useStreamLabsIntegration();
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
