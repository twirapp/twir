import { defineStore } from 'pinia'
import { computed } from 'vue'

import { useTheme } from '@/composables/use-theme.js'

export const useCommunityChartStyles = defineStore('features/community-chart-styles', () => {
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
