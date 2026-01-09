import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import type { DudesOverlaySettings, DudesOverlaySettingsInput } from '@/gql/graphql'

import { graphql } from '@/gql'
import { useMutation } from '~/composables/use-mutation'

export const useDudesOverlayApi = createGlobalState(() => {
	const cacheKey = ['dudesOverlay']

	const useDudesQuery = () =>
		useQuery({
			query: graphql(`
				query DudesOverlays {
					dudesGetAll {
						id
						dudeSettings {
							color
							eyesColor
							cosmeticsColor
							maxLifeTime
							gravity
							scale
							soundsEnabled
							soundsVolume
							visibleName
							growTime
							growMaxScale
							maxOnScreen
							defaultSprite
						}
						messageBoxSettings {
							enabled
							borderRadius
							boxColor
							fontFamily
							fontSize
							padding
							showTime
							fill
						}
						nameBoxSettings {
							fontFamily
							fontSize
							fill
							lineJoin
							strokeThickness
							stroke
							fillGradientStops
							fillGradientType
							fontStyle
							fontVariant
							fontWeight
							dropShadow
							dropShadowAlpha
							dropShadowAngle
							dropShadowBlur
							dropShadowDistance
							dropShadowColor
						}
						ignoreSettings {
							ignoreCommands
							ignoreUsers
							users
						}
						spitterEmoteSettings {
							enabled
						}
					}
				}
			`),
			context: {
				additionalTypenames: cacheKey,
			},
			variables: {},
		})

	const useDudesByIdQuery = (id: string) =>
		useQuery({
			query: graphql(`
				query DudesOverlayById($id: UUID!) {
					dudesGetById(id: $id) {
						id
						dudeSettings {
							color
							eyesColor
							cosmeticsColor
							maxLifeTime
							gravity
							scale
							soundsEnabled
							soundsVolume
							visibleName
							growTime
							growMaxScale
							maxOnScreen
							defaultSprite
						}
						messageBoxSettings {
							enabled
							borderRadius
							boxColor
							fontFamily
							fontSize
							padding
							showTime
							fill
						}
						nameBoxSettings {
							fontFamily
							fontSize
							fill
							lineJoin
							strokeThickness
							stroke
							fillGradientStops
							fillGradientType
							fontStyle
							fontVariant
							fontWeight
							dropShadow
							dropShadowAlpha
							dropShadowAngle
							dropShadowBlur
							dropShadowDistance
							dropShadowColor
						}
						ignoreSettings {
							ignoreCommands
							ignoreUsers
							users
						}
						spitterEmoteSettings {
							enabled
						}
					}
				}
			`),
			context: {
				additionalTypenames: cacheKey,
			},
			variables: { id },
		})

	const useDudesCreate = () =>
		useMutation(
			graphql(`
				mutation DudesOverlayCreate($input: DudesOverlaySettingsInput!) {
					dudesCreate(input: $input)
				}
			`),
			cacheKey
		)

	const useDudesUpdate = () =>
		useMutation(
			graphql(`
				mutation DudesOverlayUpdate($id: UUID!, $input: DudesOverlaySettingsInput!) {
					dudesUpdate(id: $id, input: $input)
				}
			`),
			cacheKey
		)

	const useDudesDelete = () =>
		useMutation(
			graphql(`
				mutation DudesOverlayDelete($id: UUID!) {
					dudesDelete(id: $id)
				}
			`),
			cacheKey
		)

	return {
		useDudesQuery,
		useDudesByIdQuery,
		useDudesCreate,
		useDudesUpdate,
		useDudesDelete,
	}
})

export const useDudesOverlayManager = () => {
	const api = useDudesOverlayApi()

	return {
		useGet: (id: string) => api.useDudesByIdQuery(id),
		useGetAll: () => api.useDudesQuery(),
		useCreate: () => api.useDudesCreate(),
		useUpdate: () => api.useDudesUpdate(),
		useDelete: () => api.useDudesDelete(),
	}
}

export type { DudesOverlaySettings, DudesOverlaySettingsInput }
