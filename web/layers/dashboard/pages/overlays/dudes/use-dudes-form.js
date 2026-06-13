import { ref } from 'vue'

export function useDudesForm() {
	const data = ref(null)

	function setData(entity: any) {
		data.value = entity
	}

	function getDefaultSettings() {
		return { id: null }
	}

	return { data, setData, getDefaultSettings }
}
