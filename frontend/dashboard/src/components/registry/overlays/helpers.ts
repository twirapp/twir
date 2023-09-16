import { OverlayLayerType } from '@twir/grpc/generated/api/api/overlays';

export const convertOverlayLayerTypeToText = (type: OverlayLayerType) => {
	switch (type) {
		case OverlayLayerType.HTML:
			return 'HTML';
		default:
			return 'UNKNOWN';
	}
};
