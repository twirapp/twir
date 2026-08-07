import { ChannelOverlayLayerType } from '~/gql/graphql.js'

export interface LayerTypeMeta {
	icon: string
	labelKey: string
	chipClass: string
}

// Chip colors follow the approved sidebar mockup (.omo/mockups/overlay-builder-sidebar.html, design 1)
export const LAYER_TYPE_META: Record<ChannelOverlayLayerType, LayerTypeMeta> = {
	[ChannelOverlayLayerType.Iframe]: {
		icon: 'lucide:panels-top-left',
		labelKey: 'overlayBuilder.layerTypes.widget',
		chipClass: 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400',
	},
	[ChannelOverlayLayerType.Text]: {
		icon: 'lucide:type',
		labelKey: 'overlayBuilder.layerTypes.text',
		chipClass: 'bg-sky-500/10 text-sky-600 dark:text-sky-400',
	},
	[ChannelOverlayLayerType.Emote]: {
		icon: 'lucide:smile',
		labelKey: 'overlayBuilder.layerTypes.emote',
		chipClass: 'bg-amber-500/10 text-amber-600 dark:text-amber-400',
	},
	[ChannelOverlayLayerType.Youtube]: {
		icon: 'simple-icons:youtube',
		labelKey: 'overlayBuilder.layerTypes.youtube',
		chipClass: 'bg-red-500/10 text-red-600 dark:text-red-400',
	},
	[ChannelOverlayLayerType.Video]: {
		icon: 'lucide:video',
		labelKey: 'overlayBuilder.layerTypes.video',
		chipClass: 'bg-violet-500/10 text-violet-600 dark:text-violet-400',
	},
	[ChannelOverlayLayerType.Image]: {
		icon: 'lucide:image',
		labelKey: 'overlayBuilder.layerTypes.image',
		chipClass: 'bg-lime-500/10 text-lime-600 dark:text-lime-400',
	},
	[ChannelOverlayLayerType.Html]: {
		icon: 'lucide:code-xml',
		labelKey: 'overlayBuilder.layerTypes.html',
		chipClass: 'bg-orange-500/10 text-orange-600 dark:text-orange-400',
	},
}

export function getLayerTypeMeta(type: ChannelOverlayLayerType): LayerTypeMeta {
	return LAYER_TYPE_META[type]
}
