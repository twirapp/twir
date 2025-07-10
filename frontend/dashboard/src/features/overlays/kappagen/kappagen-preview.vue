<script setup lang="ts">
import { computed, ref, toRaw, watch } from 'vue'
import { useFormValues } from 'vee-validate'
import { useI18n } from 'vue-i18n'

import { type KappagenFormSchema } from '@/features/overlays/kappagen/kappagen-form-schema.ts'
import { Button } from '@/components/ui/button'

import KappagenOverlay from '@twirapp/kappagen'
import type { KappagenMethods } from '@twirapp/kappagen/types'

import { useProfile } from '@/api'

const twirEmote: Emote = {
	url: 'https://cdn.7tv.app/emote/6548b7074789656a7be787e1/4x.webp',
	zwe: [
		{
			url: 'https://cdn.7tv.app/emote/6128ed55a50c52b1429e09dc/4x.webp',
		},
	],
}

const { t } = useI18n()
const { data: profile } = useProfile()
const kappagen = ref<KappagenMethods>()
const kappagenIframeRef = ref<HTMLIFrameElement | null>(null)
const kappagenIframeUrl = computed(() => {
	if (!profile.value) return null

	return `${window.location.origin}/overlays/${profile.value.apiKey}/kappagen`
})
const { value: formSettings } = useFormValues<KappagenFormSchema>()

watch(formSettings, (v) => {
	if (!v) return

	return sendIframeMessage('settings', {
		...toRaw(formSettings.value),
		channelName: profile.value?.login,
		channelId: profile.value?.id,
	})
})

function sendIframeMessage(key: string, data?: any) {
	if (!kappagenIframeRef.value) return
	const win = kappagenIframeRef.value

	win.contentWindow?.postMessage(
		JSON.stringify({
			key,
			data: toRaw(data),
		})
	)
}

function playAnimation(animation: KappagenAnimations) {
	if (!kappagen.value) return Promise.resolve()

	// const randomAnimation = formSettings.animations.filter((a) => a.enabled)[
	// 	Math.floor(Math.random() * formSettings.animations.length)
	// ]

	const randomAnimation = formSettings.animations.find((a) => a.style === 'TheCube')

	console.log(randomAnimation)

	return kappagen.value.playAnimation([twirEmote], randomAnimation)
}

function showEmotes(emotes: Emote[]) {
	if (!kappagen.value) return
	kappagen.value.showEmotes(emotes)
}

function clear() {
	if (!kappagen.value) return
	kappagen.value.clear()
}
</script>

<template>
	<div class="flex flex-col gap-6">
		<div class="flex flex-wrap justify-between">
			<span class="text-2xl font-bold">Preview</span>
			<div class="flex flex-wrap gap-2 justify-end">
				<Button variant="secondary" @click="sendIframeMessage('kappa', 'EZ')" type="button">
					{{ t('overlays.kappagen.testKappagen') }}
				</Button>
				<Button variant="default" type="button" @click="sendIframeMessage('spawn', ['EZ'])">
					{{ t('overlays.kappagen.testSpawn') }}
				</Button>
				<Button
					variant="destructive"
					@click="sendIframeMessage('clear')"
					type="button"
					class="z-[999999]"
				>
					{{ t('overlays.kappagen.clear') }}
				</Button>
			</div>
		</div>
		<div class="border rounded-md bg-sidebar/80 max-h-[50%] transform aspect-[16/9]">
			<iframe
				v-if="kappagenIframeUrl"
				ref="kappagenIframeRef"
				:src="kappagenIframeUrl"
				class="w-full h-full"
			/>

			<!--			<KappagenOverlay ref="kappagen" :config="formSettings" is-rave />-->
		</div>
	</div>
</template>
