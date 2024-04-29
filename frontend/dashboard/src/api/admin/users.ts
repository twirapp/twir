import { useQuery } from '@urql/vue'

import type { TwirUsersSearchParams, UsersGetAllQuery } from '@/gql/graphql'
import type { Ref } from 'vue'

import { useMutation } from '@/composables/use-mutation.js'
import { graphql } from '@/gql'

export type User = UsersGetAllQuery['twirUsers']['users'][0]

const invalidationKey = 'AdminUsersInvalidateKey'

export function useAdminUsers() {
	const useQueryUsers = (variables: Ref<TwirUsersSearchParams>) => useQuery({
		context: {
			additionalTypenames: [invalidationKey]
		},
		get variables() {
			return {
				opts: variables.value
			}
		},
		query: graphql(`
			query UsersGetAll($opts: TwirUsersSearchParams!) {
				twirUsers(opts: $opts) {
					total
					users {
						id
						isBanned
						isBotAdmin
						isBotEnabled
						twitchProfile {
							login
							displayName
							profileImageUrl
						}
					}
				}
			}
		`)
	})

	const useMutationUserSwitchBan = () => useMutation(graphql(`
		mutation UserSwitchBan($userId: ID!) {
			switchUserBan(userId: $userId)
		}
	`), [invalidationKey])

	const useMutationUserSwitchAdmin = () => useMutation(graphql(`
		mutation UserSwitchAdmin($userId: ID!) {
			switchUserAdmin(userId: $userId)
		}
	`), [invalidationKey])

	return {
		useQueryUsers,
		useMutationUserSwitchBan,
		useMutationUserSwitchAdmin
	}
}
