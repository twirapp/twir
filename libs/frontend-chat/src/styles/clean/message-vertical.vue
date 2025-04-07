<script setup lang="ts">
import MessageContent from '../../components/message-content.vue'
import { useMappedBadges } from '../../composables/mapped-badges'
import { normalizeDisplayName } from '../../helpers'

import type { MessageComponentProps } from '../../types'

defineProps<MessageComponentProps>()

const { globalMappedBadges, channelMappedBadges } = useMappedBadges()
</script>

<template>
	<div class="message">
		<div class="profile">
			<div v-if="settings.showBadges && msg.badges" class="badges">
				<template
					v-for="(badgeValue, badgeName) of msg.badges"
					:key="badgeName + badgeValue"
				>
					<img
						v-if="channelMappedBadges[badgeName]?.versions[badgeValue]"
						:src="channelMappedBadges[badgeName]!.versions[badgeValue].image_url_4x"
						class="badge"
					/>

					<img
						v-if="globalMappedBadges[badgeName]?.versions[badgeValue]"
						:src="globalMappedBadges[badgeName]!.versions[badgeValue].image_url_4x"
						class="badge"
					/>
				</template>
			</div>
			<div v-if="msg.sender" class="username">
				{{ normalizeDisplayName(msg.sender!, msg.senderDisplayName!) }}
			</div>
			<span v-if="msg.sender">
				{{ msg.isItalic ? '' : ':' }}
			</span>
		</div>
		<MessageContent
			:chunks="msg.chunks"
			:is-italic="msg.isItalic"
			:text-shadow-color="settings.textShadowColor"
			:text-shadow-size="settings.textShadowSize"
			:user-color="userColor"
		/>
	</div>
</template>

<style scoped>
.message {
	display: block;
	width: 100%;
	margin-left: 0.2em;
}

.profile {
	display: inline-flex;
	margin-right: 4px;
}

.badges {
	display: flex;
	gap: 4px;
	align-items: center;
	margin-right: 4px;
}

.badge {
	height: 1em;
	width: 1em;
}

.username {
	color: v-bind(userColor);
	font-weight: 700;
}
</style>
