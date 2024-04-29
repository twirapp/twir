import { defineStore } from 'pinia'
import { computed } from 'vue'

import { useStreamers as useStreamersApi } from '@/api/admin/streamers'

export const useStreamers = defineStore('admin-panel/streamers', () => {
	const { data } = useStreamersApi()

	const streamers = computed(() => {
		if (!data.value) return []
		return data.value.streamers
	})

	return {
		streamers
	}
})
