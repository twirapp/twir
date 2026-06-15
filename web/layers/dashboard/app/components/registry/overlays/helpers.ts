import { ChannelOverlayLayerType } from '~/gql/graphql.js'

export function convertOverlayLayerTypeToText(type: ChannelOverlayLayerType): string {
	switch (type) {
		case ChannelOverlayLayerType.Html:
			return 'HTML'
		case ChannelOverlayLayerType.Image:
			return 'Image'
		default:
			return 'UNKNOWN'
	}
}
