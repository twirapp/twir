import { DudesSprite } from '@twir/types/overlays'
import { DudesLayers } from '@twirapp/dudes-vue'
import { createGlobalState } from '@vueuse/core'
import { computed, ref, watch } from 'vue'

import { getSprite } from './dudes-config.js'
import { useDudesSettings } from './use-dudes-settings.js'

import type { MessageChunk } from '@twir/frontend-chat'
import type DudesOverlay from '@twirapp/dudes-vue'
import type { Dude } from '@twirapp/dudes-vue/types'

import { randomRgbColor } from '@/helpers.js'

export const useDudes = createGlobalState(() => {
	const { dudesSettings, dudesUserSettings } = useDudesSettings()

	const dudes = ref<InstanceType<typeof DudesOverlay> | null>(null)
	const isDudeReady = ref(false)
	const isDudeOverlayReady = computed(() => {
		return dudes.value && dudesSettings.value && isDudeReady.value
	})

	function createDudeInstance(dude: Dude) {
		return {
			dude,
			isCreated: false,
			showMessage(messageChunks: MessageChunk[]) {
				if (this.isCreated) {
					setTimeout(() => showMessageDude(this.dude, messageChunks), 1000)
				} else {
					showMessageDude(this.dude, messageChunks)
				}
			},
		}
	}

	async function createDude(
		{ userId, userName, color }: { userId: string, userName?: string, color?: string },
	) {
		if (!dudes.value?.dudes || !dudesSettings.value) return

		const actualDude = dudes.value.dudes.getDude(userId)
		if (actualDude) {
			return createDudeInstance(actualDude)
		}

		if (
			dudesSettings.value.overlay.maxOnScreen !== 0
			&& dudes.value.dudes.dudes.size === dudesSettings.value.overlay.maxOnScreen
		) return

		const userSettings = requestDudeUserSettings(userId)
		if (!userSettings) {
			dudesUserSettings.set(userId, {
				userId,
				userDisplayName: userName,
				dudeColor: color,
			})
			return
		}

		const dudeColor = userSettings.dudeColor
		  ?? color
		  ?? dudesSettings.value.dudes.dude.bodyColor

		const dudeSprite = getSprite(userSettings?.dudeSprite ?? dudesSettings.value.overlay.defaultSprite)
		const dude = await dudes.value.dudes.createDude({
			id: userSettings.userId,
			name: userSettings.userDisplayName!,
			sprite: dudeSprite,
		})

		updateDudeColors(dude, dudeColor)

		const dudeInstance = createDudeInstance(dude)
		dudeInstance.isCreated = true

		return dudeInstance
	}

	function updateDudeColors(dude: Dude, color?: string): void {
		if (color) {
			dude.updateColor(DudesLayers.body, color)
		}

		if (dude.config.sprite.name.startsWith(DudesSprite.girl)) {
			dude.updateColor(DudesLayers.hat, '#FF0000')
		}

		if (dude.config.sprite.name.startsWith(DudesSprite.santa)) {
			dude.updateColor(DudesLayers.hat, '#FFF')
		}

		if (dude.config.sprite.name.startsWith(DudesSprite.agent)) {
			dude.updateColor(DudesLayers.cosmetics, '#8a2be2')
		}

		if (dude.config.sprite.name.startsWith(DudesSprite.sith)) {
			dude.updateColor(DudesLayers.cosmetics, randomRgbColor())
		}
	}

	function showMessageDude(dude: Dude, messageChunks: MessageChunk[]): void {
		if (
			dudesSettings.value?.ignore.ignoreCommands
			&& messageChunks?.at(0)?.value.startsWith('!')
		) {
			return
		}

		const message = messageChunks
			.filter((chunk) => chunk.type === 'text')
			.map((chunk) => chunk.value)
			.join(' ')

		dude.addMessage(message)

		const emotes = messageChunks
			.filter((chunk) => chunk.type !== 'text')
			.map(getProxiedEmoteUrl)

		if (emotes.length) {
			dude.addEmotes([...new Set(emotes)])
		}
	}

	function getProxiedEmoteUrl(messageChunk: MessageChunk): string {
		if (messageChunk.type === 'emoji') {
			const code = messageChunk.value.codePointAt(0)?.toString(16)
			return `https://cdn.frankerfacez.com/static/emoji/images/twemoji/${code}.png`
		}

		return `${window.location.origin}/api-old/proxy?url=${messageChunk.value}`
	}

	function requestDudeUserSettings(userId: string) {
		const userSettings = dudesUserSettings.get(userId)
		if (userSettings) return userSettings
		document.dispatchEvent(new CustomEvent<string>('get-user-settings', { detail: userId }))
	}

	watch(() => dudes.value, async (dudes) => {
		if (!dudes) return
		await dudes.initDudes()
		isDudeReady.value = true
	})

	return {
		dudes,
		createDude,
		showMessageDude,
		updateDudeColors,
		getProxiedEmoteUrl,
		isDudeOverlayReady,
	}
})
