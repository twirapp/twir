import { createGlobalState, useWebSocket } from '@vueuse/core'
import { computed, onMounted, ref, watch } from 'vue'

import { getSprite } from './dudes-config.js'
import { useDudesSettings } from './use-dudes-settings'
import { useDudes } from './use-dudes.js'

import type { TwirWebSocketEvent } from '@/api.js'
import type { DudesSprite, DudesUserSettings } from '@twir/types/overlays'

import { generateSocketUrlWithParams } from '@/helpers.js'
import { useSubscription } from '@urql/vue'
import { graphql } from '@/gql'
import type { DudesSettingsSubscriptionData } from '@/gql/graphql.ts'

declare global {
	interface GlobalEventHandlersEventMap {
		'get-user-settings': CustomEvent<string>
	}
}

export const useDudesSocket = createGlobalState(() => {
	const { dudes, createDude, updateDudeColors } = useDudes()
	const dudesSettingsStore = useDudesSettings()

	const overlayId = ref('')
	const apiKey = ref('')
	const dudesUrl = ref('')
	const { data, send, open, close, status } = useWebSocket(dudesUrl, {
		immediate: false,
		autoReconnect: {
			delay: 500,
		},
		onConnected() {
			send(JSON.stringify({ eventName: 'getSettings' }))
		},
	})

	const pauseGqlSub = computed(() => {
		return !overlayId.value || !apiKey.value
	})
	const { data: gqlData } = useSubscription({
		query: graphql(`
			subscription DudesSettings($id: UUID!, $apiKey: String!) {
				dudesSettings(id: $id, apiKey: $apiKey) {
					channelId
					channelDisplayName
					channelName
					settings {
						id
						dudeSettings {
							color
							eyesColor
							cosmeticsColor
							maxLifeTime
							gravity
							scale
							soundsEnabled
							soundsVolume
							visibleName
							growTime
							growMaxScale
							maxOnScreen
							defaultSprite
						}
						messageBoxSettings {
							enabled
							borderRadius
							boxColor
							fontFamily
							fontSize
							padding
							showTime
							fill
						}
						nameBoxSettings {
							fontFamily
							fontSize
							fill
							lineJoin
							strokeThickness
							stroke
							fillGradientStops
							fillGradientType
							fontStyle
							fontVariant
							fontWeight
							dropShadow
							dropShadowAlpha
							dropShadowAngle
							dropShadowBlur
							dropShadowDistance
							dropShadowColor
						}
						ignoreSettings {
							ignoreCommands
							ignoreUsers
							users
						}
						spitterEmoteSettings {
							enabled
						}
					}
				}
			}
		`),
		pause: pauseGqlSub,
		get variables() {
			return {
				id: overlayId.value,
				apiKey: apiKey.value,
			}
		},
		context: {},
	})

	watch(gqlData, (v) => {
		if (!v?.dudesSettings?.settings) return

		dudesSettingsStore.updateChannelData({
			channelId: v.dudesSettings.channelId,
			channelName: v.dudesSettings.channelName,
			channelDisplayName: v.dudesSettings.channelDisplayName,
		})

		updateSettingFromSocket(v.dudesSettings.settings)
	})

	watch(data, async (recieviedData) => {
		if (!dudes.value?.dudes) return

		const parsedData = JSON.parse(recieviedData) as TwirWebSocketEvent

		const data = parsedData.data as DudesUserSettings
		const dude = dudes.value.dudes.getDude(data?.userId)

		if (parsedData.eventName === 'userSettings') {
			const dudeSettings = dudesSettingsStore.dudesUserSettings.get(data.userId)
			if (!dudeSettings?.userDisplayName) return

			dudesSettingsStore.dudesUserSettings.set(data.userId, {
				...dudeSettings,
				...data,
				dudeColor: data.dudeColor ?? dudeSettings.dudeColor,
			})

			const spriteData = getSprite(
				data.dudeSprite ?? dudesSettingsStore.dudesSettings.value?.overlay.defaultSprite
			)

			const createdDude = await createDude({
				userId: dudeSettings.userId,
			})

			if (createdDude?.dude) {
				await createdDude.dude.updateSpriteData(spriteData)
				updateDudeColors(createdDude.dude, data.dudeColor)
			}
		} else if (parsedData.eventName === 'jump') {
			dude?.jump()
		} else if (parsedData.eventName === 'grow') {
			dude?.grow()
		} else if (parsedData.eventName === 'leave') {
			dude?.leave()
		} else if (parsedData.eventName === 'punished') {
			dudes.value.dudes.removeDude(data.userId)
			dudesSettingsStore.dudesUserSettings.delete(data.userId)
		}
	})

	async function updateSettingFromSocket(data: DudesSettingsSubscriptionData['settings']) {
		const fontFamily = await dudesSettingsStore.loadFont(
			data.nameBoxSettings.fontFamily,
			data.nameBoxSettings.fontWeight,
			data.nameBoxSettings.fontStyle
		)

		dudesSettingsStore.updateSettings({
			ignore: data.ignoreSettings,
			overlay: {
				defaultSprite: data.dudeSettings.defaultSprite as keyof typeof DudesSprite,
				maxOnScreen: data.dudeSettings.maxOnScreen,
			},
			dudes: {
				dude: {
					...data.dudeSettings,
					// TODO: rename and deprecate `eyes_color`, `cosmetics_color`
					bodyColor: data.dudeSettings.color,
				},
				sounds: {
					enabled: data.dudeSettings.soundsEnabled,
					volume: data.dudeSettings.soundsVolume,
				},
				name: {
					...data.nameBoxSettings,
					// TODO: move to nameBoxSettings
					enabled: data.dudeSettings.visibleName,
					fontFamily,
				},
				message: {
					...data.messageBoxSettings,
					fontFamily,
				},
				emote: data.spitterEmoteSettings,
			},
		})
	}

	function destroy(): void {
		if (status.value === 'OPEN') {
			close()
		}
	}

	function connect(connectionApiKey: string, id: string): void {
		overlayId.value = id
		apiKey.value = connectionApiKey
		dudesUrl.value = generateSocketUrlWithParams('/overlays/dudes', {
			apiKey: connectionApiKey,
			id,
		})
		open()
	}

	onMounted(() => {
		document.addEventListener('get-user-settings', (event) => {
			if (status.value !== 'OPEN') return
			send(JSON.stringify({ eventName: 'getUserSettings', data: event.detail }))
		})
	})

	return {
		destroy,
		connect,
		updateSettingFromSocket,
	}
})
