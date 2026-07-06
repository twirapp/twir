import { useDebounceFn } from '@vueuse/core'
import { ref } from 'vue'

export function useDebouncedSeek(
	getChannelId: () => string,
	updatePositionMutation: (vars: { channelId: string; position: number }) => any,
) {
	const localPosition = ref(0)
	const isSeeking = ref(false)

	const debouncedSync = useDebounceFn((position: number) => {
		const id = getChannelId()
		isSeeking.value = false
		if (!id) return
		updatePositionMutation({ channelId: id, position })
	}, 300)

	function handleSeekInput(position: number) {
		localPosition.value = position
		isSeeking.value = true
		debouncedSync(position)
	}

	function syncFromServer(position: number) {
		if (isSeeking.value) return
		localPosition.value = position
	}

	return {
		localPosition,
		handleSeekInput,
		syncFromServer,
	}
}
