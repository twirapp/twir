import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import { graphql } from '@/gql/gql'
import { useMutation } from '~/composables/use-mutation'

const invalidationKey = 'ChatTranslation'

export const useChatTranslationApi = createGlobalState(() => {
	const useQueryChatTranslation = () =>
		useQuery({
			query: graphql(`
      query GetChatTranslation {
        chatTranslation {
          id
          channelID
          enabled
          targetLanguage
          excludedLanguages
          useItalic
          excludedUsersIDs
          createdAt
          updatedAt
        }
      }
    `),
			context: { additionalTypenames: [invalidationKey] },
		})

	const useMutationCreateChatTranslation = () =>
		useMutation(
			graphql(`
    mutation CreateChatTranslation($input: ChatTranslationCreateInput!) {
      chatTranslationCreate(input: $input) {
        id
        channelID
        enabled
        targetLanguage
        excludedLanguages
        useItalic
        excludedUsersIDs
        createdAt
        updatedAt
      }
    }
  `),
			[invalidationKey]
		)

	const useMutationUpdateChatTranslation = () =>
		useMutation(
			graphql(`
    mutation UpdateChatTranslation($id: String!, $input: ChatTranslationUpdateInput!) {
      chatTranslationUpdate(id: $id, input: $input) {
        id
        channelID
        enabled
        targetLanguage
        excludedLanguages
        useItalic
        excludedUsersIDs
        createdAt
        updatedAt
      }
    }
  `),
			[invalidationKey]
		)

	return {
		useQueryChatTranslation,
		useMutationCreateChatTranslation,
		useMutationUpdateChatTranslation,
	}
})
