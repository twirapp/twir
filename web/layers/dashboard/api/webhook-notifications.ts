import { useQuery } from '@urql/vue';
import { createGlobalState } from '@vueuse/core';

import { useMutation } from '~~/layers/dashboard/composables/use-mutation';
import { graphql } from '~/gql/gql.js';

const invalidationKey = 'WebhookNotifications';

export const useWebhookNotificationsApi = createGlobalState(() => {
	const useQueryWebhookNotifications = () =>
		useQuery({
			query: graphql(`
				query GetWebhookNotifications {
					webhookNotifications {
						id
						channelID
						enabled
						githubIssues
						githubPullRequests
						githubCommits
						githubIssuesOnline
						githubIssuesOffline
						githubPullRequestsOnline
						githubPullRequestsOffline
						githubCommitsOnline
						githubCommitsOffline
						createdAt
						updatedAt
					}
				}
			`),
			context: { additionalTypenames: [invalidationKey] },
		});

	const useMutationCreateWebhookNotifications = () =>
		useMutation(
			graphql(`
				mutation CreateWebhookNotifications($input: WebhookNotificationsCreateInput!) {
					webhookNotificationsCreate(input: $input) {
						id
						channelID
						enabled
						githubIssues
						githubPullRequests
						githubCommits
						githubIssuesOnline
						githubIssuesOffline
						githubPullRequestsOnline
						githubPullRequestsOffline
						githubCommitsOnline
						githubCommitsOffline
						createdAt
						updatedAt
					}
				}
			`),
			[invalidationKey],
		);

	const useMutationUpdateWebhookNotifications = () =>
		useMutation(
			graphql(`
				mutation UpdateWebhookNotifications(
					$id: String!
					$input: WebhookNotificationsUpdateInput!
				) {
					webhookNotificationsUpdate(id: $id, input: $input) {
						id
						channelID
						enabled
						githubIssues
						githubPullRequests
						githubCommits
						githubIssuesOnline
						githubIssuesOffline
						githubPullRequestsOnline
						githubPullRequestsOffline
						githubCommitsOnline
						githubCommitsOffline
						createdAt
						updatedAt
					}
				}
			`),
			[invalidationKey],
		);

	return {
		useQueryWebhookNotifications,
		useMutationCreateWebhookNotifications,
		useMutationUpdateWebhookNotifications,
	};
});
