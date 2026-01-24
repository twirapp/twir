import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import { useMutation } from '@/composables/use-mutation'
import { graphql } from '@/gql/gql'

const invalidationKey = 'WebhookNotifications'

export const useWebhookNotificationsApi = createGlobalState(() => {
	const useQueryWebhookNotifications = () => useQuery({
		query: graphql(`
      query GetWebhookNotifications {
        webhookNotifications {
          id
          channelID
          enabled
          githubIssues
          githubPullRequests
          githubCommits
          createdAt
          updatedAt
        }
      }
    `),
		context: { additionalTypenames: [invalidationKey] },
	})

	const useMutationCreateWebhookNotifications = () => useMutation(graphql(`
    mutation CreateWebhookNotifications($input: WebhookNotificationsCreateInput!) {
      webhookNotificationsCreate(input: $input) {
        id
        channelID
        enabled
        githubIssues
        githubPullRequests
        githubCommits
        createdAt
        updatedAt
      }
    }
  `), [invalidationKey])

	const useMutationUpdateWebhookNotifications = () => useMutation(graphql(`
    mutation UpdateWebhookNotifications($id: String!, $input: WebhookNotificationsUpdateInput!) {
      webhookNotificationsUpdate(id: $id, input: $input) {
        id
        channelID
        enabled
        githubIssues
        githubPullRequests
        githubCommits
        createdAt
        updatedAt
      }
    }
  `), [invalidationKey])

	return {
		useQueryWebhookNotifications,
		useMutationCreateWebhookNotifications,
		useMutationUpdateWebhookNotifications,
	}
})
