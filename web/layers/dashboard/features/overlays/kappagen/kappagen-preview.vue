<script setup lang="ts">
import KappagenOverlay from '@twirapp/kappagen'
import '@twirapp/kappagen/styles'

import { useFormValues } from 'vee-validate'
import { computed } from 'vue'

import type { KappagenFormSchema } from '@/features/overlays/kappagen/kappagen-form-schema.ts'
import type { KappagenConfig, KappagenMethods } from '@twirapp/kappagen/types'

import {
	twirEmote,
	useKappagenInstance,
} from '@/features/overlays/kappagen/composables/use-kappagen-instance.ts'

const kappagen = useKappagenInstance()

const { value: formSettings } = useFormValues<KappagenFormSchema>()

const transofrmedSettings = computed<KappagenConfig>(() => {
	return {
		time: formSettings.emotes?.time ?? 5000,
		max: formSettings.emotes?.max ?? 100,
		queue: formSettings.emotes?.queue ?? 100,
		animation: {
			fade: {
				in: formSettings.animation?.fadeIn ? 1 : 0,
				out: formSettings.animation?.fadeOut ? 1 : 0,
			},
			zoom: {
				in: formSettings.animation?.zoomIn ? 1 : 0,
				out: formSettings.animation?.zoomOut ? 1 : 0,
			},
		},
		in: {
			fade: formSettings.animation?.fadeIn,
			zoom: formSettings.animation?.zoomIn,
		},
		out: {
			fade: formSettings.animation?.fadeOut,
			zoom: formSettings.animation?.zoomOut,
		},
	}
})
</script>

<template>
	<div class="h-full">
		<div class="absolute top-4 right-4 flex gap-2 z-50">
			<button
				class="px-4 py-2 rounded-md border-stone-700/50 border bg-indigo-600 shadow-lg"
				type="button"
				@click="kappagen.showEmotes([twirEmote])"
			>
				Spawn emote
			</button>
			<button
				class="px-4 py-2 rounded-md border-stone-700/50 border bg-stone-700/40 shadow-lg"
				type="button"
				@click="kappagen.clear"
			>
				Clear overlay
			</button>
		</div>
		<KappagenOverlay
			:ref="(el) => kappagen.setKappagenInstance(el as unknown as KappagenMethods)"
			:config="transofrmedSettings"
			:is-rave="formSettings.enableRave"
		/>
	</div>
</template>
