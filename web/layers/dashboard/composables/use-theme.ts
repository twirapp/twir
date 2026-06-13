import { createGlobalState, useColorMode } from '@vueuse/core'
import { computed } from 'vue'

export type Theme = 'light' | 'dark'

function _useTheme(key?: string) {
	return createGlobalState(() => {
		const theme = useColorMode({
			storage: localStorage,
			storageKey: key,
			initialValue: 'dark',
		})

		function toggleTheme() {
			theme.value = theme.value === 'light' ? 'dark' : 'light'
		}

		function changeTheme(newTheme: Theme) {
			theme.value = newTheme
		}

		const isDark = computed(() => theme.value === 'dark')

		return {
			theme,
			isDark,
			toggleTheme,
			changeTheme,
		}
	})
}

export const useTheme = _useTheme('twirTheme')
export const useThemeTwitchChat = _useTheme('twirTwitchChatTheme')
