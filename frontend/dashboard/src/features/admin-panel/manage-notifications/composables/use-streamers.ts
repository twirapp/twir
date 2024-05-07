import { createGlobalState } from '@vueuse/core'
import { computed } from 'vue'

import { useStreamers as useStreamersApi } from '@/api/admin/streamers'

export const useStreamers = createGlobalState(() => {
	const { data } = useStreamersApi()

	const streamers = computed(() => {
		if (!data.value) return []
		return data.value.streamers
	})

	return {
		streamers,
	}
})
