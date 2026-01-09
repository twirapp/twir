<script setup lang="ts">
import { IconBroadcast } from '@tabler/icons-vue'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

import DialogOrSheet from '@/components/dialog-or-sheet.vue'
import Card from '@/components/overlays/card.vue'
import { Dialog } from '@/components/ui/dialog'
import NowPlaying from '@/pages/overlays/now-playing/now-playing.vue'

const { t } = useI18n()
const isModalOpened = ref(false)
</script>

<template>
	<Card
		:icon="IconBroadcast"
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
