import { useSubscription } from '@urql/vue'
import { computed, ref, watch } from 'vue'

import type { BrbOnStartFn, BrbOnStopFn } from '@/types.js'

import { graphql } from '@/gql'

import { useBrbSettings } from './use-brb-settings.js'

type Options = {
	onStart: BrbOnStartFn
	onStop: BrbOnStopFn
}

export function useBeRightBackOverlayGraphQL(options: Options) {
	const apiKey = ref<string>('')
	const paused = computed(() => !apiKey.value)
	const { setSettings } = useBrbSettings()

	// Subscribe to settings updates
	const {
		data: settingsData,
		executeSubscription: connectSettings,
		pause: pauseSettings,
	} = useSubscription({
		query: graphql(`
			subscription BeRightBackSettings($apiKey: String!) {
				overlaysBeRightBack(apiKey: $apiKey) {
					id
					text
					late {
						enabled
						text
						displayBrbTime
					}
					backgroundColor
					fontSize
					fontColor
					fontFamily
					createdAt
					updatedAt
					channelId
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

	// Subscribe to start events
	const {
		data: startData,
		executeSubscription: connectStart,
		pause: pauseStart,
	} = useSubscription({
		query: graphql(`
			subscription BeRightBackStart($apiKey: String!) {
				overlaysBeRightBackStart(apiKey: $apiKey) {
					time
					text
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

	// Subscribe to stop events
	const {
		data: stopData,
		executeSubscription: connectStop,
		pause: pauseStop,
	} = useSubscription({
		query: graphql(`
			subscription BeRightBackStop($apiKey: String!) {
				overlaysBeRightBackStop(apiKey: $apiKey)
			}
		`),
		get variables() {
			return {
				apiKey: apiKey.value,
			}
		},
		pause: paused,
	})

	// Watch for settings updates
	watch(settingsData, (data) => {
		if (!data?.overlaysBeRightBack) return

		const overlay = data.overlaysBeRightBack
		setSettings({
			text: overlay.text,
			late: {
				enabled: overlay.late.enabled,
				text: overlay.late.text,
				displayBrbTime: overlay.late.displayBrbTime,
			},
			backgroundColor: overlay.backgroundColor,
			fontSize: overlay.fontSize,
			fontColor: overlay.fontColor,
			fontFamily: overlay.fontFamily,
			channelId: overlay.channelId,
		})
	})

	// Watch for start events
	watch(startData, (data) => {
		if (!data?.overlaysBeRightBackStart) return

		const { time, text } = data.overlaysBeRightBackStart
		options.onStart(time, text || '')
	})

	// Watch for stop events
	watch(stopData, (data) => {
		if (!data?.overlaysBeRightBackStop) return
		options.onStop()
	})

	function destroy() {
		pauseSettings()
		pauseStart()
		pauseStop()
	}

	async function connect(key: string) {
		apiKey.value = key
		connectSettings()
		connectStart()
		connectStop()
	}

	return {
		connect,
		destroy,
		settingsData,
	}
}
