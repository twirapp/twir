import { useRouter } from 'next/router';
import { useEffect } from 'react';

import { useSpotifyIntegration } from '@/services/api/integrations';

export default function SpotifyLogin() {
  const router = useRouter();
  const manager = useSpotifyIntegration();

  useEffect(() => {
    const code = router.query.code;
    console.log(router.pathname, router.query);
    if (typeof code !== 'string') {
      router.push('/integrations');
    } else {
      manager.postCode(code).finally(() => {
        router.push('/integrations');
      });
    }
  }, [router.query]);
}
