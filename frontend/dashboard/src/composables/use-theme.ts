import { createGlobalState, useColorMode } from '@vueuse/core'

export type Theme = 'light' | 'dark'

function _useTheme(key?: string) {
	return createGlobalState(() => {
		const mode = useColorMode({
			storage: localStorage,
			storageKey: key,
			initialValue: 'dark',
		})

		return {
			theme: mode,
			toggleTheme: () => {
				mode.value = mode.value === 'light' ? 'dark' : 'light'
			},
			changeTheme: (newTheme: Theme) => {
				mode.value = newTheme
			},
		}
	})
}

export const useTheme = _useTheme('twirTheme')
export const useThemeTwitchChat = _useTheme('twirTwitchChatTheme')
