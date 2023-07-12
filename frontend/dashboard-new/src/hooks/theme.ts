import {useLocalStorage} from "@vueuse/core";

export const useTheme = () => useLocalStorage<'light' | 'dark'>('twirTheme', 'dark')
