<script lang="ts" setup>
import { computed } from 'vue';

import MessageContent from '../components/messageContent.vue';
import { getColorFromMsg, getMessageAlign, normalizeDisplayName } from '../helpers.js';
import type { Direction } from '../helpers.js';
import type { Settings, Message } from '../types.js';

const props = defineProps<{
	msg: Message,
	settings: Settings,
	direction: Direction
}>();

const messageAlign = computed(() => getMessageAlign(props.settings.direction));
const messageFlexWrap = computed(() => props.direction === 'horizontal' ? 'nowrap' : 'wrap');
const messageDirection = computed(() => props.direction === 'horizontal' ? 'row' : 'column');
const messageWidth = computed(() => props.direction === 'vertical' ? '100%' : 'auto');
const profileDirection = computed(() => props.direction === 'vertical' ? 'row' : 'row-reverse');
const userColor = computed(() => getColorFromMsg(props.msg));
</script>

<template>
	<div class="message">
		<div class="profile">
			<div v-if="msg.sender" :style="{ color: userColor, fontWeight: 700 }">
				{{ normalizeDisplayName(msg.sender!, msg.senderDisplayName!) }}
			</div>
			<div v-if="settings.showBadges" class="badges">
				<span v-if="settings.showAnnounceBadge && msg.isAnnounce" class="text-badge">Announce</span>
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
	display: flex;
	flex-direction: v-bind(messageDirection);
	align-items: v-bind(messageAlign);
	padding: 0.5em;
	gap: 0.2em;
	border-radius: 8px;
	background-color: #252525;
}

.message .badges {
	display: inline-flex;
	gap: 4px;
	align-items: center;
}

.message .badges .text-badge {
	padding-top: 4px;
	padding-bottom: 4px;
	padding-right: 8px;
	padding-left: 8px;
	font-size: 11px;
	background-color: #6d6767;
	border-radius: 4px;
	text-transform: uppercase;
}

.message .badges .badge {
	height: 1em;
	width: 1em;
}

.message .profile {
	display: flex;
	flex-direction: v-bind(profileDirection);
	justify-content: space-between;
	gap: 4px;
	width: v-bind(messageWidth);
}

.message > .text {
	flex-wrap: v-bind(messageFlexWrap);
}
</style>
