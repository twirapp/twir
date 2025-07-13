<script setup lang="ts">
import { CopyIcon } from 'lucide-vue-next'
import { useThemeVars } from 'naive-ui'
import { useForm } from 'vee-validate'
import { onMounted } from 'vue'
import { useI18n } from 'vue-i18n'

import { kappagenFormSchema } from './kappagen-form-schema'
import KappagenForm from './kappagen-form.vue'

import type { KappagenOverlaySettingsFragment } from '@/gql/graphql'

import { useKappagenApi } from '@/api/overlays/kappagen'
import { useCopyOverlayLink } from '@/components/overlays/copyOverlayLink.ts'
import { Button } from '@/components/ui/button'
import { useToast } from '@/components/ui/toast/use-toast'
import KappagenPreview from '@/features/overlays/kappagen/kappagen-preview.vue'
import { KappagenEmojiStyle } from '@/gql/graphql'
import PageLayout from '@/layout/page-layout.vue'

const { t } = useI18n()
const { toast } = useToast()
const { kappagen, isLoading, isUpdating, updateKappagen, refetch } = useKappagenApi()

// Initialize form with proper variable name
const kappagenForm = useForm({
	validationSchema: kappagenFormSchema,
	initialValues: {
		enableSpawn: true,
		excludedEmotes: [],
		enableRave: false,
		animation: {
			fadeIn: true,
			fadeOut: true,
			zoomIn: false,
			zoomOut: false,
		},
		animations: [],
		emotes: {
			time: 5000,
			max: 100,
			queue: 100,
			ffzEnabled: true,
			bttvEnabled: true,
			sevenTvEnabled: true,
			emojiStyle: KappagenEmojiStyle.Twemoji,
		},
		size: {
			rationNormal: 1,
			rationSmall: 0.5,
			min: 20,
			max: 150,
		},
		events: [],
	},
	keepValuesOnUnmount: true,
})

onMounted(async () => {
	await refetch()
	if (!kappagen.value) return

	const value = kappagen.value as KappagenOverlaySettingsFragment

	kappagenForm.setValues({
		enableSpawn: value.enableSpawn,
		excludedEmotes: value.excludedEmotes,
		enableRave: value.enableRave,
		animation: {
			fadeIn: value.animation.fadeIn,
			fadeOut: value.animation.fadeOut,
			zoomIn: value.animation.zoomIn,
			zoomOut: value.animation.zoomOut,
		},
		animations: value.animations.map((anim) => ({
			style: anim.style,
			prefs: anim.prefs
				? {
					size: anim.prefs.size,
					center: anim.prefs.center,
					speed: anim.prefs.speed,
					faces: anim.prefs.faces,
					message: anim.prefs.message,
					// time: anim.prefs.time,
				}
				: null,
			count: anim.count ?? null,
			enabled: anim.enabled,
		})),
		emotes: {
			time: value.emotes.time,
			max: value.emotes.max,
			queue: value.emotes.queue,
			ffzEnabled: value.emotes.ffzEnabled,
			bttvEnabled: value.emotes.bttvEnabled,
			sevenTvEnabled: value.emotes.sevenTvEnabled,
			emojiStyle: value.emotes.emojiStyle,
		},
		size: {
			rationNormal: value.size.rationNormal,
			rationSmall: value.size.rationSmall,
			min: value.size.min,
			max: value.size.max,
		},
		events: value.events.map((event) => ({
			event: event.event,
			disabledAnimations: event.disabledAnimations,
			enabled: event.enabled,
		})),
	})
})

const onSubmit = kappagenForm.handleSubmit(async (values) => {
	try {
		const { error } = await updateKappagen({
			input: {
				enableSpawn: values.enableSpawn,
				excludedEmotes: values.excludedEmotes,
				enableRave: values.enableRave,
				animation: {
					fadeIn: values.animation.fadeIn,
					fadeOut: values.animation.fadeOut,
					zoomIn: values.animation.zoomIn,
					zoomOut: values.animation.zoomOut,
				},
				animations: values.animations.map((anim) => ({
					style: anim.style,
					prefs: anim.prefs
						? {
							size: anim.prefs.size,
							center: anim.prefs.center,
							speed: anim.prefs.speed,
							faces: anim.prefs.faces,
							message: anim.prefs.message,
							time: anim.prefs.time,
						}
						: null,
					count: anim.count,
					enabled: anim.enabled,
				})),
				emotes: {
					time: values.emotes.time,
					max: values.emotes.max,
					queue: values.emotes.queue,
					ffzEnabled: values.emotes.ffzEnabled,
					bttvEnabled: values.emotes.bttvEnabled,
					sevenTvEnabled: values.emotes.sevenTvEnabled,
					emojiStyle: values.emotes.emojiStyle,
				},
				size: {
					rationNormal: values.size.rationNormal,
					rationSmall: values.size.rationSmall,
					min: values.size.min,
					max: values.size.max,
				},
				events: values.events.map((event) => ({
					event: event.event,
					disabledAnimations: event.disabledAnimations,
					enabled: event.enabled,
				})),
			},
		})

		if (error) {
			throw new Error(error.message)
		}

		toast({
			title: t('sharedTexts.saved'),
			description: 'Kappagen settings updated successfully',
			duration: 2500,
		})
	} catch (error) {
		console.error('Failed to update Kappagen settings:', error)
		toast({
			description: 'Failed to update Kappagen settings',
			variant: 'destructive',
		})
	}
})

const themeVars = useThemeVars()

const { copyOverlayLink } = useCopyOverlayLink('kappagen')
</script>

<template>
	<PageLayout clean-body>
		<template #title>
			Kappagen overlay
		</template>

		<template #title-footer>
			<p class="text-sm text-muted-foreground">
				Flying emotes on your screen!
			</p>
		</template>

		<template #action>
			<div class="flex flex-col gap-2">
				<Button :loading="isUpdating" @click="onSubmit">
					{{ isUpdating ? 'Saving...' : 'Save Changes' }}
				</Button>
				<Button variant="outline" class="flex items-center gap-2" @click="copyOverlayLink()">
					<CopyIcon class="size-4" />
					Copy overlay link
				</Button>
			</div>
		</template>

		<template #content>
			<div v-if="isLoading" class="flex items-center justify-center h-64">
				<div class="text-muted-foreground">
					Loading Kappagen settings...
				</div>
			</div>
			<div v-else class="p-8">
				<form
					class="relative w-full rounded-lg h-[80dvh] border-2 border-border bg-background/60 shadow-lg"
					:style="{ backgroundColor: themeVars.cardColor }"
					@submit.prevent="onSubmit"
				>
					<KappagenForm class="absolute top-2 left-4 h-full" />
					<KappagenPreview class="w-full" />
				</form>
			</div>
		</template>
	</PageLayout>
</template>
