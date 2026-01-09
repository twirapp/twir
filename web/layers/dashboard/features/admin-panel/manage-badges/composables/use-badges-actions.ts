import { createGlobalState } from '@vueuse/core'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

import { useBadgesForm } from './use-badges-form.js'
import { useBadges } from './use-badges.js'
import { useUsersTableFilters } from '../../manage-users/composables/use-users-table-filters.js'

import type { Badge } from '~/gql/graphql'

import { useLayout } from '~/composables/use-layout'

export const useBadgesActions = createGlobalState(() => {
	const layout = useLayout()
	const badgesForm = useBadgesForm()
	const { badgesDelete, badgesUpdate } = useBadges()

	const router = useRouter()
	const userFilters = useUsersTableFilters()

	const isShowModalDelete = ref(false)
	const deleteBadgeId = ref<string | null>(null)

	async function deleteBadge() {
		if (!deleteBadgeId.value)
			return
		await badgesDelete.executeMutation({ id: deleteBadgeId.value })
		deleteBadgeId.value = null
	}

	function editBadge(badge: Badge) {
		badgesForm.editableBadgeId.value = badge.id
		badgesForm.nameField.fieldModel.value = badge.name
		badgesForm.fileField.fieldModel.value = badge.fileUrl
		badgesForm.slotField.fieldModel.value = badge.ffzSlot
		layout.scrollToTop()
	}

	function toggleBadgeEnabled(badge: Badge) {
		badgesUpdate.executeMutation({
			id: badge.id,
			opts: { enabled: !badge.enabled },
		})
	}

	function showModalDeleteBadge(badge: Badge): void {
		isShowModalDelete.value = true
		deleteBadgeId.value = badge.id
	}

	function applyUserSearchBadgeFilter(badge: Badge): void {
		userFilters.clearFilters()
		userFilters.selectedBadges.value.push(badge.id)
		router.push({ query: { tab: 'users' } })
	}

	return {
		isShowModalDelete,
		editBadge,
		deleteBadge,
		toggleBadgeEnabled,
		showModalDeleteBadge,
		applyUserSearchBadgeFilter,
	}
})
