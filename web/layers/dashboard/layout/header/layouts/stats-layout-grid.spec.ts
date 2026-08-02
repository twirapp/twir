import { mount } from '@vue/test-utils'
import { afterEach, describe, expect, it } from 'vitest'
import { createI18n } from 'vue-i18n'

import de from '../../../../../i18n/locales/de.json'
import en from '../../../../../i18n/locales/en.json'
import es from '../../../../../i18n/locales/es.json'
import ja from '../../../../../i18n/locales/ja.json'
import pt from '../../../../../i18n/locales/pt.json'
import ru from '../../../../../i18n/locales/ru.json'
import sk from '../../../../../i18n/locales/sk.json'
import uk from '../../../../../i18n/locales/uk.json'
import { Platform } from '~/gql/graphql.js'
import StatsLayoutGrid from './stats-layout-grid.vue'

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

describe('StatsLayoutGrid', () => {
	afterEach(() => {
		document.body.replaceChildren()
	})

	it('renders one card per platform, live first, without global stats widgets', () => {
		const wrapper = mount(StatsLayoutGrid, {
			...mountOptions,
			props: { stats },
		})

		const text = wrapper.text()
		expect(text).toContain('Twitch title')
		expect(text).toContain('Kick title')
		expect(text).toContain('VK title')
		// live platforms show their viewers, offline VK shows a dash
		expect(text).toContain('21')
		expect(text).toContain('12')
		expect(text).toContain('—')
		// no global stats widgets in the grid layout
		expect(text).not.toContain('Messages')
		expect(wrapper.find('[draggable="true"]').exists()).toBe(false)
	})

	it('does not dim offline platform strips, but keeps the offline badge', () => {
		const wrapper = mount(StatsLayoutGrid, {
			...mountOptions,
			props: { stats },
		})

		const strips = wrapper.findAll('.header-widget')
		expect(strips).toHaveLength(3)
		for (const strip of strips) {
			expect(strip.classes()).not.toContain('opacity-60')
		}
		expect(strips[2]!.text()).toContain('Offline')
	})

	it('emits editStreamInfo from Twitch and Kick pencils only', async () => {
		const wrapper = mount(StatsLayoutGrid, {
			...mountOptions,
			props: { stats },
		})

		// VK (canEditInfo=false) renders no pencil at all
		const editButtons = wrapper.findAll('button[title]')
		expect(editButtons).toHaveLength(2)

		await editButtons[0]!.trigger('click')
		expect(wrapper.emitted('editStreamInfo')).toHaveLength(1)
		expect(wrapper.emitted('editStreamInfo')?.[0]).toEqual([Platform.Twitch])

		await editButtons[1]!.trigger('click')
		expect(wrapper.emitted('editStreamInfo')?.[1]).toEqual([Platform.Kick])
	})
})
