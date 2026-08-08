import { createGlobalState, useWebSocket } from '@vueuse/core'
import { ref, watch } from 'vue'

import type { TwirWebSocketEvent } from '@/api.js'

import { generateSocketUrlWithParams } from '@/helpers.js'

export interface DotaState {
	channelId: string
	inGame: boolean
	mmr: number
	sessionWins: number
	sessionLosses: number
	winProbability: number
	heroName: string
	matchId: number
	teamIsRadiant: boolean
	teamKnown: boolean
}

export const useDotaSocket = createGlobalState(() => {
	const state = ref<DotaState | null>(null)

	const socketUrl = ref('')
	const { data, open, close, status } = useWebSocket(socketUrl, {
		immediate: false,
		autoReconnect: {
			delay: 1000,
		},
	})

	watch(data, (raw) => {
		if (typeof raw !== 'string') return

		let event: TwirWebSocketEvent<DotaState>
		try {
			event = JSON.parse(raw)
		} catch {
			return
		}

		if (event.eventName !== 'dotaStateUpdate') return
		state.value = event.data
	})

	function connect(apiKey: string) {
		socketUrl.value = generateSocketUrlWithParams('/overlays/dota', { apiKey })
		open()
	}

	return {
		state,
		status,
		connect,
		close,
	}
})
