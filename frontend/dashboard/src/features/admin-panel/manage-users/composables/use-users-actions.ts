import { defineStore } from 'pinia';

import { useAdminUserSwitcher } from '@/api/admin/users';

export const useUsersActions = defineStore('manage-users/users-actions', () => {
	const userSwitcher = useAdminUserSwitcher();
	const adminUserSwitch = userSwitcher.useUserSwitchAdmin();
	const banUserSwitch = userSwitcher.useUserSwitchBan();

	async function switchAdmin(userId: string) {
		await adminUserSwitch.mutateAsync(userId);
	}

	async function switchBan(userId: string) {
		await banUserSwitch.mutateAsync(userId);
	}

	return {
		switchAdmin,
		switchBan,
	};
});
