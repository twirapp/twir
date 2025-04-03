import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import { graphql } from '@/gql'

export const useModerationAvailableLanguages = createGlobalState(() => {
	const query = () => useQuery({
		query: graphql(`
			query ModerationAvailableLanguages {
				moderationLanguagesAvailableLanguages {
					languages {
						name
						iso_639_1
						nativeName
					}
				}
			}
		`),
	})

	return {
		query,
	}
})
