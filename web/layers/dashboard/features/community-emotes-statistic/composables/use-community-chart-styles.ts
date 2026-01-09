import { useTheme } from '#layers/dashboard/composables/use-theme'
import { createGlobalState } from '@vueuse/core'
import { computed } from 'vue'

export const useCommunityChartStyles = createGlobalState(() => {
	const { theme } = storeToRefs(useTheme())

	const chartStyles = computed(() => {
		const isDark = theme.value.value === 'dark'
		const styles = getComputedStyle(document.documentElement)

		const textColor = `hsl(${styles.getPropertyValue('--foreground')})`
		const borderColor = `hsl(${styles.getPropertyValue('--border')})`

		return {
			isDark,
			textColor,
			borderColor,
		}
	})

	return {
		chartStyles,
	}
})
