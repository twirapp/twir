import { useLocalStorage } from '@vueuse/core';

export type Theme = 'light' | 'dark'

export const useTheme = () => {
	const theme = useLocalStorage<Theme>('twirTheme', 'dark');

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
