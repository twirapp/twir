import { useQuery as useUrqlQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import { useMutation } from '~~/layers/dashboard/composables/use-mutation.js'
import { graphql } from '~/gql/gql.js'

export const channelPlatformsInvalidationKey = 'ChannelPlatforms'

export const useChannelPlatformsApi = createGlobalState(() => {
	const useQuery = () =>
		useUrqlQuery({
			query: graphql(`
				query ChannelPlatforms {
					channelPlatformBindings {
						id
						platform
						userId
						platformChannelId
						enabled
						platformUserId
						platformLogin
						platformDisplayName
						platformAvatar
						capabilities {
							name
						}
					}
					channelPlatformOptions {
						platform
						capabilities {
							name
						}
					}
				}
			`),
			context: {
				additionalTypenames: [channelPlatformsInvalidationKey],
			},
		})

	const useConnect = () =>
		useMutation(
			graphql(`
				mutation ChannelPlatformConnect($platform: Platform!) {
					channelPlatformConnect(platform: $platform)
				}
			`),
			[channelPlatformsInvalidationKey],
		)

	const useDisconnect = () =>
		useMutation(
			graphql(`
				mutation ChannelPlatformDisconnect($platform: Platform!) {
					channelPlatformDisconnect(platform: $platform)
				}
			`),
			[channelPlatformsInvalidationKey],
		)

	const useSetEnabled = () =>
		useMutation(
			graphql(`
				mutation ChannelPlatformSetEnabled($platform: Platform!, $enabled: Boolean!) {
					channelPlatformSetEnabled(platform: $platform, enabled: $enabled) {
						id
						platform
						enabled
					}
				}
			`),
			[channelPlatformsInvalidationKey],
		)

	return {
		useQuery,
		useConnect,
		useDisconnect,
		useSetEnabled,
	}
})
