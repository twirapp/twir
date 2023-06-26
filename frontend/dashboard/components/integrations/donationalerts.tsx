import { OAuthComponent } from './oauthComponent';

import { useDonationAlerts } from '@/services/api/integrations';

export const DonationAlertsIntegration: React.FC = () => {
	const manager = useDonationAlerts();
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
			integrationKey={'donationalerts'}
			logout={logout.mutate}
			login={login}
			profile={profile ? {
				name: profile?.name,
				avatar: profile?.avatar,
			} : undefined}
		/>
	);
};
