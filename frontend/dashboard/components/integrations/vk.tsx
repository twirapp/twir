import { OAuthComponent } from './oauthComponent';

import { useVK } from '@/services/api/integrations';

export const VKIntegration: React.FC = () => {
	const manager = useVK();
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
			integrationKey={'vk'}
			logout={logout.mutate}
			login={login}
			profile={profile ? {
				name: `${profile.first_name} ${profile.last_name}`,
				avatar: profile?.photo_max_orig,
			} : undefined}
		/>
	);
};
