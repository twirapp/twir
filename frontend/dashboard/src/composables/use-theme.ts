import { useColorMode } from '@vueuse/core';

export type Theme = 'light' | 'dark'

export const useTheme = (key?: string) => {
	const mode = useColorMode({
		storage: localStorage,
		storageKey: key ?? 'twirTheme',
		initialValue: 'dark',
	});

	return {
		theme: mode,
		toggleTheme: () => {
			mode.value = mode.value === 'light' ? 'dark' : 'light';
		},
		changeTheme: (newTheme: Theme) => {
			mode.value = newTheme;
		},
	};
};
