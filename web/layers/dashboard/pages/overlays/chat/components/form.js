import { ref } from 'vue'

export function useChatOverlayForm() {
	const data = ref(null)

	function setData(entity: any) {
		data.value = entity
	}

	function getDefaultSettings() {
		return {}
	}

	return { data, setData, getDefaultSettings }
}
