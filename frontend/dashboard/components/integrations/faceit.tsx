import { OAuthComponent } from './oauthComponent';

import { useFaceit } from '@/services/api/integrations';

export const FaceitIntegration: React.FC = () => {
	const manager = useFaceit();
	const { data: profile } = manager.useData();
	const auth = manager.useGetAuthLink();
	const logout = manager.useLogout();

	async function login() {
		if (auth.data) {
			window.location.replace(auth.data);
		}
	}

	return (
		<OAuthComponent
			integrationKey={'faceit'}
			logout={logout.mutate}
			login={login}
			profile={profile ? {
				name: profile?.name,
				avatar: profile?.avatar,
			} : undefined}
		/>
	);
};
