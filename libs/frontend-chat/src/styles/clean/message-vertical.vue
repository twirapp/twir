<script setup lang="ts">

import MessageContent from '../../components/message-content.vue';
import { normalizeDisplayName } from '../../helpers';
import type { MessageComponentProps } from '../../types';

defineProps<MessageComponentProps>();
</script>

<template>
	<div class="message">
		<div class="profile">
			<div v-if="settings.showBadges && msg.badges" class="badges">
				<template
					v-for="(badgeValue, badgeName) of msg.badges"
					:key="badgeName+badgeValue"
				>
					<img
						v-if="settings.channelBadges.get(`${badgeName}-${badgeValue}`)"
						:src="settings.channelBadges.get(`${badgeName}-${badgeValue}`)!.image_url_4x"
						class="badge"
					/>

					<img
						v-else-if="settings.globalBadges.get(badgeName)?.versions.length"
						:src="settings.globalBadges.get(badgeName)!.versions.at(-1)!.image_url_4x"
						class="badge"
					/>
				</template>
			</div>
			<div v-if="msg.sender" class="username">
				{{ normalizeDisplayName(msg.sender!, msg.senderDisplayName!) }}{{ msg.isItalic ? '' : ':' }}
			</div>
		</div>
		<message-content
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
	gap: 8px;
	margin-right: 4px;
}

.badges {
	display: inline-flex;
	gap: 4px;
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
