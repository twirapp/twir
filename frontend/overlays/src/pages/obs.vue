<script lang="ts" setup>
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'

import { useObs } from '@/composables/obs/use-obs.js'
import {
	type ObsCommand,
	type ObsSettings,
	useObsOverlayGraphQL,
} from '@/composables/obs/use-obs-graphql.js'
import { ObsWebsocketCommandAction } from '@/gql/graphql.js'

const route = useRoute()
const obs = useObs()

const apiKey = route.params.apiKey as string
const settings = ref<ObsSettings | null>(null)
const isObsConnected = ref(false)

const { connect, destroy, updateSources, startHeartbeat, stopHeartbeat } = useObsOverlayGraphQL({
	onSettings: async (newSettings) => {
		settings.value = newSettings
		await handleSettingsChange(newSettings)
	},
	onCommand: (command) => {
		handleCommand(command)
	},
})

function handleCommand(command: ObsCommand) {
	switch (command.action) {
		case ObsWebsocketCommandAction.SetScene:
			obs.setScene(command.target)
			break
		case ObsWebsocketCommandAction.ToggleSource:
			obs.toggleSource(command.target)
			break
		case ObsWebsocketCommandAction.ToggleAudio:
			obs.toggleAudioSource(command.target)
			break
		case ObsWebsocketCommandAction.SetVolume:
			if (command.volumeValue !== null && command.volumeValue !== undefined) {
				obs.setVolume(command.target, command.volumeValue)
			}
			break
		case ObsWebsocketCommandAction.IncreaseVolume:
			if (command.volumeStep !== null && command.volumeStep !== undefined) {
				obs.changeVolume(command.target, command.volumeStep, 'increase')
			}
			break
		case ObsWebsocketCommandAction.DecreaseVolume:
			if (command.volumeStep !== null && command.volumeStep !== undefined) {
				obs.changeVolume(command.target, command.volumeStep, 'decrease')
			}
			break
		case ObsWebsocketCommandAction.EnableAudio:
			obs.toggleAudioSource(command.target, true)
			break
		case ObsWebsocketCommandAction.DisableAudio:
			obs.toggleAudioSource(command.target, false)
			break
		case ObsWebsocketCommandAction.StartStream:
			obs.startStream()
			break
		case ObsWebsocketCommandAction.StopStream:
			obs.stopStream()
			break
	}
}

async function handleSettingsChange(newSettings: ObsSettings) {
	if (!newSettings.serverAddress || !newSettings.serverPort || !newSettings.serverPassword) {
		// Stop heartbeat and disconnect if settings are invalid
		stopHeartbeat()
		isObsConnected.value = false
		await obs.disconnect()
		return
	}

	try {
		await obs.connect(newSettings.serverAddress, newSettings.serverPort, newSettings.serverPassword)
		console.log('Twir OBS WebSocket connected')
		isObsConnected.value = true

		// Start heartbeat only after successful OBS connection
		startHeartbeat()

		// Send initial data after connect
		await sendObsData()

		// Setup listeners for OBS changes
		setupObsListeners()
	} catch (error) {
		console.error('Failed to connect to OBS:', error)
		isObsConnected.value = false
		stopHeartbeat()
	}
}

async function sendObsData() {
	const [sources, audioSources] = await Promise.all([obs.getSources(), obs.getAudioSources()])

	if (!sources) return

	// Extract scene names from sources object keys
	const scenes = Object.keys(sources)
	// Flatten all sources from all scenes
	const allSources = Object.values(sources)
		.flat()
		.map((s) => s.name)

	await updateSources({
		scenes,
		sources: allSources,
		audioSources: audioSources ?? [],
	})
}

function setupObsListeners() {
	const updateHandler = async () => {
		await sendObsData()
	}

	// Handle OBS disconnect
	obs.instance.value.on('ConnectionClosed', () => {
		console.log('OBS WebSocket disconnected')
		isObsConnected.value = false
		stopHeartbeat()
	})

	obs.instance.value
		.on('SceneListChanged', updateHandler)
		.on('InputCreated', updateHandler)
		.on('InputRemoved', updateHandler)
		.on('InputNameChanged', updateHandler)
		.on('SceneItemCreated', updateHandler)
		.on('SceneItemRemoved', updateHandler)
}

onMounted(() => {
	connect(apiKey)
})

onBeforeUnmount(() => {
	destroy()
})
</script>

<template>
	<!-- OBS overlay has no visible UI, it only controls OBS -->
</template>
