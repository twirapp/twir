import { useLocalStorage } from '@vueuse/core';

export type Theme = 'light' | 'dark'

export const useTheme = (key?: string) => {
	const theme = useLocalStorage<Theme>(key ?? 'twirTheme', 'dark');

	return {
		theme,
		toggleTheme: () => {
			theme.value = theme.value === 'light' ? 'dark' : 'light';
		},
		changeTheme: (newTheme: Theme) => {
			theme.value = newTheme;
		},
	};
};
