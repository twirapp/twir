import { type Ref, computed } from 'vue'

import { useProfile } from '~~/layers/dashboard/api/auth.js'
import { useFilesApi } from '~~/layers/dashboard/api/files.js'

import type { Layer, LayerSettings } from '../types'

type UpdateSettings = (updates: Partial<LayerSettings>) => void

export function useImageLayerEditor(layer: Ref<Layer>, updateSettings: UpdateSettings) {
	const { data: profile } = useProfile()
	const filesApi = useFilesApi()

	const imageUrl = computed({
		get: () => layer.value.settings.imageUrl,
		set: (value: string) => updateSettings({ imageUrl: value }),
	})

	function setPlaceholder() {
		updateSettings({ imageUrl: 'https://via.placeholder.com/300x200' })
	}

	function selectUploadedFile(fileId: string) {
		const channelId = profile.value?.selectedDashboardId
		if (!channelId) return

		updateSettings({ imageUrl: filesApi.computeFileUrl(channelId, fileId) })
	}

	function onUploadedFileDelete(deletedFileId: string) {
		if (layer.value.settings.imageUrl.endsWith(`/files/content/${deletedFileId}`)) {
			updateSettings({ imageUrl: '' })
		}
	}

	return { imageUrl, setPlaceholder, selectUploadedFile, onUploadedFileDelete }
}
