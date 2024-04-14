import { useQuery } from '@urql/vue';
import type { Ref } from 'vue';

import { useMutation } from '@/composables/use-mutation.js';
import { graphql } from '@/gql';
import type { TwirUsersSearchParams } from '@/gql/graphql';

const invalidationKey = 'AdminUsersInvalidateKey';

export const useAdminUsers = () => {
	const useQueryUsers = (variables: Ref<TwirUsersSearchParams>) => useQuery({
		context: {
			additionalTypenames: [invalidationKey],
		},
		get variables() {
			return {
				opts: variables.value,
			};
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
		`),
	});

	const useMutationUserSwitchBan = () => useMutation(graphql(`
		mutation UserSwitchBan($userId: ID!) {
			switchUserBan(userId: $userId)
		}
	`), [invalidationKey]);

	const useMutationUserSwitchAdmin = () => useMutation(graphql(`
		mutation UserSwitchAdmin($userId: ID!) {
			switchUserAdmin(userId: $userId)
		}
	`), [invalidationKey]);

	return {
		useQueryUsers,
		useMutationUserSwitchBan,
		useMutationUserSwitchAdmin,
	};
};
