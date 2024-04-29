import { defineStore } from 'pinia'
import { computed } from 'vue'

import { useAdminBadges, useQueryBadges } from '@/api/admin/badges'

export const useBadges = defineStore('admin-panel/badges', () => {
	const badgesApi = useAdminBadges()
	const badgesCreate = badgesApi.useMutationCreateBadge()
	const badgesDelete = badgesApi.useMutationDeleteBadge()
	const badgesUpdate = badgesApi.useMutationUpdateBadge()
	const badgesAddUser = badgesApi.useMutationsAddUserBadge()
	const badgesRemoveUser = badgesApi.useMutationsRemoveUserBadge()

	const { data } = useQueryBadges()
	const badges = computed(() => data.value?.twirBadges ?? [])

	return {
		badges,
		badgesCreate,
		badgesDelete,
		badgesUpdate,
		badgesAddUser,
		badgesRemoveUser
	}
})
