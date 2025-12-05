import { useMutation, useSubscription } from '@urql/vue'
import { computed, ref, watch } from 'vue'

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

export function useObsOverlayGraphQL(options: Options) {
	const apiKey = ref<string>('')
	const paused = computed(() => !apiKey.value)

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

	// Watch for settings updates
	watch(settingsData, (data) => {
		if (!data?.obsWebsocketData) return

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
		pauseSettings()
		pauseCommands()
	}

	async function connect(key: string) {
		apiKey.value = key
		connectSettings()
		connectCommands()
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

	return {
		connect,
		destroy,
		updateSources,
		settingsData,
	}
}
