import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { computed } from 'vue'

import type { Platform } from '~/gql/graphql.js'

import { graphql } from '~/gql/gql.js'

export const useDefaultCommandsApi = createGlobalState(() => {
	const query = useQuery({
		variables: {},
		query: graphql(`
			query GetDefaultCommandsInfo {
				commandsDefault {
					name
					description
					module
					aliases
					platforms
				}
			}
		`),
	})

	const defaultCommands = computed(() => query.data.value?.commandsDefault ?? [])

	const defaultCommandPlatforms = computed(() => {
		const platformsByName = new Map<string, Platform[]>()
		for (const command of defaultCommands.value) {
			platformsByName.set(command.name, command.platforms)
		}
		return platformsByName
	})

	return {
		defaultCommands,
		defaultCommandPlatforms,
	}
})
