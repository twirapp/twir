import { ref } from 'vue'

import type { BadgeVersion, ChatBadge } from '../types'

interface Badge {
	id: string
	versions: Record<string, BadgeVersion>
}

const globalMappedBadges = ref<Record<string, Badge>>({})
const channelMappedBadges = ref<Record<string, Badge>>({})

export function useMappedBadges() {
	function setGlobalBadges(badges: ChatBadge[]) {
		globalMappedBadges.value = badges.reduce((acc, badge) => {
			acc[badge.set_id] = {
				id: badge.set_id,
				versions: badge.versions.reduce((acc, version) => {
					acc[version.id] = version
					return acc
				}, {} as Record<string, BadgeVersion>),
			}
			return acc
		}, {} as Record<string, Badge>)
	}

	function setChannelBadges(badges: ChatBadge[]) {
		channelMappedBadges.value = badges.reduce((acc, badge) => {
			acc[badge.set_id] = {
				id: badge.set_id,
				versions: badge.versions.reduce((acc, version) => {
					acc[version.id] = version
					return acc
				}, {} as Record<string, BadgeVersion>),
			}
			return acc
		}, {} as Record<string, Badge>)
	}

	return {
		setGlobalBadges,
		setChannelBadges,

		globalMappedBadges,
		channelMappedBadges,
	}
}
