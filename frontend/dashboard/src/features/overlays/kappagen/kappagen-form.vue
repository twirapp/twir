<script setup lang="ts">
import { useForm } from 'vee-validate'
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'

import AnimationsTab from './components/animations-tab.vue'
import EventsTab from './components/events-tab.vue'
import SettingsTab from './components/settings-tab.vue'

import { useKappagenApi } from '@/api/overlays/kappagen'
import { kappagenFormSchema } from './kappagen-form-schema'

import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useToast } from '@/components/ui/toast/use-toast'
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

const isFormValid = computed(() => kappagenForm.meta.valid)
const isFormDirty = computed(() => kappagenForm.meta.dirty)

onMounted(async () => {
	await refetch()
	if (kappagen.value) {
		kappagenForm.setValues({
			enableSpawn: kappagen.value.enableSpawn,
			excludedEmotes: kappagen.value.excludedEmotes,
			enableRave: kappagen.value.enableRave,
			animation: {
				fadeIn: kappagen.value.animation.fadeIn,
				fadeOut: kappagen.value.animation.fadeOut,
				zoomIn: kappagen.value.animation.zoomIn,
				zoomOut: kappagen.value.animation.zoomOut,
			},
			animations: kappagen.value.animations.map((anim) => ({
				style: anim.style,
				prefs: {
					size: anim.prefs.size,
					center: anim.prefs.center,
					speed: anim.prefs.speed,
					faces: anim.prefs.faces,
					message: anim.prefs.message,
					time: anim.prefs.time,
				},
				count: anim.count,
				enabled: anim.enabled,
			})),
			emotes: {
				time: kappagen.value.emotes.time,
				max: kappagen.value.emotes.max,
				queue: kappagen.value.emotes.queue,
				ffzEnabled: kappagen.value.emotes.ffzEnabled,
				bttvEnabled: kappagen.value.emotes.bttvEnabled,
				sevenTvEnabled: kappagen.value.emotes.sevenTvEnabled,
				emojiStyle: kappagen.value.emotes.emojiStyle,
			},
			size: {
				rationNormal: kappagen.value.size.rationNormal,
				rationSmall: kappagen.value.size.rationSmall,
				min: kappagen.value.size.min,
				max: kappagen.value.size.max,
			},
			events: kappagen.value.events.map((event) => ({
				event: event.event,
				disabledAnimations: event.disabledAnimations,
				enabled: event.enabled,
			})),
		})
	}
})

const onSubmit = kappagenForm.handleSubmit(async (values) => {
	try {
		await updateKappagen({
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
					prefs: {
						size: anim.prefs.size,
						center: anim.prefs.center,
						speed: anim.prefs.speed,
						faces: anim.prefs.faces,
						message: anim.prefs.message,
						time: anim.prefs.time,
					},
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

		toast({
			title: t('sharedTexts.success'),
			description: 'Kappagen settings updated successfully',
		})
	} catch (error) {
		console.error('Failed to update Kappagen settings:', error)
		toast({
			title: t('sharedTexts.error'),
			description: 'Failed to update Kappagen settings',
			variant: 'destructive',
		})
	}
})
</script>

<template>
	<PageLayout>
		<template #title> Kappagen overlay </template>

		<template #title-footer>
			<p class="text-sm text-muted-foreground">
				Configure your Kappagen overlay settings, animations, and events
			</p>
		</template>

		<template #action>
			<Button
				:disabled="!isFormValid || !isFormDirty || isUpdating"
				:loading="isUpdating"
				@click="onSubmit"
			>
				{{ isUpdating ? 'Saving...' : 'Save Changes' }}
			</Button>
		</template>

		<template #content>
			<div v-if="isLoading" class="flex items-center justify-center h-64">
				<div class="text-muted-foreground">Loading Kappagen settings...</div>
			</div>

			<form v-else @submit.prevent="onSubmit" class="space-y-6">
				<Tabs default-value="settings" class="w-full">
					<TabsList class="grid w-full grid-cols-3">
						<TabsTrigger value="settings">Settings</TabsTrigger>
						<TabsTrigger value="animations">Animations</TabsTrigger>
						<TabsTrigger value="events">Events</TabsTrigger>
					</TabsList>

					<TabsContent force-mount value="settings" class="mt-6">
						<SettingsTab />
					</TabsContent>

					<TabsContent force-mount value="animations" class="mt-6">
						<AnimationsTab />
					</TabsContent>

					<TabsContent force-mount value="events" class="mt-6">
						<EventsTab />
					</TabsContent>
				</Tabs>
			</form>
		</template>
	</PageLayout>
</template>
