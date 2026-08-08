import { useMutation, useQuery } from '@urql/vue'
import { type MaybeRef, computed, unref } from 'vue'

import { graphql } from '~/gql/gql.js'

const _channelOverlayLayerSettingsFragment = graphql(`
	fragment ChannelOverlayLayerSettingsFields on ChannelOverlayLayerSettings {
		htmlOverlayHtml
		htmlOverlayCss
		htmlOverlayJs
		htmlOverlayDataPollSecondsInterval
		imageUrl
		textContent
		textFontFamily
		textFontSize
		textFontWeight
		textColor
		textAlign
		videoUrl
		videoLoop
		videoMuted
		iframeUrl
		iframeScale
		widgetKey
		youtubeVideoId
		youtubeAutoplay
		youtubeLoop
		youtubeMuted
		emoteUrl
		emoteName
		emoteProvider
	}
`)

const channelOverlaysQuery = graphql(`
	query ChannelOverlays {
		channelOverlays {
			id
			channelId
			name
			createdAt
			updatedAt
			width
			height
			instaSave
			layers {
				id
				type
				name
			settings {
					...ChannelOverlayLayerSettingsFields
				}
				overlayId
				posX
				posY
				width
				height
				rotation
				createdAt
				updatedAt
				periodicallyRefetchData
				locked
				visible
				opacity
				zIndex
			}
		}
	}
`)

const channelOverlayByIdQuery = graphql(`
	query ChannelOverlayById($id: UUID!) {
		channelOverlayById(id: $id) {
			id
			channelId
			name
			createdAt
			updatedAt
			width
			height
			instaSave
			layers {
				id
				type
				name
				settings {
					...ChannelOverlayLayerSettingsFields
				}
				overlayId
				posX
				posY
				width
				height
				rotation
				createdAt
				updatedAt
				periodicallyRefetchData
				locked
				visible
				opacity
				zIndex
			}
		}
	}
`)

const channelOverlayCreateMutation = graphql(`
	mutation ChannelOverlayCreate($input: ChannelOverlayCreateInput!) {
		channelOverlayCreate(input: $input) {
			id
			channelId
			name
			createdAt
			updatedAt
			width
			height
			instaSave
			layers {
				id
				type
				name
				settings {
					...ChannelOverlayLayerSettingsFields
				}
				overlayId
				posX
				posY
				width
				height
				rotation
				createdAt
				updatedAt
				periodicallyRefetchData
				locked
				visible
				zIndex
			}
		}
	}
`)

const channelOverlayUpdateMutation = graphql(`
	mutation ChannelOverlayUpdate($id: UUID!, $input: ChannelOverlayUpdateInput!) {
		channelOverlayUpdate(id: $id, input: $input) {
			id
			channelId
			name
			createdAt
			updatedAt
			width
			height
			instaSave
			layers {
				id
				type
				name
				settings {
					...ChannelOverlayLayerSettingsFields
				}
				overlayId
				posX
				posY
				width
				height
				rotation
				createdAt
				updatedAt
				periodicallyRefetchData
				locked
				visible
				zIndex
			}
		}
	}
`)

const channelOverlayDeleteMutation = graphql(`
	mutation ChannelOverlayDelete($id: UUID!) {
		channelOverlayDelete(id: $id)
	}
`)

const channelOverlayParseHtmlMutation = graphql(`
	mutation ChannelOverlayParseHtml($html: String!) {
		channelOverlayParseHtml(html: $html)
	}
`)

export const useChannelOverlaysQuery = () => {
	return useQuery({
		query: channelOverlaysQuery,
	})
}

export const useChannelOverlayByIdQuery = (id: MaybeRef<string>) => {
	return useQuery({
		query: channelOverlayByIdQuery,
		variables: computed(() => ({
			id: unref(id),
		})),
		pause: computed(() => !unref(id)),
	})
}

export const useChannelOverlayCreate = () => useMutation(channelOverlayCreateMutation)

export const useChannelOverlayUpdate = () => useMutation(channelOverlayUpdateMutation)

export const useChannelOverlayDelete = () => useMutation(channelOverlayDeleteMutation)

export const useChannelOverlayParseHtml = () => useMutation(channelOverlayParseHtmlMutation)
