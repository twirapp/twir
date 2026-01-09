import { useMutation, useQuery } from '@urql/vue'
import { type MaybeRef, computed, unref } from 'vue'

import { graphql } from '@/gql'

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
				settings {
					htmlOverlayHtml
					htmlOverlayCss
					htmlOverlayJs
					htmlOverlayDataPollSecondsInterval
					imageUrl
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
				settings {
					htmlOverlayHtml
					htmlOverlayCss
					htmlOverlayJs
					htmlOverlayDataPollSecondsInterval
					imageUrl
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
				settings {
					htmlOverlayHtml
					htmlOverlayCss
					htmlOverlayJs
					htmlOverlayDataPollSecondsInterval
					imageUrl
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
				settings {
					htmlOverlayHtml
					htmlOverlayCss
					htmlOverlayJs
					htmlOverlayDataPollSecondsInterval
					imageUrl
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
