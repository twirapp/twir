<script lang="ts" setup>
import { useWebSocket } from '@vueuse/core'
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'

import { useObs } from '@/composables/obs/use-obs.js'
import { generateSocketUrlWithParams } from '@/helpers.js'

const obs = useObs()
const route = useRoute()

const apiKey = route.params.apiKey as string
const obsUrl = generateSocketUrlWithParams('/overlays/obs', {
	apiKey,
})

const internalSocket = useWebSocket(obsUrl, {
	immediate: true,
	autoReconnect: {
		delay: 500,
	},
	onConnected(ws) {
		ws.send(JSON.stringify({ eventName: 'requestSettings' }))
	},
})

const settings = ref<Record<string, any> | null>(null)

watch(internalSocket.data, (message) => {
	const { eventName, data } = JSON.parse(message)

	switch (eventName) {
		case 'settings':
			settings.value = data
			break
		case 'setScene':
			obs.setScene(data.sceneName)
			break
		case 'toggleSource':
			obs.toggleSource(data.sourceName)
			break
		case 'toggleAudioSource':
			obs.toggleAudioSource(data.audioSourceName)
			break
		case 'setVolume':
			obs.setVolume(data.audioSourceName, data.volume)
			break
		case 'increaseVolume':
			obs.changeVolume(data.audioSourceName, data.step, 'increase')
			break
		case 'decreaseVolume':
			obs.changeVolume(data.audioSourceName, data.step, 'decrease')
			break
		case 'enableAudio':
			obs.toggleAudioSource(data.audioSourceName, true)
			break
		case 'disableAudio':
			obs.toggleAudioSource(data.audioSourceName, false)
			break
		case 'startStream':
			obs.startStream()
			break
		case 'stopStream':
			obs.stopStream()
			break
	}
})

watch(settings, async (settings) => {
	if (!settings) {
		await obs.disconnect()
		return
	}

	await obs.connect(settings.serverAddress, settings.serverPort, settings.serverPassword)
	console.log('Twir obs socket opened')

	internalSocket.send(JSON.stringify({ eventName: 'obsConnected' }))

	obs.getSources().then((sources) => {
		if (!sources) return
		internalSocket.send(JSON.stringify({
			eventName: 'setSources',
			data: sources,
		}))
	})

	obs.getAudioSources().then((sources) => {
		if (!sources) return
		internalSocket.send(JSON.stringify({
			eventName: 'setAudioSources',
			data: sources,
		}))
	})

	const scenesHandler = async () => {
		const sources = await obs.getSources()
		internalSocket.send(JSON.stringify({
			eventName: 'setSources',
			data: sources,
		}))
	}

	const audioHandler = async () => {
		const sources = await obs.getAudioSources()
		internalSocket.send(JSON.stringify({
			eventName: 'setAudioSources',
			data: sources,
		}))
	}

	obs.instance.value
		.on('SceneListChanged', scenesHandler)

		.on('InputCreated', audioHandler)
		.on('InputRemoved', audioHandler)
		.on('InputNameChanged', audioHandler)

		.on('SceneItemCreated', scenesHandler)
		.on('SceneItemRemoved', scenesHandler)
})
</script>
