import { useMutation, useQuery } from '@urql/vue';

import { graphql } from '@/gql';

const invalidationKey = 'AdminBadgesInvalidateKey';
export const useAdminBadges = () => {

	const useMutationCreateBadge = () => useMutation(graphql(`
		mutation CreateBadge($opts: TwirBadgeCreateOpts!) {
			badgesCreate(opts: $opts) {
				id
			}
		}
	`));

	const useMutationDeleteBadge = () => useMutation(graphql(`
		mutation DeleteBadge($id: ID!) {
			badgesDelete(id: $id)
		}
	`));

	const useMutationUpdateBadge = () => useMutation(graphql(`
		mutation UpdateBadge($id: ID!, $opts: TwirBadgeUpdateOpts!) {
			badgesUpdate(id: $id, opts: $opts) {
				id
			}
		}
	`));

	const useMutationsAddUserBadge = () => useMutation(graphql(`
		mutation AddUserBadge($id: ID!, $userId: String!) {
			badgesAddUser(id: $id, userId: $userId)
		}
	`));

	const useMutationsRemoveUserBadge = () => useMutation(graphql(`
		mutation RemoveUserBadge($id: ID!, $userId: String!) {
			badgesRemoveUser(id: $id, userId: $userId)
		}
	`));

	return {
		useMutationCreateBadge,
		useMutationDeleteBadge,
		useMutationUpdateBadge,
		useMutationsAddUserBadge,
		useMutationsRemoveUserBadge,
	};
};

export const useQueryBadges = () => useQuery({
	context: {
		additionalTypenames: [invalidationKey],
	},
	query: graphql(`
		query BadgesGetAll {
			twirBadges {
				id
				name
				createdAt
				fileUrl
				enabled
				ffzSlot
				users
			}
		}
	`),
});
