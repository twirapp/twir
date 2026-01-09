import { useAdminUsers } from '#layers/dashboard/api/admin/users'
import { createGlobalState } from '@vueuse/core'

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
