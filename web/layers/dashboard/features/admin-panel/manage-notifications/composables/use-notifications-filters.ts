import { createGlobalState, refDebounced } from '@vueuse/core'
import { ref } from 'vue'

import { NotificationType } from '~/gql/graphql'

export const useNotificationsFilters = createGlobalState(() => {
	const searchInput = ref('')
	const debounceSearchInput = refDebounced<string>(searchInput, 500)

	const filterInput = ref<NotificationType>(NotificationType.Global)

	return {
		searchInput,
		debounceSearchInput,
		filterInput,
	}
})
