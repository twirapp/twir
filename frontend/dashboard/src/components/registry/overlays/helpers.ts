import { ChannelOverlayLayerType } from '@/gql/graphql'

export function convertOverlayLayerTypeToText(type: ChannelOverlayLayerType): string {
	switch (type) {
		case ChannelOverlayLayerType.Html:
			return 'HTML'
		default:
			return 'UNKNOWN'
	}
}
