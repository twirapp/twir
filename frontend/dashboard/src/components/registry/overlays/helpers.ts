import { OverlayLayerType } from '@twir/api/messages/overlays/overlays';

export const convertOverlayLayerTypeToText = (type: OverlayLayerType) => {
	switch (type) {
		case OverlayLayerType.HTML:
			return 'HTML';
		default:
			return 'UNKNOWN';
	}
};
