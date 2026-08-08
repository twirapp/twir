import { Clock, Eye, Heart, MessageSquare, Star } from 'lucide-vue-next'

import type { FunctionalComponent } from 'vue'

import type { StreamStatsCounterKey } from '@/composables/stream-stats/use-stream-stats-counters.js'
import { StreamStatsOverlayVariant } from '@/gql/graphql'

export const counterIcons: Record<StreamStatsCounterKey, FunctionalComponent> = {
	viewers: Eye,
	messages: MessageSquare,
	uptime: Clock,
	subscribers: Star,
	followers: Heart,
}

export const counterColors: Record<StreamStatsCounterKey, string> = {
	viewers: '#8b5cf6',
	messages: '#38bdf8',
	uptime: '#fbbf24',
	subscribers: '#f472b6',
	followers: '#34d399',
}

export function variantRootClass(variant: StreamStatsOverlayVariant | undefined): string {
	switch (variant) {
		case StreamStatsOverlayVariant.HorizontalCompact:
			return 'root--compact'
		case StreamStatsOverlayVariant.Vertical:
			return 'root--vertical'
		case StreamStatsOverlayVariant.VerticalCompact:
			return 'root--vertical root--compact'
		case StreamStatsOverlayVariant.Large:
			return 'root--large'
		default:
			return ''
	}
}
