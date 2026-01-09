import { createGlobalState } from '@vueuse/core'
import { computed, readonly } from 'vue'

import type { AdminToxicMessagesInput, AdminToxicMessagesQuery } from '~/gql/graphql'

import { useToxicMessagesAdminApi } from '#layers/dashboard/api/admin/toxic-messages.ts'
import { usePagination } from '~/composables/use-pagination.ts'

export const useAdminToxicMessagesApi = createGlobalState(() => {
	const api = useToxicMessagesAdminApi()

	const { pagination, setPagination } = usePagination()
	const params = computed<AdminToxicMessagesInput>(() => ({
		page: pagination.value.pageIndex,
		perPage: pagination.value.pageSize,
	}))

	const { data, fetching } = api.useDataQuery(params)

	const list = computed<AdminToxicMessagesQuery['adminToxicMessages']['items']>(() => {
		if (!data.value) return []
		return data.value.adminToxicMessages.items
	})

	const totalItems = computed(() => {
		return data.value?.adminToxicMessages.total ?? 0
	})

	return {
		list,
		totalItems,
		isLoading: readonly(fetching),

		pagination,
		setPagination,
	}
})
