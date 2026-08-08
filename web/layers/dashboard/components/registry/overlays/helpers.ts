import { ChannelOverlayLayerType } from '~/gql/graphql.js'

export function convertOverlayLayerTypeToText(type: ChannelOverlayLayerType): string {
	switch (type) {
		case ChannelOverlayLayerType.Html:
			return 'HTML'
		case ChannelOverlayLayerType.Image:
			return 'Image'
		case ChannelOverlayLayerType.Text:
			return 'Text'
		case ChannelOverlayLayerType.Video:
			return 'Video'
		case ChannelOverlayLayerType.Iframe:
			return 'Widget'
		case ChannelOverlayLayerType.Youtube:
			return 'YouTube'
		case ChannelOverlayLayerType.Emote:
			return 'Emote'
		default:
			return 'UNKNOWN'
	}
}
