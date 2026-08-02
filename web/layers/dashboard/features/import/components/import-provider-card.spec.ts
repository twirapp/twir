import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'

import { createImportsTestI18n } from './test-i18n.js'
import ImportProviderCard from './import-provider-card.vue'

const i18n = createImportsTestI18n()

const stubs = {
	Alert: { template: '<div><slot /></div>' },
	AlertDescription: { template: '<p><slot /></p>' },
	AlertTitle: { template: '<strong><slot /></strong>' },
	Avatar: { template: '<div><slot /></div>' },
	AvatarFallback: { template: '<span><slot /></span>' },
	AvatarImage: { props: ['src'], template: '<img :src="src" />' },
	Badge: { template: '<span><slot /></span>' },
	Button: { props: ['disabled'], template: '<button :disabled="disabled"><slot /></button>' },
	Card: { template: '<section><slot /></section>' },
	CardContent: { template: '<div><slot /></div>' },
	CardDescription: { template: '<p><slot /></p>' },
	CardFooter: { template: '<footer><slot /></footer>' },
	CardHeader: { template: '<header><slot /></header>' },
	CardTitle: { template: '<h2><slot /></h2>' },
	Dialog: { template: '<div><slot /></div>' },
	DialogDescription: { template: '<p><slot /></p>' },
	DialogHeader: { template: '<header><slot /></header>' },
	DialogOrSheet: { template: '<div><slot /></div>' },
	DialogTitle: { template: '<h2><slot /></h2>' },
	NuxtIcon: { props: ['name'], template: '<i :data-icon="name" />' },
}

describe('ImportProviderCard', () => {
	it('opens provider OAuth for a disconnected account', async () => {
		const open = vi.spyOn(window, 'open').mockReturnValue(null)
		const wrapper = mount(ImportProviderCard, {
			props: {
				title: 'Nightbot',
				icon: 'twir-integrations:nightbot',
				description: 'Import Nightbot data.',
				authLink: 'https://nightbot.example/oauth',
				account: null,
				canManageIntegration: true,
			},
			global: { plugins: [i18n], stubs },
		})

		expect(wrapper.text()).toContain('Not connected')
		await wrapper.get('[data-test="provider-auth"]').trigger('click')
		expect(open).toHaveBeenCalledWith(
			'https://nightbot.example/oauth',
			'Twir connect Nightbot',
			'width=800,height=600'
		)
	})

	it('shows account controls and settings only when import is available', async () => {
		const wrapper = mount(ImportProviderCard, {
			props: {
				title: 'StreamElements',
				icon: 'twir-integrations:streamelements',
				description: 'Import data.',
				authLink: 'https://streamelements.example/oauth',
				account: { userName: 'caster', avatar: 'https://example.com/avatar.png' },
				canManageIntegration: true,
			},
			slots: { settings: '<p>Import settings</p>' },
			global: { plugins: [i18n], stubs },
		})

		expect(wrapper.text()).toContain('caster')
		expect(wrapper.find('[data-test="provider-settings"]').exists()).toBe(true)
		await wrapper.get('[data-test="provider-auth"]').trigger('click')
		expect(wrapper.emitted('logout')).toHaveLength(1)

		await wrapper.setProps({ importAvailable: false, unavailableDescription: 'Cloudbot API unavailable.' })
		expect(wrapper.text()).toContain('Unavailable')
		expect(wrapper.text()).toContain('Cloudbot API unavailable.')
		expect(wrapper.find('[data-test="provider-settings"]').exists()).toBe(false)
	})

	it('disables OAuth actions without integration permission', () => {
		const wrapper = mount(ImportProviderCard, {
			props: {
				title: 'Nightbot',
				icon: 'twir-integrations:nightbot',
				description: 'Import data.',
				authLink: 'https://nightbot.example/oauth',
				account: null,
				canManageIntegration: false,
			},
			global: { plugins: [i18n], stubs },
		})

		expect(wrapper.get('[data-test="provider-auth"]').attributes('disabled')).toBeDefined()
	})
})
