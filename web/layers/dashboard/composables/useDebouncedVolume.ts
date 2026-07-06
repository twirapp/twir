import { useDebounceFn } from '@vueuse/core'
import { ref } from 'vue'

export function useDebouncedVolume(
	getChannelId: () => string,
	setVolumeMutation: (vars: { channelId: string; volume: number }) => any,
) {
	const localVolume = ref(100)

	const debouncedSync = useDebounceFn((volume: number) => {
		const id = getChannelId()
		if (!id) return
		setVolumeMutation({ channelId: id, volume })
	}, 300)

	function handleVolumeInput(value: number) {
		localVolume.value = value
		debouncedSync(value)
	}

	function syncFromServer(volume: number) {
		localVolume.value = volume
	}

	return {
		localVolume,
		handleVolumeInput,
		syncFromServer,
	}
}
