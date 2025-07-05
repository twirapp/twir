import { createGlobalState } from '@vueuse/core'

import { dudeMock, getSprite } from './dudes-config.js'
import { useDudesSettings } from './use-dudes-settings.js'
import { useDudesSocket } from './use-dudes-socket.js'
import { useDudes } from './use-dudes.js'
import type { DudesOverlaySettings } from '@/gql/graphql'
import { randomEmoji } from '@/helpers.js'

interface DudesPostMessage {
	action: string
	data?: any
}

export const useDudesIframe = createGlobalState(() => {
	const isIframe = Boolean(window.frameElement)
	const { dudes, updateDudeColors, getProxiedEmoteUrl } = useDudes()
	const dudesSocketStore = useDudesSocket()
	const dudesSettingsStore = useDudesSettings()

	async function onPostMessage(msg: MessageEvent<string>) {
		if (!dudes.value?.dudes) return

		const parsedData = JSON.parse(msg.data) as DudesPostMessage
		const dude = dudes.value.dudes.getDude(dudeMock.id)

		if (parsedData.action === 'reset') {
			dudes.value.dudes.removeAllDudes()
			spawnIframeDude()
			return
		}

		if (parsedData.data) {
			if (parsedData.action === 'update-settings') {
				const settings = parsedData.data as DudesOverlaySettings
				dudesSocketStore.updateSettingFromSocket(settings)
				return
			}

			if (!dude) return

			if (parsedData.action === 'update-sprite') {
				const spriteData = getSprite(parsedData.data)
				await dude.updateSpriteData(spriteData)
				updateDudeColors(dude)
				return
			}

			if (parsedData.action === 'update-color') {
				updateDudeColors(dude, parsedData.data)
				return
			}
		}

		if (dude) {
			if (parsedData.action === 'jump') {
				dude.jump()
			}

			if (parsedData.action === 'grow') {
				dude.grow()
			}

			if (parsedData.action === 'leave') {
				dude.leave()
			}

			if (parsedData.action === 'spawn-emote') {
				const emote = getProxiedEmoteUrl({
					type: '3rd_party_emote',
					value: 'https://cdn.7tv.app/emote/60b00d1f0d3a78a196f803e3/1x.gif',
				})
				dude.addEmotes([emote])
			}

			if (parsedData.action === 'show-message') {
				dude.addMessage(
					`Hello, ${dudesSettingsStore.channelData.value!.channelDisplayName}! ${randomEmoji('emoticons')}`
				)
			}
		}
	}

	async function spawnIframeDude() {
		if (
			!dudes.value?.dudes ||
			!dudesSettingsStore.dudesSettings.value ||
			!dudesSettingsStore.channelData.value ||
			dudes.value.dudes.getDude(dudeMock.id)
		)
			return

		const emote = getProxiedEmoteUrl({
			type: '3rd_party_emote',
			value: 'https://cdn.7tv.app/emote/65413498dc0468e8c1fbcdc6/1x.gif',
		})

		const dudeSprite = getSprite(dudesSettingsStore.dudesSettings.value.overlay.defaultSprite)
		const dude = await dudes.value.dudes.createDude({
			id: dudeMock.id,
			name: dudeMock.name,
			sprite: dudeSprite,
		})

		updateDudeColors(dude, dudeMock.color)
		dude.addMessage(
			`Hello, ${dudesSettingsStore.channelData.value.channelDisplayName}! ${randomEmoji('emoticons')}`
		)
		dude.addEmotes([emote])
	}

	function connect() {
		if (!isIframe) return
		window.addEventListener('message', onPostMessage)
	}

	function destroy() {
		if (!isIframe) return
		window.removeEventListener('message', onPostMessage)
	}

	return {
		isIframe,
		spawnIframeDude,
		connect,
		destroy,
	}
})
