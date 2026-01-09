export function useDashboardLocale() {
	const { locale, setLocale } = useI18n()
	const savedLocale = useLocalStorage('twirLocale', 'en')

	// Sync with localStorage for backward compatibility
	watch(locale, (newLocale) => {
		savedLocale.value = newLocale
	})

	// Apply saved locale on mount
	onMounted(() => {
		if (savedLocale.value && savedLocale.value !== locale.value) {
			setLocale(savedLocale.value)
		}
	})

	return { locale, setLocale }
}
