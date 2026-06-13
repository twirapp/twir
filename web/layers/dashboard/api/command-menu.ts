import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { computed } from 'vue'

import { graphql } from '@/gql'

export const commandMenuCacheKey = 'CommandMenuCacheKey'

const CommandMenuQuery = graphql(`
	query CommandMenuData {
		commands {
			id
			name
			description
			enabled
			module
		}
		keywords {
			id
			text
			enabled
		}
		variables {
			id
			name
			description
		}
	}
`)

export const useCommandMenuData = createGlobalState(() => {
	const query = useQuery({
		query: CommandMenuQuery,
		context: { additionalTypenames: [commandMenuCacheKey] },
		variables: {},
	})

	const commands = computed(() => query.data.value?.commands ?? [])
	const keywords = computed(() => query.data.value?.keywords ?? [])
	const variables = computed(() => query.data.value?.variables ?? [])

	const isLoading = computed(() => query.fetching.value)

	return {
		query,
		commands,
		keywords,
		variables,
		isLoading,
	}
})
