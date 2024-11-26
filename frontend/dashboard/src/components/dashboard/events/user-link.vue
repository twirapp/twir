<script setup lang="ts">
import { computed } from 'vue'

import { resolveUserName } from '@/helpers/resolveUserName.js'

const props = defineProps<{
	name: string
	displayName: string
}>()

const userName = computed(() => resolveUserName(props.name, props.displayName))
const channelLink = computed(() => `https://twitch.tv/${userName.value}`)
</script>

<template>
	<a
		:href="channelLink"
		target="_blank"
		class="hover:underline"
	>
		{{ userName }}
	</a>
</template>
