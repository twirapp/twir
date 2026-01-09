import { createGlobalState, refDebounced } from '@vueuse/core'
import { computed, ref } from 'vue'

import type { ChatMessageInput } from '~/gql/graphql.js'

export const useChatMessagesFilters = createGlobalState(() => {
	const userSearchInput = ref('')
	const debouncedUserSearchInput = refDebounced(userSearchInput, 500)

	const textSearchInput = ref('')
	const debouncedTextSearchInput = refDebounced(textSearchInput, 500)

	const page = ref(0)
	const perPage = ref(500)

	const computedFilters = computed<ChatMessageInput>(() => {
		return {
			page: 0,
			perPage: 500,
			userNameLike: debouncedUserSearchInput.value,
			textLike: debouncedTextSearchInput.value,
		}
	})

	return {
		userSearchInput,
		textSearchInput,
		computedFilters,
		page,
		perPage,
	}
})
