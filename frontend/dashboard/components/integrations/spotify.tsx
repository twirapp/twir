import { useTranslation } from 'next-i18next';

import { useSpotify } from '@/services/api/integrations';
import { OAuthComponent } from './oauthComponent';

export const SpotifyIntegration: React.FC = () => {
	const manager = useSpotify();
	const logout = manager.useLogout();
	const { t } = useTranslation('integrations');
	const auth = manager.useGetAuthLink();

	const { data: profile } = manager.useData();

	async function login() {
		if (auth.data) {
			window.location.replace(auth.data);
		}
	}

	return (
		<OAuthComponent
			integrationKey={'spotify'}
			logout={logout.mutate}
			login={login}
			profile={profile ? {
				name: profile?.display_name,
				avatar: profile?.images?.at(0)?.url,
			} : undefined}
		/>
	);
};
