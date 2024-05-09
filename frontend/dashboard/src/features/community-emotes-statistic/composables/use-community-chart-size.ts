import { ref } from 'vue'

const chartSizes = ref({
	width: 0,
	height: 0,
})

export function useCommunityChartSize() {
	function setChartSize(width: number, height: number) {
		chartSizes.value.width = width
		chartSizes.value.height = height
	}

	return {
		chartSizes,
		setChartSize,
	}
}
