import { defineStore } from 'pinia';

import { useAdminUserSwitcher } from '@/api/manage-users';

export const useUsersActions = defineStore('manage-users/users-actions', () => {
	const userSwitcher = useAdminUserSwitcher();
	const adminUserSwitch = userSwitcher.useUserSwitchAdmin();
	const banUserSwitch = userSwitcher.useUserSwitchBan();

	function switchAdmin(userId: string) {
		adminUserSwitch.mutate(userId);
	}

	function switchBan(userId: string) {
		banUserSwitch.mutate(userId);
	}

	return {
		switchAdmin,
		switchBan,
	};
});
