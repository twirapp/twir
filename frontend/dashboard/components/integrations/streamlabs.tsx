import { OAuthComponent } from './oauthComponent';

import { useStreamlabs } from '@/services/api/integrations';

export const StreamlabsIntegration: React.FC = () => {
  const manager = useStreamlabs();
  const logout = manager.useLogout();
  const { data: profile } = manager.useData();
  const auth = manager.useGetAuthLink();

  async function login() {
    if (auth.data) {
      window.location.replace(auth.data);
    }
  }

  return (
		<OAuthComponent
			integrationKey={'streamlabs'}
			logout={logout.mutate}
			login={login}
			profile={profile ? {
				name: profile?.name,
				avatar: profile?.avatar,
			} : undefined}
		/>
  );
};
