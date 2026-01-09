import { createGlobalState } from '@vueuse/core'
import { useQuery } from '@urql/vue'

import { graphql } from '~/gql/gql'

export const useLanguagesApi = createGlobalState(() => {
  const useAvailableLanguages = () => useQuery({
    query: graphql(`
      query GetAvailableLanguages {
        moderationLanguagesAvailableLanguages {
          languages {
            iso_639_1
            name
          }
        }
      }
    `),
  })

  return {
    useAvailableLanguages,
  }
})
