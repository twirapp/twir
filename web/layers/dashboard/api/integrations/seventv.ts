import { useMutation, useSubscription } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import { graphql } from '@/gql'

export const useSevenTvIntegration = createGlobalState(() => {
	const subscription = useSubscription({
		query: graphql(`
			subscription SevenTvData {
				sevenTvData {
					botSeventvProfile {
						id
						username
						displayName
						avatarUri
					}
					userSeventvProfile {
						id
						username
						displayName
						avatarUri
					}
					deleteEmotesOnlyAddedByApp
					emoteSetId
					isEditor
					rewardIdForAddEmote
					rewardIdForRemoveEmote
				}
			}
		`),
	})

	const updater = useMutation(graphql(`
		mutation SevenTvUpdate($input: SevenTvUpdateInput!) {
			sevenTvUpdate(input: $input)
		}
	`))

	return {
		subscription,
		updater,
	}
})
