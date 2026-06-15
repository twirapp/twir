<script setup lang="ts">
import { ref } from 'vue'

import DialogOrSheet from '~~/layers/dashboard/app/components/dialog-or-sheet.vue'
import Card from '~~/layers/dashboard/app/components/overlays/card.vue'
import { Dialog } from '@/components/ui/dialog'
import NowPlaying from '~~/layers/dashboard/app/pages/overlays/now-playing.vue'

const { t } = useI18n()
const isModalOpened = ref(false)
</script>

<template>
	<Card
		icon="lucide:broadcast"
		title="Now Playing"
		:description="t('overlays.nowPlaying.description')"
		:show-copy="false"
		@open-settings="isModalOpened = true"
	>
	</Card>

	<Dialog v-model:open="isModalOpened">
		<DialogOrSheet>
			<NowPlaying />
		</DialogOrSheet>
	</Dialog>
</template>
