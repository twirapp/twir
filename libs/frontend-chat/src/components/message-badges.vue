<script setup lang="ts">
import KickBadgeIcon from './kick-badge-icon.vue'
import PlatformIcon from './platform-icon.vue'
import { useMappedBadges } from '../composables/mapped-badges.js'

import type { Message, Settings } from '../types.js'

const props = defineProps<{
	msg: Message
	settings: Settings
}>()

const { globalMappedBadges, channelMappedBadges } = useMappedBadges()

function shouldRender() {
	return (
		(props.settings.showPlatformIcon && props.msg.platform) ||
		props.settings.showBadges
	)
}
</script>

<template>
	<div v-if="shouldRender()" class="badges">
		<PlatformIcon
			v-if="settings.showPlatformIcon && msg.platform"
			:platform="msg.platform"
		/>
		<template v-if="settings.showBadges">
			<template v-if="msg.platform === 'kick'">
				<KickBadgeIcon
					v-for="badge of msg.kickBadges ?? []"
					:key="badge.type"
					:type="badge.type"
					:text="badge.text"
				/>
			</template>
			<template v-else-if="msg.badges">
				<template v-for="(badgeValue, badgeName) of msg.badges" :key="badgeName + badgeValue">
					<img
						v-if="channelMappedBadges[badgeName]?.versions[badgeValue]"
						:src="channelMappedBadges[badgeName]!.versions[badgeValue].image_url_4x"
						class="badge"
					/>

					<img
						v-else-if="globalMappedBadges[badgeName]?.versions[badgeValue]"
						:src="globalMappedBadges[badgeName]!.versions[badgeValue].image_url_4x"
						class="badge"
					/>
				</template>
			</template>
		</template>
	</div>
</template>

<style scoped>
.badges {
	display: inline-flex;
	gap: 4px;
	align-items: center;
	margin-right: var(--message-badges-margin, 4px);
}

.badge {
	height: 1em;
	width: 1em;
}
</style>
