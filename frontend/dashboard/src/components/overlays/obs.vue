<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

import Settings from './obs/settings.vue'

import { useObsOverlayManager } from '@/api/index.js'
import BroadcastIcon from '@/assets/overlays/broadcast.svg?use'
import DialogOrSheet from '@/components/dialog-or-sheet.vue'
import Card from '@/components/overlays/card.vue'
import { Dialog } from '@/components/ui/dialog'

const isModalOpened = ref(false)
const obsManager = useObsOverlayManager()
const { data: obsSettings, isError } = obsManager.getSettings(false)

const { t } = useI18n()
</script>

<template>
	<Card
		:icon="BroadcastIcon"
		:icon-stroke="2"
		title="OBS"
		:description="t('overlays.obs.description')"
		overlay-path="obs"
		:copy-disabled="!obsSettings || isError"
		@open-settings="isModalOpened = true"
	>
	</Card>

	<Dialog v-model:open="isModalOpened">
		<DialogOrSheet class="w-full md:max-w-xl">
			<Settings />
		</DialogOrSheet>
	</Dialog>
</template>
