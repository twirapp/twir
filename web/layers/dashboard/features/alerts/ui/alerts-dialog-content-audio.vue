<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useProfile } from '~~/layers/dashboard/api/auth.js'
import { useFilesApi } from '~~/layers/dashboard/api/files.js'
import DialogOrSheet from '~~/layers/dashboard/components/dialog-or-sheet.vue'
import FilesPicker from '~~/layers/dashboard/components/files/files.vue'

import { Button } from '@/components/ui/button'
import { Dialog, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import { Slider } from '@/components/ui/slider'

const audioId = defineModel<string | null>('audioId')
const volume = defineModel<number>('volume', {
	default: 30,
})

const { t } = useI18n()
const { data: profile } = useProfile()
const filesApi = useFilesApi()
const { data: files } = filesApi.useQuery()

const selectedAudio = computed(() => {
	return files.value?.files.find((file) => file.id === audioId.value)
})

const showAudioDialog = ref(false)
const isAudioPlaying = ref(false)
const loadedAudios = new Map<string, HTMLAudioElement>()

watch(audioId, (audioId, prevValue) => {
	if (audioId === undefined) {
		loadedAudios.get(prevValue!)?.pause()
		isAudioPlaying.value = false
	}
})

async function playAudio(audio: HTMLAudioElement) {
	setVolume(audioId.value!, volume.value)
	audio.currentTime = 0
	audio.play()
}

async function getAudio() {
	if (!selectedAudio.value?.id || !profile.value) return
	const audioId = selectedAudio.value.id

	if (loadedAudios.has(audioId)) {
		const audio = loadedAudios.get(audioId)!
		if (isAudioPlaying.value) {
			audio.pause()
			isAudioPlaying.value = false
		} else {
			playAudio(audio)
		}
		return
	}

	const audio = new Audio(filesApi.computeFileUrl(profile.value.selectedDashboardId, audioId))
	audio.addEventListener('error', (error) => {
		console.error(error)
	})

	audio.addEventListener('loadstart', () => {
		loadedAudios.set(audioId, audio)
		playAudio(audio)
	})

	audio.addEventListener('play', () => {
		isAudioPlaying.value = true
	})

	audio.addEventListener('ended', () => {
		isAudioPlaying.value = false
	})
}

function setVolume(audioId: string, v: number) {
	const audio = loadedAudios.get(audioId)
	if (!audio) return
	audio.volume = v / 100
	volume.value = v
}
</script>

<template>
	<div class="flex flex-col gap-2">
		<span>{{ t('alerts.select.audio') }}</span>
		<div class="flex w-full gap-2">
			<Dialog
				v-model:open="showAudioDialog"
				@update:open="showAudioDialog = false"
			>
				<DialogTrigger as-child>
					<Button
						class="w-[80%]"
						@click="showAudioDialog = true"
					>
						{{ selectedAudio?.name ?? t('sharedButtons.select') }}
					</Button>
				</DialogTrigger>

				<DialogOrSheet class="h-[80dvh] min-w-[50%] gap-0 p-0 md:h-auto">
					<DialogHeader class="border-b p-6">
						<DialogTitle>
							{{ t('alerts.select.audio') }}
						</DialogTitle>
					</DialogHeader>

					<FilesPicker
						class="h-auto md:max-h-[50dvh]"
						mode="picker"
						tab="audios"
						@select="
							(id) => {
								audioId = id
								showAudioDialog = false
							}
						"
						@delete="
							(id) => {
								if (id === audioId) {
									audioId = undefined
								}
							}
						"
					/>
				</DialogOrSheet>
			</Dialog>

		<Button
			class="min-w-10"
			size="icon"
			variant="secondary"
			:disabled="!audioId"
			@click.stop.prevent="getAudio"
		>
				<Icon
					name="lucide:play"
					v-if="!isAudioPlaying"
					class="size-4"
				/>
				<Icon
					name="lucide:pause"
					v-else
					class="size-4"
				/>
			</Button>

		<Button
			class="min-w-10"
			size="icon"
			variant="destructive"
			:disabled="!audioId"
			@click.stop.prevent="audioId = undefined"
		>
				<Icon
					name="lucide:trash"
					class="size-4"
				/>
			</Button>
		</div>
	</div>

	<div class="flex flex-col gap-2">
		<div class="flex justify-between">
			<span>{{ t('alerts.audioVolume') }}</span>
			<span>{{ volume }}%</span>
		</div>
		<Slider
			:model-value="[volume]"
			@update:model-value="
				(val) => {
					if (!val) return

					volume = val[0]
					if (audioId) {
						setVolume(audioId, volume)
					}
				}
			"
			:max="100"
			:min="0"
			:step="1"
		/>
	</div>
</template>
