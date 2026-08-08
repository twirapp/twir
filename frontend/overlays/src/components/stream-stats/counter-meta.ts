import { Clock, Eye, Heart, MessageSquare, Star } from 'lucide-vue-next'

import type { FunctionalComponent } from 'vue'

import type { StreamStatsCounterKey } from '@/composables/stream-stats/use-stream-stats-counters.js'

export const counterIcons: Record<StreamStatsCounterKey, FunctionalComponent> = {
	viewers: Eye,
	messages: MessageSquare,
	uptime: Clock,
	subscribers: Star,
	followers: Heart,
}
