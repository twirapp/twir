interface UseAppSeoOptions {
	titleKey: string
	descriptionKey: string
}

export function useAppSeo(options: UseAppSeoOptions) {
	const { t } = useI18n()

	const title = computed(() => t(options.titleKey))
	const description = computed(() => t(options.descriptionKey))

	useSeoMeta({
		title,
		description,
		ogTitle: title,
		ogDescription: description,
		ogType: 'website',
		twitterCard: 'summary_large_image',
		twitterTitle: title,
		twitterDescription: description,
	})
}
