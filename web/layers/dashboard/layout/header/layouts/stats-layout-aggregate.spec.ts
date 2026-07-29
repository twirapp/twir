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
import StatsLayoutAggregate from './stats-layout-aggregate.vue'

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

describe('StatsLayoutAggregate', () => {
	afterEach(() => {
		document.body.replaceChildren()
	})

	it('shows the primary platform title and aggregated totals', () => {
		const wrapper = mount(StatsLayoutAggregate, {
			...mountOptions,
			props: { stats },
		})

		const text = wrapper.text()
		// primary platform is the first live one (Twitch)
		expect(text).toContain('Twitch title')
		expect(text).toContain('Software and Game Development')
		// totals: 21 + 12 viewers, 126 + 48 followers (VK null skipped)
		expect(text).toContain('33')
		expect(text).toContain('174')
	})

	it('emits editStreamInfo with the primary platform when the title widget is clicked', async () => {
		const wrapper = mount(StatsLayoutAggregate, {
			...mountOptions,
			props: { stats },
		})

		const titleWidget = wrapper.get('.header-widget-stream')
		await titleWidget.trigger('click')
		expect(wrapper.emitted('editStreamInfo')).toHaveLength(1)
		expect(wrapper.emitted('editStreamInfo')?.[0]).toEqual([Platform.Twitch])
	})

	it('emits the Kick platform when Kick is the editable primary platform', async () => {
		const kickPrimary: DashboardStats = {
			...stats,
			platforms: stats.platforms.map((platform) =>
				platform.platform === Platform.Twitch
					? { ...platform, isLive: false }
					: platform
			),
		}

		const wrapper = mount(StatsLayoutAggregate, {
			...mountOptions,
			props: { stats: kickPrimary },
		})

		const titleWidget = wrapper.get('.header-widget-stream')
		await titleWidget.trigger('click')
		expect(wrapper.emitted('editStreamInfo')?.[0]).toEqual([Platform.Kick])
	})
})
