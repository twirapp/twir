<script setup lang="ts">
import { PauseIcon, PlayIcon, TrashIcon } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'


// ...existing code...
import { useFilesApi } from '#layers/dashboard/api/files'
import DialogOrSheet from '#layers/dashboard/components/dialog-or-sheet.vue'
import FilesPicker from '#layers/dashboard/components/files/files.vue'




const audioId = defineModel<string | null>('audioId')
const volume = defineModel<number>('volume', {
	default: 30,
})

const { t } = useI18n()
const { user: profile } = storeToRefs(useDashboardAuth())
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
			<UiDialog v-model:open="showAudioDialog" @update:open="showAudioDialog = false">
				<UiDialogTrigger as-child>
					<UiButton class="w-[80%]" @click="showAudioDialog = true">
						{{ selectedAudio?.name ?? t('sharedButtons.select') }}
					</UiButton>
				</UiDialogTrigger>

				<DialogOrSheet class="p-0 gap-0 h-[80dvh] md:h-auto">
					<UiDialogHeader class="p-6 border-b">
						<UiDialogTitle>
							{{ t('alerts.select.audio') }}
						</UiDialogTitle>
					</UiDialogHeader>

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
			</UiDialog>

			<UiButton
				class="min-w-10"
				size="icon"
				variant="secondary"
				:disabled="!audioId"
				@click="getAudio"
			>
				<PlayIcon v-if="!isAudioPlaying" class="size-4" />
				<PauseIcon v-else class="size-4" />
			</UiButton>

			<UiButton
				class="min-w-10"
				size="icon"
				variant="destructive"
				:disabled="!audioId"
				@click="audioId = undefined"
			>
				<TrashIcon class="size-4" />
			</UiButton>
		</div>
	</div>

	<div class="flex flex-col gap-2">
		<div class="flex justify-between">
			<span>{{ t('alerts.audioVolume') }}</span>
			<span>{{ volume }}%</span>
		</div>
		<UiSlider
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
