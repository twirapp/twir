export function useUploadsPagination() {
	const currentPage = ref(0)
	const perPage = ref(10)
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

	function setTotal(value: number) {
		total.value = value
	}

	return {
		currentPage: readonly(currentPage),
		perPage: readonly(perPage),
		total: readonly(total),
		totalPages,
		goToPage,
		nextPage,
		previousPage,
		setTotal,
	}
}
