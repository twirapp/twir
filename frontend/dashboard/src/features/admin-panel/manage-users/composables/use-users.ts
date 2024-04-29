import { defineStore } from 'pinia'

import { useAdminUsers } from '@/api/admin/users.js'

export const useUsers = defineStore('admin-panel/users', () => {
	const usersApi = useAdminUsers()
	const switchBan = usersApi.useMutationUserSwitchBan()
	const switchAdmin = usersApi.useMutationUserSwitchAdmin()

	return {
		usersApi,
		switchBan,
		switchAdmin
	}
})
