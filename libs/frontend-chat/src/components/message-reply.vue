<script setup lang="ts">
import type { MessageReply } from '../types.js'

defineProps<{
	reply: MessageReply
	variant: 'compact' | 'card' | 'inline'
}>()
</script>

<template>
	<div v-if="variant === 'compact'" class="reply reply-compact">
		<svg
			xmlns="http://www.w3.org/2000/svg"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="2"
			stroke-linecap="round"
			stroke-linejoin="round"
			class="reply-icon"
		>
			<path d="M9 14 4 9l5-5" />
			<path d="M20 20v-7a4 4 0 0 0-4-4H4" />
		</svg>
		<span class="reply-line">
			Replying to <span class="reply-name">@{{ reply.parentUserName }}</span>:
			{{ reply.parentMessageBody }}
		</span>
	</div>

	<div v-else-if="variant === 'card'" class="reply reply-card">
		<div class="reply-card-header">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="2"
				stroke-linecap="round"
				stroke-linejoin="round"
				class="reply-icon"
			>
				<path d="M9 14 4 9l5-5" />
				<path d="M20 20v-7a4 4 0 0 0-4-4H4" />
			</svg>
			<span class="reply-card-title">
				Replying to <span class="reply-name">@{{ reply.parentUserName }}</span>
			</span>
		</div>
		<div class="reply-card-text">{{ reply.parentMessageBody }}</div>
	</div>

	<span v-else class="reply reply-inline">
		<svg
			xmlns="http://www.w3.org/2000/svg"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="2"
			stroke-linecap="round"
			stroke-linejoin="round"
			class="reply-icon"
		>
			<path d="M9 14 4 9l5-5" />
			<path d="M20 20v-7a4 4 0 0 0-4-4H4" />
		</svg>
		<span class="reply-inline-text">
			<span class="reply-name">@{{ reply.parentUserName }}</span>: {{ reply.parentMessageBody }}
		</span>
	</span>
</template>

<style scoped>
.reply-icon {
	display: inline-block;
	height: 1em;
	width: 1em;
	flex: none;
	vertical-align: -0.15em;
}

.reply-name {
	color: rgba(255, 255, 255, 0.85);
	font-weight: 600;
}

/* Compact variant (clean preset): single muted truncated line */
.reply-compact {
	display: block;
	width: 100%;
	max-width: 100%;
	margin-bottom: 0.15em;
	font-size: 0.7em;
	line-height: 1.4;
	color: rgba(255, 255, 255, 0.55);
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
	box-sizing: border-box;
}

.reply-compact .reply-icon {
	margin-right: 0.3em;
}

/* Card variant (boxed preset): quote-style block inside the box */
.reply-card {
	display: flex;
	flex-direction: column;
	gap: 0.1em;
	width: 100%;
	max-width: 100%;
	margin-bottom: 0.2em;
	padding: 0.35em 0.55em;
	font-size: 0.8em;
	line-height: 1.35;
	background-color: rgba(255, 255, 255, 0.06);
	border-left: 3px solid rgba(255, 255, 255, 0.25);
	border-radius: 6px;
	white-space: normal;
	box-sizing: border-box;
}

.reply-card-header {
	display: flex;
	align-items: center;
	gap: 0.3em;
	font-size: 0.9em;
	color: rgba(255, 255, 255, 0.7);
	white-space: nowrap;
	overflow: hidden;
}

.reply-card-title {
	overflow: hidden;
	text-overflow: ellipsis;
}

.reply-card-text {
	color: rgba(255, 255, 255, 0.6);
	display: -webkit-box;
	-webkit-line-clamp: 2;
	line-clamp: 2;
	-webkit-box-orient: vertical;
	overflow: hidden;
	word-break: break-word;
}

/* Inline variant (horizontal layouts): prefix inside the message line */
.reply-inline {
	display: inline-flex;
	align-items: center;
	gap: 0.25em;
	max-width: 15em;
	margin-right: 0.35em;
	font-size: 0.7em;
	line-height: 1.4;
	color: rgba(255, 255, 255, 0.55);
	white-space: nowrap;
	overflow: hidden;
	vertical-align: middle;
}

.reply-inline-text {
	overflow: hidden;
	text-overflow: ellipsis;
}
</style>
