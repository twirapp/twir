import { createGlobalState } from '@vueuse/core'
import { ref, toRaw } from 'vue'

import { useModerationManager } from '@/api'
import { availableSettings, useEditableItem } from '@/features/moderation/ui/form/helpers.ts'

export const useModerationRules = createGlobalState(() => {
	const manager = useModerationManager()
	const { data: settings } = manager.getAll({})
	const { editableItem } = useEditableItem()
	const settingsOpened = ref(false)

	function showSettings(id: string) {
		const item = settings.value?.body.find(i => i.id === id)
		if (!item) return

		editableItem.value = toRaw(item)
		settingsOpened.value = true
	}

	async function createNewItem(itemType: string) {
		const defaultSettings = availableSettings.find(s => s.type === itemType)
		if (!defaultSettings) return
		editableItem.value = {
			data: structuredClone(defaultSettings),
		}
		settingsOpened.value = true
	}

	return {
		settings,
		showSettings,
		createNewItem,
		settingsOpened,
	}
})
