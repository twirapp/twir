import { createGlobalState } from '@vueuse/core'
import { computed } from 'vue'

import { useTheme } from '@/composables/use-theme.js'

export const useCommunityChartStyles = createGlobalState(() => {
	const { theme } = useTheme()

	const chartStyles = computed(() => {
		const isDark = theme.value === 'dark'
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
