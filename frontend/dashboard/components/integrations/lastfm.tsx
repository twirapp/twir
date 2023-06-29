import { OAuthComponent } from './oauthComponent';

import { useLastfm } from '@/services/api/integrations';

export const LastfmIntegration: React.FC = () => {
	const manager = useLastfm();
	const { data: profile } = manager.useData();
	const logout = manager.useLogout();
	const auth = manager.useGetAuthLink();

	async function login() {
		if (auth.data) {
			window.location.replace(auth.data);
		}
	}

	return (
		<OAuthComponent
			integrationKey={'lastfm'}
			logout={logout.mutate}
			login={login}
			profile={profile ? {
				name: profile?.name,
				avatar: profile?.image,
			} : undefined}
		/>
	);
};
