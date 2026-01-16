import { ShortUrlProfileParamsSortByEnum } from '@twir/api/openapi'

export function useLinksPagination() {
	const currentPage = ref(0)
	const perPage = ref(10)
	const sortBy = ref<ShortUrlProfileParamsSortByEnum>(ShortUrlProfileParamsSortByEnum.Views)
	const total = ref(0)

	const totalPages = computed(() => Math.ceil(total.value / perPage.value))

	function goToPage(page: number) {
		if (page >= 0 && page < totalPages.value) {
			currentPage.value = page
		}
	}

	function nextPage() {
		if (currentPage.value < totalPages.value - 1) {
			currentPage.value++
		}
	}

	function previousPage() {
		if (currentPage.value > 0) {
			currentPage.value--
		}
	}

	function setSortBy(value: ShortUrlProfileParamsSortByEnum) {
		sortBy.value = value
		currentPage.value = 0 // Reset to first page when sorting changes
	}

	function setTotal(value: number) {
		total.value = value
	}

	return {
		currentPage: readonly(currentPage),
		perPage: readonly(perPage),
		sortBy: readonly(sortBy),
		total: readonly(total),
		totalPages,
		goToPage,
		nextPage,
		previousPage,
		setSortBy,
		setTotal,
	}
}
