<script setup lang="ts">
import KappagenOverlay from '@twirapp/kappagen'
import { useFormValues } from 'vee-validate'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

import type { KappagenFormSchema } from '@/features/overlays/kappagen/kappagen-form-schema.ts'
import type { Emote, KappagenMethods } from '@twirapp/kappagen/types'

import { Button } from '@/components/ui/button'
import {
	twirEmote,
	useKappagenInstance,
} from '@/features/overlays/kappagen/composables/use-kappagen-instance.ts'

const { t } = useI18n()

const kappagen = useKappagenInstance()

// const kappagen = ref<KappagenMethods>()
const { value: formSettings } = useFormValues<KappagenFormSchema>()

// function playAnimation(animation: KappagenAnimations) {
// 	if (!kappagen.value) return Promise.resolve()
//
// 	// const randomAnimation = formSettings.animations.filter((a) => a.enabled)[
// 	// 	Math.floor(Math.random() * formSettings.animations.length)
// 	// ]
//
// 	console.log(formSettings.animations)
//
// 	const randomAnimation = formSettings.animations.find((a) => a.style === 'THE_CUBE')
//
// 	console.log(randomAnimation)
//
// 	return kappagen.value.playAnimation([twirEmote], {
// 		...randomAnimation,
// 		style: 'TheCube',
// 	})
// }
//
// function showEmotes(emotes: Emote[]) {
// 	console.log(kappagen.value.showEmotes, emotes)
// 	if (!kappagen.value) return
// 	kappagen.value.showEmotes(emotes)
// }
//
// function clear() {
// 	if (!kappagen.value) return
// 	kappagen.value.clear()
// }
</script>

<template>
	<div class="h-full">
		<div class="absolute top-4 right-4 flex gap-2 z-50">
			<button
				@click="kappagen.showEmotes([twirEmote])"
				class="px-4 py-2 rounded-md border-stone-700/50 border bg-indigo-600 shadow-lg"
			>
				Spawn emote
			</button>
			<button
				@click="kappagen.clear"
				class="px-4 py-2 rounded-md border-stone-700/50 border bg-stone-700/40 shadow-lg"
			>
				Clear overlay
			</button>
		</div>
		<KappagenOverlay
			:ref="kappagen.setKappagenInstance"
			:config="formSettings"
			:is-rave="formSettings.enableRave"
		/>
	</div>
</template>
