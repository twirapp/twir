<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';

import { useChatSocket } from '../sockets/chat.js';
import { useTmiChat } from '../sockets/chat_tmi.js';
import { normalizeDisplayName } from '../sockets/chat_tmi_helpers.js';

const route = useRoute();
const apiKey = route.params.apiKey as string;

const chat = useChatSocket(apiKey);

const channelName = computed(() => chat.settings.channelName);
const channelId = computed(() => chat.settings.channelId);
const messageTimeout = computed(() => chat.settings.messageTimeout);
const { messages } = useTmiChat(
	channelName,
	channelId,
	messageTimeout,
);
</script>

<template>
	<div class="chat">
		<TransitionGroup name="list" tag="div" class="messages">
			<div v-for="(msg, index) of messages" :key="index" class="message">
				<div v-if="msg.badges?.length || msg.sender" class="profile">
					<div v-if="msg.badges" class="badges">
						<template
							v-for="(badgeValue, badgeName) of msg.badges"
							:key="badgeName+badgeValue"
						>
							<img
								v-if="chat.settings.channelBadges.get(`${badgeName}-${badgeValue}`)"
								:src="chat.settings.channelBadges.get(`${badgeName}-${badgeValue}`)!.image_url_4x"
								class="badge"
							/>

							<img
								v-else-if="chat.settings.globalBadges.get(badgeName)?.versions.length"
								:src="chat.settings.globalBadges.get(badgeName)!.versions.at(-1)!.image_url_4x"
								class="badge"
							/>
						</template>
					</div>
					<div v-if="msg.senderDisplayName" :style="{ color: msg.senderColor }">
						{{ normalizeDisplayName(msg.sender!, msg.senderDisplayName!) }}{{ msg.isItalic ? '' : ':' }}
					</div>
				</div>
				<span class="text" :style="{ fontStyle: msg.isItalic ? 'italic' : 'normal' }">
					<template v-for="(chunk, _) of msg.chunks" :key="_">
						<img
							v-if="chunk.type === 'emote'"
							:src="`https://static-cdn.jtvnw.net/emoticons/v2/${chunk.value}/default/dark/3.0`"
							class="emote"
						/>

						<img
							v-else-if="chunk.type === 'ffz_emote'"
							:src="chunk.value"
							class="emote"
						/>

						<img
							v-else-if="chunk.type === 'bttv_emote'"
							:src="chunk.value"
							class="emote"
						/>

						<img
							v-if="chunk.type === '7tv_emote'"
							:src="chunk.value"
							class="emote"
						/>

						<template v-else-if="chunk.type === 'text'">
							{{ chunk.value }}
						</template>
						{{ ' ' }}
					</template>
				</span>
			</div>
		</TransitionGroup>
	</div>
</template>

<style scoped>
@import url(https://fonts.googleapis.com/css?family=Roboto:700);

.chat {
  height: 100%;
  width: 100%;
  background-color: #000;
  color: #fff;
  font-size: 20px;
	font-family: 'Roboto';
}

.chat .messages {
	display: flex;
	flex-direction: column;
	gap: 8px;
	overflow: hidden;
}

.chat .messages .message {
  display: flex;
  gap: 8px;
  justify-content: flex-start;
  align-items: flex-start;
	width: 100%;
}

.chat .messages .message .badges {
  display: flex;
  gap: 4px;
}

.chat .messages .message .badges .badge {
	height: 20px;
	width: 20px;
}

.chat .messages .message .profile {
	display: flex;
  flex-wrap: nowrap;
  gap: 4px;
  align-items: center;
}

.chat .messages .message .text {

}

.chat .messages .message .text .emote {
	height: 20px;
	width: 20px;
}

.list-move, /* apply transition to moving elements */
.list-enter-active,
.list-leave-active {
  transition: all 0.5s ease;
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateX(30px);
}

/* ensure leaving items are taken out of layout flow so that moving
   animations can be calculated correctly. */
.list-leave-active {
  position: absolute;
}
</style>
