import { createGlobalState } from '@vueuse/core'

import { useAdminUsers } from '@/api/admin/users.js'

export const useUsers = createGlobalState(() => {
	const usersApi = useAdminUsers()
	const switchBan = usersApi.useMutationUserSwitchBan()
	const switchAdmin = usersApi.useMutationUserSwitchAdmin()

	return {
		usersApi,
		switchBan,
		switchAdmin,
	}
})
