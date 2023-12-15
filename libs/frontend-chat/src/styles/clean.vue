<script lang="ts" setup>
import { computed } from 'vue';

import MessageContent from '../components/messageContent.vue';
import { getColorFromMsg, getMessageAlign, normalizeDisplayName } from '../helpers.js';
import type { Settings, Message } from '../types.js';

const props = defineProps<{
	msg: Message,
	settings: Settings
}>();

const messageAlign = computed(() => getMessageAlign(props.settings.direction));
const messageFlexWrap = computed(() => messageAlign.value === 'center' ? 'nowrap' : 'wrap');
const userColor = computed(() => getColorFromMsg(props.msg));
</script>

<template>
	<div class="message">
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
		<div v-if="msg.sender" class="profile" :style="{ color: userColor, fontWeight: 700 }">
			{{ normalizeDisplayName(msg.sender!, msg.senderDisplayName!) }}{{ msg.isItalic ? '' : ':' }}
		</div>
		<message-content
			:chunks="msg.chunks"
			:is-italic="msg.isItalic"
			:text-shadow-color="settings.textShadowColor"
			:text-shadow-size="settings.textShadowSize"
			:message-align="messageAlign"
			:user-color="userColor"
		/>
	</div>
</template>

<style scoped>
.message {
	display: inline-flex;
	align-items: center;
}

.message .badges {
	display: inline-flex;
	gap: 4px;
	margin-right: 4px;
}

.message .badges .badge {
	height: 1em;
	width: 1em;
}

.message .profile {
	display: inline-flex;
}

.message .text {
	margin-left: 4px;
}

.message .text .emote {
	height: 1em;
	width: 1em;
}

.message > .text {
	flex-wrap: v-bind(messageFlexWrap);
}
</style>
