import { mount } from '@vue/test-utils'
import { afterEach, describe, expect, it } from 'vitest'
import { createI18n } from 'vue-i18n'

import de from '../../../../../locales/de.json'
import en from '../../../../../locales/en.json'
import es from '../../../../../locales/es.json'
import ja from '../../../../../locales/ja.json'
import pt from '../../../../../locales/pt.json'
import ru from '../../../../../locales/ru.json'
import sk from '../../../../../locales/sk.json'
import uk from '../../../../../locales/uk.json'
import { Platform } from '~/gql/graphql.js'
import StatsLayoutRows from './stats-layout-rows.vue'

import type { DashboardStats } from '../composables/use-platform-stats.js'

const stats: DashboardStats = {
	categoryId: '1',
	categoryName: 'Software and Game Development',
	viewers: 33,
	startedAt: new Date(Date.now() - 3_600_000).toISOString(),
	title: 'Twitch title',
	chatMessages: 164,
	followers: 686,
	usedEmotes: 0,
	requestedSongs: 0,
	subs: 3,
	platforms: [
		{
			platform: Platform.VkVideoLive,
			isLive: false,
			title: 'VK title',
			categoryId: null,
			categoryName: 'Programming',
			viewers: null,
			followers: null,
			startedAt: null,
			chatMessages: 0,
			usedEmotes: 0,
			canEditInfo: false,
		},
		{
			platform: Platform.Kick,
			isLive: true,
			title: 'Kick title',
			categoryId: '2',
			categoryName: 'Just Chatting',
			viewers: 12,
			followers: 48,
			startedAt: new Date(Date.now() - 3_600_000).toISOString(),
			chatMessages: 10,
			usedEmotes: 2,
			canEditInfo: true,
		},
		{
			platform: Platform.Twitch,
			isLive: true,
			title: 'Twitch title',
			categoryId: '1',
			categoryName: 'Software and Game Development',
			viewers: 21,
			followers: 126,
			startedAt: new Date(Date.now() - 7_200_000).toISOString(),
			chatMessages: 154,
			usedEmotes: 5,
			canEditInfo: true,
		},
	],
}

const i18n = createI18n({
	legacy: false,
	globalInjection: true,
	locale: 'en',
	messages: { de, en, es, ja, pt, ru, sk, uk },
})

const mountOptions = {
	global: {
		plugins: [i18n],
		stubs: {
			NuxtIcon: true,
		},
	},
}

describe('StatsLayoutRows', () => {
	afterEach(() => {
		document.body.replaceChildren()
	})

	it('renders the titles widget and global widgets with per-platform rows', () => {
		const wrapper = mount(StatsLayoutRows, {
			...mountOptions,
			props: { stats },
		})

		const text = wrapper.text()
		expect(text).toContain('Stream titles')
		expect(text).toContain('Viewers')
		expect(text).toContain('Followers')
		expect(text).toContain('Messages')
		expect(text).toContain('21')
		expect(text).toContain('12')
		// offline VK viewers shown as a dimmed dash
		expect(text).toContain('—')
	})

	it('keeps the edit-mode button for widget config', () => {
		const wrapper = mount(StatsLayoutRows, {
			...mountOptions,
			props: { stats },
		})

		expect(wrapper.text()).toContain('Edit')
	})
})
