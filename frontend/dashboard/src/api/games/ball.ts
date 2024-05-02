import { useQuery } from '@urql/vue'
import { defineStore } from 'pinia'

import { useMutation } from '@/composables/use-mutation'
import { graphql } from '@/gql'

const eightBallInvalidationKey = '8ballInvalidationKey'

export const useGamesApi = defineStore('games/8ball', () => {
	const useGamesQuery = () => useQuery({
		query: graphql(`
			query Games {
				gamesEightBall {
					answers
					enabled
				}
			}
		`),
		variables: {},
		context: {
			additionalTypenames: [eightBallInvalidationKey]
		}
	})

	const useEightBallMutation = () => useMutation(
		graphql(`
			mutation UpdateEightBallSettings($opts: EightBallGameOpts!) {
				gamesEightBallUpdate(opts: $opts) {
					answers
					enabled
				}
			}
		`),
		[eightBallInvalidationKey]
	)

	return {
		useGamesQuery,
		useEightBallMutation
	}
})

// export const use8ballSettings = () => useQuery({
// 	queryKey: ['8ballSettings'],
// 	queryFn: async () => {
// 		const req = await protectedApiClient.gamesGetEightBallSettings({});
// 		return req.response;
// 	},
// });
//
// export const use8ballUpdateSettings = () => {
// 	const queryClient = useQueryClient();
//
// 	return useMutation({
// 		mutationKey: ['8ballSettings'],
// 		mutationFn: async (opts: { answers: string[], enabled: boolean }) => {
// 			const req = await protectedApiClient.gamesUpdateEightBallSettings(opts);
// 			return req.response;
// 		},
// 		onSuccess: async () => {
// 			await queryClient.invalidateQueries(['8ballSettings']);
// 			await queryClient.invalidateQueries(['commands']);
// 		},
// 	});
// };
