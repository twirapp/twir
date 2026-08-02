import { useQuery } from '@urql/vue'

import { graphql } from '~/gql/gql.js'

const channelApiKeyQuery = graphql(`
	query McpGuideChannelApiKey {
		channelApiKey
	}
`)

export function useMcpChannelApiKey() {
	return useQuery({
		query: channelApiKeyQuery,
		variables: {},
		context: {
			key: 'McpGuideChannelApiKey',
		},
	})
}
