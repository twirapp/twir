import { createGlobalState, refDebounced } from '@vueuse/core'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { AuditLogSystem, AuditOperationType } from '@/gql/graphql.js'

export type AuditSearchType = 'channel' | 'actor'
export type AuditFilterType = 'operation-type' | 'system'

interface Filter {
	group: string
	type: AuditFilterType
	list: {
		label: string
		key: string
	}[]
}

function convertFilterKey(key: string): string {
	return key.toLowerCase().replaceAll('_', '-')
}

const defaultFilters = {
	'operation-type': [],
	'system': [],
}

export const useAuditFilters = createGlobalState(() => {
	const { t } = useI18n()

	const searchUserId = ref('')
	const debounceSearchUserId = refDebounced(searchUserId, 500)

	const searchType = ref<AuditSearchType>('channel')
	const searchOptions = computed<{ label: string, value: AuditSearchType }[]>(() => [
		{
			label: t('dashboard.widgets.audit-logs.search.channel'),
			value: 'channel',
		},
		{
			label: t('dashboard.widgets.audit-logs.search.actor'),
			value: 'actor',
		},
	])

	const selectedFilters = ref<Record<AuditFilterType, string[]>>(structuredClone(defaultFilters))
	const selectedFiltersCount = computed(() => {
		return selectedFilters.value['operation-type'].length + selectedFilters.value.system.length
	})

	const filtersList = computed<Filter[]>(() => {
		const systemList: Filter['list'] = Object.values(AuditLogSystem).map((system) => {
			return {
				label: t(`dashboard.widgets.audit-logs.systems.${convertFilterKey(system)}`),
				key: system,
			}
		})

		const operationTypeList: Filter['list'] = Object.values(AuditOperationType).map((operationType) => {
			return {
				label: t(`dashboard.widgets.audit-logs.operation-type.${convertFilterKey(operationType)}`),
				key: operationType,
			}
		})

		return [
			{
				group: t('dashboard.widgets.audit-logs.operation-type-label'),
				type: 'operation-type',
				list: operationTypeList,
			},
			{
				group: t('dashboard.widgets.audit-logs.systems-label'),
				type: 'system',
				list: systemList,
			},
		]
	})

	function clearFilters() {
		selectedFilters.value = structuredClone(defaultFilters)
	}

	function setFilterValue(type: AuditFilterType, value: string) {
		const filter = selectedFilters.value[type]
		if (filter.includes(value)) {
			const valueIndex = filter.findIndex((filter) => filter === value)
			filter.splice(valueIndex, 1)
		} else {
			filter.push(value)
		}
	}

	function isFilterApplied(type: AuditFilterType, value: string): boolean {
		return selectedFilters.value[type].includes(value)
	}

	return {
		searchUserId: debounceSearchUserId,
		searchType,
		searchOptions,

		filtersList,
		selectedFilters,
		selectedFiltersCount,
		setFilterValue,
		isFilterApplied,
		clearFilters,
	}
})
