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

const { value: formSettings } = useFormValues<KappagenFormSchema>()
</script>

<template>
	<div class="h-full">
		<div class="absolute top-4 right-4 flex gap-2 z-50">
			<button
				@click="kappagen.showEmotes([twirEmote])"
				class="px-4 py-2 rounded-md border-stone-700/50 border bg-indigo-600 shadow-lg"
				type="button"
			>
				Spawn emote
			</button>
			<button
				@click="kappagen.clear"
				class="px-4 py-2 rounded-md border-stone-700/50 border bg-stone-700/40 shadow-lg"
				type="button"
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
