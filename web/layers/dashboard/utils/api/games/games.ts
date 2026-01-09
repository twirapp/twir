import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import { useMutation } from '~/composables/use-mutation'
import { graphql } from '~/gql'

const gamesInvalidationKey = 'gamesInvalidationKey'

export const useGamesApi = createGlobalState(() => {
	const useGamesQuery = () => useQuery({
		query: graphql(`
			query Games {
				gamesEightBall {
					answers
					enabled
				}
				gamesDuel {
					enabled
					userCooldown
					globalCooldown
					timeoutSeconds
					startMessage
					resultMessage
					secondsToAccept
					pointsPerWin
					pointsPerLose
					bothDiePercent
					bothDieMessage
				}
				gamesRussianRoulette {
					enabled
					canBeUsedByModerator
					timeoutSeconds
					decisionSeconds
					initMessage
					surviveMessage
					deathMessage
					chargedBullets
					tumberSize
				}
				gamesSeppuku {
					enabled
					message
					messageModerators
					timeoutModerators
					timeoutSeconds
				}
				gamesVoteban {
					enabled
					timeoutSeconds
					timeoutModerators
					initMessage
					banMessage
					banMessageModerators
					surviveMessage
					surviveMessageModerators
					neededVotes
					voteDuration
					votingMode
					chatVotesWordsPositive
					chatVotesWordsNegative
				}
			}
		`),
		variables: {},
		context: {
			additionalTypenames: [gamesInvalidationKey],
		},
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
		[gamesInvalidationKey],
	)

	const useDuelMutation = () => useMutation(
		graphql(`
			mutation UpdateDuelSettings($opts: DuelGameOpts!) {
				gamesDuelUpdate(opts: $opts) {
					bothDieMessage
				}
			}
		`),
		[gamesInvalidationKey],
	)

	const useRussianRouletteMutation = () => useMutation(
		graphql(`
			mutation UpdateRussianRouletteSettings($opts: RussianRouletteGameOpts!) {
				gamesRussianRouletteUpdate(opts: $opts) {
					chargedBullets
				}
			}
		`),
		[gamesInvalidationKey],
	)

	const useSeppukuMutation = () => useMutation(
		graphql(`
			mutation UpdateSeppukuSettings($opts: SeppukuGameOpts!) {
				gamesSeppukuUpdate(opts: $opts) {
					message
				}
			}
		`),
		[gamesInvalidationKey],
	)

	const useVotebanMutation = () => useMutation(
		graphql(`
			mutation UpdateVotebanSettings($opts: VotebanGameOpts!) {
				gamesVotebanUpdate(opts: $opts) {
					timeoutSeconds
				}
			}
		`),
		[gamesInvalidationKey],
	)

	return {
		useGamesQuery,
		useEightBallMutation,
		useDuelMutation,
		useRussianRouletteMutation,
		useSeppukuMutation,
		useVotebanMutation,
	}
})
