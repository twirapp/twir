<script setup lang="ts">
import { PauseIcon, PlayIcon, TrashIcon } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { useProfile } from '@/api/auth.js'
import { useFilesApi } from '@/api/files.js'
import DialogOrSheet from '@/components/dialog-or-sheet.vue'
import FilesPicker from '@/components/files/files.vue'
import { Button } from '@/components/ui/button'
import {
	Dialog,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '@/components/ui/dialog'
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
	return files.value?.files
		.find((file) => file.id === audioId.value)
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
		<div class="flex gap-2 w-full">
			<Dialog v-model:open="showAudioDialog" @update:open="showAudioDialog = false">
				<DialogTrigger as-child>
					<Button class="w-[80%]" @click="showAudioDialog = true">
						{{ selectedAudio?.name ?? t('sharedButtons.select') }}
					</Button>
				</DialogTrigger>

				<DialogOrSheet class="p-0 gap-0 h-[80dvh] md:h-auto">
					<DialogHeader class="p-6 border-b">
						<DialogTitle>
							{{ t('alerts.select.audio') }}
						</DialogTitle>
					</DialogHeader>

					<FilesPicker
						class="h-auto md:max-h-[50dvh]"
						mode="picker"
						tab="audios"
						@select="(id) => {
							audioId = id
							showAudioDialog = false
						}"
						@delete="(id) => {
							if (id === audioId) {
								audioId = undefined
							}
						}"
					/>
				</DialogOrSheet>
			</Dialog>

			<Button
				class="min-w-10"
				size="icon"
				variant="secondary"
				:disabled="!audioId"
				@click="getAudio"
			>
				<PlayIcon v-if="!isAudioPlaying" class="size-4" />
				<PauseIcon v-else class="size-4" />
			</Button>

			<Button
				class="min-w-10"
				size="icon"
				variant="destructive"
				:disabled="!audioId"
				@click="audioId = undefined"
			>
				<TrashIcon class="size-4" />
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
			@update:model-value="(val) => {
				if (!val) return;

				volume = val[0]
				if (audioId) {
					setVolume(audioId, volume)
				}
			}"
			:max="100"
			:min="0"
			:step="1"
		/>
	</div>
</template>
