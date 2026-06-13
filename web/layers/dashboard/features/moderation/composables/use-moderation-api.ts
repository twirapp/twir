import { createGlobalState } from '@vueuse/core'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import type { ModerationCreateOrUpdateInput, ModerationItem } from '@/api/moderation.ts'

import { useChannelModerationSettingsApi } from '@/api/moderation.ts'
import { toast } from 'vue-sonner'

export const useModerationApi = createGlobalState(() => {
	const api = useChannelModerationSettingsApi()
	const deleter = api.useDelete()
	const updater = api.useUpdate()
	const creator = api.useCreate()

	const { t } = useI18n()

	const { data: moderationItems, fetching, executeQuery: refetchItems } = api.useQuery()
	const items = computed<ModerationItem[]>(() => {
		return moderationItems?.value?.moderationSettings ?? []
	})

	async function update(
		id: string,
		input: ModerationCreateOrUpdateInput
	): Promise<{ id: string } | undefined> {
		//
		// @ts-expect-error
		delete input.createdAt
		//
		// @ts-expect-error
		delete input.updatedAt
		//
		// @ts-expect-error
		delete input.id

		try {
			const result = await updater.executeMutation({ id, input })
			if (result.error) {
				throw result.error
			}

			toast.success('Updated', {
				duration: 2500,
			})

			return result.data?.moderationSettingsUpdate
		} catch (e) {
			toast.error(`Cannot update: ${e}`, {
				duration: 2500,
			})
		}
	}

	async function create(input: ModerationCreateOrUpdateInput): Promise<{ id: string } | undefined> {
		try {
			const result = await creator.executeMutation({ input })
			if (result.error) {
				throw result.error
			}

			toast.success('Created', {
				duration: 2500,
			})

			return result.data?.moderationSettingsCreate
		} catch (e) {
			toast.error(`Cannot create: ${e}`, {
				duration: 2500,
			})
		}
	}

	async function remove(id: string) {
		try {
			const result = await deleter.executeMutation({ id })
			if (result.error) {
				throw result.error
			}

			toast.success(t('sharedTexts.deleted'), {
				duration: 2500,
			})
		} catch (e) {
			toast.error(`Cannot delete: ${e}`, {
				duration: 2500,
			})
		}
	}

	async function fetchItems() {
		await refetchItems()
	}

	return {
		items,
		isLoading: fetching,
		fetchItems,

		update,
		create,
		remove,
	}
})
