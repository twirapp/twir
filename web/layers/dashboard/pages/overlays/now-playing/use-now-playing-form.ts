import { ref } from 'vue'

export const defaultSettings = {
	preset: 'TRANSPARENT',
	channelId: null,
	id: null,
}

export function useNowPlayingForm() {
	const data = ref(null)

	function setData(entity: any) {
		data.value = entity
	}

	return { data, setData, defaultSettings }
}
