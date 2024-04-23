import { useQuery } from '@urql/vue';
import { defineStore } from 'pinia';
import type { Ref } from 'vue';

import { useMutation } from '@/composables/use-mutation.js';
import { graphql } from '@/gql/gql.js';
import type { CommunityUsersOpts, GetAllCommunityUsersQuery } from '@/gql/graphql';

export type CommunityUser = GetAllCommunityUsersQuery['communityUsers']['users'][0]

const invalidationKey = 'CommunityInvalidateKey';

export const useCommunityUsersApi = defineStore('api/community-users', () => {
	const useCommunityUsers = (variables: Ref<CommunityUsersOpts>) => useQuery({
		context: {
			additionalTypenames: [invalidationKey],
		},
		get variables() {
			return {
				opts: variables.value,
			};
		},
		query: graphql(`
			query GetAllCommunityUsers($opts: CommunityUsersOpts!) {
				communityUsers(opts: $opts) {
					total
					users {
						id
						twitchProfile {
							login
							displayName
							profileImageUrl
						}
						watchedMs
						messages
						usedEmotes
						usedChannelPoints
					}
				}
			}
		`),
	});

	const useMutationCommunityReset = () => useMutation(graphql(`
		mutation CommunityReset($type: CommunityUsersResetType!) {
			communityResetStats(type: $type)
		}
	`), [invalidationKey]);

	return {
		useCommunityUsers,
		useMutationCommunityReset,
	};
});
