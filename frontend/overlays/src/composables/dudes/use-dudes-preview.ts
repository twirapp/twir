import { type DudeSprite, getSprite } from './dudes-config.js'
import { useDudes } from './use-dudes.js'

import { randomEmoji, randomRgbColor } from '@/helpers.js'

interface PreviewDude {
	id: string
	name: string
	sprite: DudeSprite
}

const previewDudes: PreviewDude[] = [
	{ id: 'preview-potato', name: 'Potato', sprite: 'dude' },
	{ id: 'preview-noodle', name: 'Noodle', sprite: 'cat' },
	{ id: 'preview-biscuit', name: 'Biscuit', sprite: 'girl' },
]

const previewMessages = ['Hello!', 'HeyGuys', 'Pog', 'Nice stream', 'lol', 'Hype!']

const previewEmoteUrls = [
	'https://cdn.7tv.app/emote/60b00d1f0d3a78a196f803e3/1x.gif',
	'https://cdn.7tv.app/emote/65413498dc0468e8c1fbcdc6/1x.gif',
]

export function useDudesPreview() {
	const { dudes, updateDudeColors, getProxiedEmoteUrl } = useDudes()

	let interval: ReturnType<typeof setInterval> | null = null

	async function spawn() {
		if (!dudes.value?.dudes) return

		for (const mock of previewDudes) {
			if (dudes.value.dudes.getDude(mock.id)) continue

			const dude = await dudes.value.dudes.createDude({
				id: mock.id,
				name: mock.name,
				sprite: getSprite(mock.sprite),
			})

			updateDudeColors(dude, randomRgbColor())
			dude.addMessage(`${mock.name} says hi! ${randomEmoji('emoticons')}`)
		}
	}

	function tick() {
		if (!dudes.value?.dudes) return

		const mock = previewDudes[Math.floor(Math.random() * previewDudes.length)]
		const dude = dudes.value.dudes.getDude(mock.id)
		if (!dude) return

		const roll = Math.random()
		if (roll < 0.45) {
			const text = previewMessages[Math.floor(Math.random() * previewMessages.length)]
			dude.addMessage(`${text} ${randomEmoji('emoticons')}`)
		} else if (roll < 0.75) {
			const emoteUrl = previewEmoteUrls[Math.floor(Math.random() * previewEmoteUrls.length)]
			dude.addEmotes([getProxiedEmoteUrl({ type: '3rd_party_emote', value: emoteUrl })])
		} else {
			dude.jump()
		}
	}

	function start() {
		stop()
		void spawn()
		interval = setInterval(tick, 4000)
	}

	function stop() {
		if (!interval) return
		clearInterval(interval)
		interval = null
	}

	return { start, stop }
}
