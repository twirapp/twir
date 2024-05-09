import { createGlobalState, refDebounced } from '@vueuse/core'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { useBadges } from '../../manage-badges/composables/use-badges.js'

export type FilterType = 'status' | 'badge'

interface Filter {
	group: string
	type: FilterType
	list: {
		label: string
		key: string
		image?: string
	}[]
}

export const useUsersTableFilters = createGlobalState(() => {
	const { t } = useI18n()

	const searchInput = ref('')
	const debounceSearchInput = refDebounced(searchInput, 500)

	const { badges } = useBadges()

	const selectedStatuses = ref<Record<string, true | undefined>>({})
	const selectedBadges = ref<string[]>([])

	const selectedFiltersCount = computed(() => {
		return Object.keys(selectedStatuses.value).length + selectedBadges.value.length
	})

	const filtersList = computed<Filter[]>(() => [
		{
			group: t('adminPanel.manageUsers.statusGroup'),
			type: 'status',
			list: [
				{
					label: t('adminPanel.manageUsers.isAdmin'),
					key: 'isBotAdmin',
				},
				{
					label: t('adminPanel.manageUsers.isBanned'),
					key: 'isBanned',
				},
				{
					label: t('adminPanel.manageUsers.isBotEnabled'),
					key: 'isBotEnabled',
				},
			],
		},
		{
			group: t('adminPanel.manageUsers.badgesGroup'),
			type: 'badge',
			list: badges.value.map((badge) => ({
				label: badge.name,
				key: badge.id,
				image: badge.fileUrl,
			})),
		},
	])

	function clearFilters() {
		selectedStatuses.value = {}
		selectedBadges.value = []
	}

	function setFilterValue(filterKey: string, type: FilterType) {
		if (type === 'status') {
			if (selectedStatuses.value[filterKey]) {
				delete selectedStatuses.value[filterKey]
				return
			}

			selectedStatuses.value[filterKey] = true
		}

		if (type === 'badge') {
			if (selectedBadges.value.includes(filterKey)) {
				selectedBadges.value = selectedBadges.value.filter((badge) => badge !== filterKey)
				return
			}

			selectedBadges.value.push(filterKey)
		}
	}

	function isFilterApplied(filterKey: string, type: FilterType): boolean {
		if (type === 'status') {
			return filterKey in selectedStatuses.value
		}

		if (type === 'badge') {
			return selectedBadges.value.includes(filterKey)
		}

		return false
	}

	return {
		searchInput,
		debounceSearchInput,
		filtersList,
		selectedStatuses,
		selectedFiltersCount,
		setFilterValue,
		isFilterApplied,
		clearFilters,
		selectedBadges,
	}
})
