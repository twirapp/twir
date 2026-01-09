import { computed } from 'vue'

export type Theme = 'light' | 'dark'

export const useTheme = defineStore('theme', () => {
	const theme = useColorMode()

	function toggleTheme() {
		theme.preference = theme.value === 'light' ? 'dark' : 'light'
	}

	function changeTheme(newTheme: Theme) {
		theme.preference = newTheme
	}

	const isDark = computed(() => theme.value === 'dark')

	return {
		theme,
		isDark,
		toggleTheme,
		changeTheme,
	}
})
