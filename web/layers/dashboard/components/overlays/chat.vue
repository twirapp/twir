<script setup lang="ts">
import { useRouter } from 'vue-router'
import Card from '~~/layers/dashboard/components/overlays/card.vue'

const { t } = useI18n()

const router = useRouter()
const localePath = useLocalePath()
</script>

<template>
	<Card
		:icon="'tabler:icon-message'"
		:icon-stroke="1"
		title="Chat"
		:description="t('overlays.chat.description')"
		overlay-path="chat"
		:show-copy="false"
		@open-settings="router.push(localePath('/dashboard/overlays/chat'))"
	>
	</Card>
</template>
