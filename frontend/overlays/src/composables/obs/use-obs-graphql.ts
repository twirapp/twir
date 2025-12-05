import { useMutation, useSubscription } from '@urql/vue'
import { computed, onUnmounted, ref, watch } from 'vue'

import type { ObsWebsocketCommandAction } from '@/gql/graphql'

import { graphql } from '@/gql'

export type ObsSettings = {
	serverAddress: string
	serverPort: number
	serverPassword: string
}

export type ObsCommand = {
	action: ObsWebsocketCommandAction
	target: string
	volumeValue?: number | null
	volumeStep?: number | null
}

type Options = {
	onSettings: (settings: ObsSettings) => void
	onCommand: (command: ObsCommand) => void
}

// Heartbeat interval in milliseconds (should be less than Redis TTL of 5 seconds)
const HEARTBEAT_INTERVAL = 1000

export function useObsOverlayGraphQL(options: Options) {
	const apiKey = ref<string>('')
	const paused = computed(() => !apiKey.value)
	let heartbeatInterval: ReturnType<typeof setInterval> | null = null

	// Subscribe to settings updates
	const {
		data: settingsData,
		executeSubscription: connectSettings,
		pause: pauseSettings,
	} = useSubscription({
		query: graphql(`
			subscription ObsWebsocketSettings($apiKey: String!) {
				obsWebsocketData(apiKey: $apiKey) {
					serverPort
					serverAddress
					serverPassword
					sources
					audioSources
					scenes
					isConnected
				}
			}
		`),
		get variables() {
			return {
				apiKey: apiKey.value,
			}
		},
		pause: paused,
	})

	// Subscribe to commands
	const {
		data: commandsData,
		executeSubscription: connectCommands,
		pause: pauseCommands,
	} = useSubscription({
		query: graphql(`
			subscription ObsWebsocketCommands($apiKey: String!) {
				obsWebsocketCommands(apiKey: $apiKey) {
					action
					target
					volumeValue
					volumeStep
				}
			}
		`),
		get variables() {
			return {
				apiKey: apiKey.value,
			}
		},
		pause: paused,
	})

	// Mutation to update overlay data
	const updateFromOverlayMutation = useMutation(
		graphql(`
			mutation ObsWebsocketUpdateFromOverlay($apiKey: String!, $input: ObsWebsocketOverlayInput!) {
				obsWebsocketUpdateFromOverlay(apiKey: $apiKey, input: $input)
			}
		`)
	)

	// Mutation to set connected state
	const setConnectedMutation = useMutation(
		graphql(`
			mutation ObsWebsocketSetConnected($apiKey: String!) {
				obsWebsocketSetConnected(apiKey: $apiKey)
			}
		`)
	)

	// Send heartbeat to keep connection status alive
	async function sendHeartbeat() {
		if (!apiKey.value) return
		try {
			await setConnectedMutation.executeMutation({ apiKey: apiKey.value })
		} catch (error) {
			console.error('Failed to send OBS connection heartbeat:', error)
		}
	}

	function startHeartbeat() {
		// Send initial heartbeat
		sendHeartbeat()
		// Start periodic heartbeat
		heartbeatInterval = setInterval(sendHeartbeat, HEARTBEAT_INTERVAL)
	}

	function stopHeartbeat() {
		if (heartbeatInterval) {
			clearInterval(heartbeatInterval)
			heartbeatInterval = null
		}
	}

	// Watch for settings updates
	watch(settingsData, (data, oldValue) => {
		if (!data?.obsWebsocketData) return

		if (
			data.obsWebsocketData.serverAddress === oldValue?.obsWebsocketData.serverAddress &&
			data.obsWebsocketData.serverPort === oldValue?.obsWebsocketData.serverPort &&
			data.obsWebsocketData.serverPassword === oldValue?.obsWebsocketData.serverPassword
		) {
			return
		}

		const settings = data.obsWebsocketData
		options.onSettings({
			serverAddress: settings.serverAddress,
			serverPort: settings.serverPort,
			serverPassword: settings.serverPassword,
		})
	})

	// Watch for commands
	watch(commandsData, (data) => {
		if (!data?.obsWebsocketCommands) return

		const cmd = data.obsWebsocketCommands
		options.onCommand({
			action: cmd.action,
			target: cmd.target,
			volumeValue: cmd.volumeValue,
			volumeStep: cmd.volumeStep,
		})
	})

	function destroy() {
		stopHeartbeat()
		pauseSettings()
		pauseCommands()
	}

	async function connect(key: string) {
		apiKey.value = key
		connectSettings()
		connectCommands()
		// Note: heartbeat should be started manually after OBS connection is established
	}

	async function updateSources(input: {
		scenes: string[]
		sources: string[]
		audioSources: string[]
	}) {
		await updateFromOverlayMutation.executeMutation({
			apiKey: apiKey.value,
			input,
		})
	}

	// Cleanup on unmount
	onUnmounted(() => {
		stopHeartbeat()
	})

	return {
		connect,
		destroy,
		updateSources,
		settingsData,
		startHeartbeat,
		stopHeartbeat,
	}
}
