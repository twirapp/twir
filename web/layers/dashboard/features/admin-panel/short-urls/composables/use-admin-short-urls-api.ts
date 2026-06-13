import { createGlobalState } from '@vueuse/core'
import { computed, readonly } from 'vue'

import type { AdminShortUrl } from '@/api/admin/short-urls.ts'
import type { AdminShortUrlsInput } from '@/gql/graphql.ts'

import { useAdminShortUrlsApi as useApiManager } from '@/api/admin/short-urls.ts'
import { toast } from 'vue-sonner'
import { usePagination } from '@/composables/use-pagination.ts'

export const useAdminShortUrlsApi = createGlobalState(() => {
	const api = useApiManager()
	const deleter = api.useDeleteMutation()
	const creator = api.useCreateMutation()

	const { pagination, setPagination } = usePagination()
	const params = computed<AdminShortUrlsInput>(() => ({
		page: pagination.value.pageIndex,
		perPage: pagination.value.pageSize,
	}))

	const { data, fetching } = api.useQuery(params)

	const list = computed<AdminShortUrl[]>(() => {
		if (!data.value?.adminShortUrls?.items) return []
		return data.value.adminShortUrls.items ?? []
	})

	const totalItems = computed(() => {
		return data.value?.adminShortUrls.total ?? 0
	})

	async function removeShortUrl(id: string) {
		try {
			const { error } = await deleter.executeMutation({ id })
			if (!error) {
				toast.success('Deleted', {
					duration: 2500,
				})
			}
		} catch (e) {
			toast.error(`${e}`, {
				duration: 2500,
			})
		}
	}

	async function createShortUrl(link: string, shortId?: string) {
		try {
			const { error } = await creator.executeMutation({
				input: {
					link,
					shortId,
				},
			})
			if (!error) {
				toast.success('Created', {
					duration: 2500,
				})
			}
		} catch (e) {
			toast.error(`${e}`, {
				duration: 2500,
			})
		}
	}

	return {
		list,
		totalItems,
		isLoading: readonly(fetching),

		pagination,
		setPagination,

		removeShortUrl,
		createShortUrl,
	}
})
