<script setup lang="ts">
import { watch } from 'vue'

import { useDudesOverlayManager } from '~~/layers/dashboard/api/overlays/dudes'
import DudesSettingsForm from '~~/layers/dashboard/pages/dashboard/overlays/dudes/dudes-settings-form.vue'
import { useDudesForm } from '~~/layers/dashboard/pages/dashboard/overlays/dudes/use-dudes-form'

const { data: entities } = useDudesOverlayManager().useGetAll()
const { setData } = useDudesForm()

watch(
	() => entities.value?.dudesGetAll,
	(overlays) => {
		const firstOverlay = overlays?.[0]
		if (firstOverlay) setData(firstOverlay)
	},
	{ immediate: true },
)
</script>

<template>
	<DudesSettingsForm embedded />
</template>
